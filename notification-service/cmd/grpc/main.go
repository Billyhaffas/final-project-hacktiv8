package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	grpchandler "notification-service/internal/delivery/grpc"
	"notification-service/internal/repository"
	"notification-service/internal/usecase"
	pbemission "notification-service/proto/emission"
	pb "notification-service/proto/generated"
)

// grpcCreds returns TLS credentials for Heroku (GRPC_USE_TLS=true) or insecure for local/docker.
func grpcCreds() grpc.DialOption {
	if os.Getenv("GRPC_USE_TLS") == "true" {
		return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, ""))
	}
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}

	emissionAddr := os.Getenv("COUNT_EMISSION_SERVICE_ADDR")
	if emissionAddr == "" {
		emissionAddr = "count-emission-service:50051"
	}

	emissionConn, err := grpc.NewClient(emissionAddr, grpcCreds())
	if err != nil {
		log.Fatalf("failed to connect to count-emission-service at %s: %v", emissionAddr, err)
	}
	defer emissionConn.Close()

	emissionClient := repository.NewEmissionGRPCClient(pbemission.NewEmissionClient(emissionConn))
	notifUsecase := usecase.NewNotificationUsecase(emissionClient)
	notifHandler := grpchandler.NewNotificationGRPCServer(notifUsecase)

	server := grpc.NewServer()
	pb.RegisterNotificationServer(server, notifHandler)
	reflection.Register(server)

	port := os.Getenv("PORT") // Heroku injects $PORT
	if port == "" {
		port = os.Getenv("NOTIFICATION_GRPC_PORT")
	}
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	log.Printf("notification-service gRPC server running on :%s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
