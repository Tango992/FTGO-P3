package main

import (
	"context"
	"log"
	"net"
	pb "ungraded_5/internal/product"

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
	
	products := []*pb.Product{}
	if err := cursor.All(context.TODO(), &products); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)		
	}

	productsResp := &pb.Products{Products: products}
	return productsResp, nil
}

func (s *Server) UpdateProduct(ctx context.Context, data *pb.UpdateProductRequest) (*emptypb.Empty, error) {
	objectId, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if _, err := s.collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": data}); err != nil {
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
		if err == mongo.ErrNilDocument {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	var product *pb.Product
	if err := result.Decode(&product); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}
	return product, nil
}

func main() {
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
	
	collection := client.Database("test_grpc_db").Collection("users")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, &Server{collection: collection})
	log.Println("Server is running on port: 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server gRPC:", err)
	}
}