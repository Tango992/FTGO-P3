package config

import (
	"log"
	pb "ungraded_5/internal/product"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitGrpcClient() (*grpc.ClientConn, pb.ProductServiceClient) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return conn, pb.NewProductServiceClient(conn)
}