package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	pb "ungraded_5/internal/product"
	"ungraded_5/middlewares"
	"ungraded_5/models"

	"github.com/dgrijalva/jwt-go"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedProductServiceServer
	collection *mongo.Collection
}

func (s *Server) AddProduct(ctx context.Context, data *pb.AddProductRequest) (*pb.Product, error) {
	res, err := s.collection.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	objectId := res.InsertedID.(primitive.ObjectID).Hex()
	product := &pb.Product{
		Id: objectId,
		Name: data.Name,
		Stock: data.Stock,
	}
	return product, nil
}

func (s *Server) GetProducts(ctx context.Context, data *emptypb.Empty) (*pb.Products, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}
	defer cursor.Close(context.Background())
	
	products := []*pb.Product{}
	for cursor.Next(context.Background()) {
		var data models.Product
		if err := cursor.Decode(&data); err != nil {
			return nil, status.Errorf(codes.Internal, "Internal error: %v", err)		
		}

		product := &pb.Product{
			Id: data.Id.Hex(),
			Name: data.Name,
			Stock: uint32(data.Stock),
		}
		products = append(products, product)
	}
	productsResp := &pb.Products{Products: products}
	return productsResp, nil
}

func (s *Server) UpdateProduct(ctx context.Context, data *pb.UpdateProductRequest) (*emptypb.Empty, error) {
	objectId, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	updateData := models.Product{
		Name: data.Name,
		Stock: uint(data.Stock),
	}

	res := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": objectId}, bson.M{"$set": updateData})
	if err := res.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, data *pb.DeleteProductRequest) (*pb.Product, error) {
	objectId, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	result := s.collection.FindOneAndDelete(ctx, bson.M{"_id": objectId})
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	var productTmp models.Product
	if err := result.Decode(&productTmp); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	product := &pb.Product{
		Id: productTmp.Id.Hex(),
		Name: productTmp.Name,
		Stock: uint32(productTmp.Stock),
	}
	return product, nil
}

func main() {
	token, err := generateToken()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("JWT Token:", token)
	
	clientOption := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	defer func(){
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	
	collection := client.Database("ungraded5_db").Collection("products")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(middlewares.JWTAuth)))
	pb.RegisterProductServiceServer(grpcServer, &Server{collection: collection})
	log.Println("Server is running on port: 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server gRPC:", err)
	}
}

func generateToken() (string, error) {
    // Atur klaim JWT Anda
    claims := jwt.StandardClaims{
        Subject:   "testuser",                           // Subjek/token penerima
        ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // Waktu kedaluwarsa token (1 jam dari sekarang)
        IssuedAt:  time.Now().Unix(),                    // Waktu pembuatan token
    }

    // Buat token dengan menggunakan metode HMAC dan kunci rahasia
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secretKey := []byte(os.Getenv("JWT_SECRET")) // Ganti dengan kunci rahasia Anda
    tokenString, err := token.SignedString(secretKey)
    if err != nil {
        return "", fmt.Errorf("failed to generate JWT: %v", err)
    }

    return tokenString, nil
}