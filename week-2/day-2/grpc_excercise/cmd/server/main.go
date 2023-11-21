package main

import (
	"context"
	pb "grpc_test/internal/user"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	users []*pb.User
	pb.UnimplementedUserServiceServer
	collection *mongo.Collection
}

func (s *Server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.User, error) {
	newUser := &pb.User{
		Id: req.Id,
		Name: req.Name,
	}

	_, err := s.collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}
	return newUser, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	filerById := bson.M{"id": req.Id}

	if _, err := s.collection.DeleteOne(ctx, filerById); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error){
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	users := []*pb.User{}
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)		
	}
	
	UserResp := &pb.GetUserResponse{Users: users}
	return UserResp, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	filerById := bson.M{"id": req.Id}
	updateData := bson.M{"$set": bson.M{
		"name": req.Name,
	}}

	if _, err := s.collection.UpdateOne(ctx, filerById, updateData); err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}
	return &pb.User{Id: req.Id, Name: req.Name}, nil
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
	pb.RegisterUserServiceServer(grpcServer, &Server{collection: collection})
	log.Println("Server is running on port: 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to server gRPC:", err)
	}
}