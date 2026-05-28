package main

import (
	"convert-emission-service/config"
	grpchandler "convert-emission-service/internal/delivery/grpc"
	"convert-emission-service/internal/repository"
	"convert-emission-service/internal/usecase"
	pb "convert-emission-service/proto/generated"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}

	col := config.ConnectMongo()

	repo := repository.NewCarbonPriceRepository(col)
	uc := usecase.NewConversionUsecase(repo)
	handler := grpchandler.NewConvertGRPCServer(uc)

	server := grpc.NewServer()
	pb.RegisterConvertServer(server, handler)
	reflection.Register(server)

	port := os.Getenv("CONVERT_GRPC_PORT")
	if port == "" {
		port = "50053"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	log.Printf("convert-emission-service gRPC server running on :%s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
