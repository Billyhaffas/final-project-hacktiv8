package main

import (
	grpchandler "count-emission-service/internal/delivery/gRPC"
	"count-emission-service/internal/infrastructure/database"
	externalapi "count-emission-service/internal/repository/external_api"
	repository "count-emission-service/internal/repository/internal_repo"
	"count-emission-service/internal/usecase"
	pb "count-emission-service/proto/generated"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := database.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	reflection.Register(server)

	emissionRepo := repository.NewEmissionCollection(db)
	preferenceRepo := repository.NewPreferenceRepository(db)
	carbonSutraRepo := externalapi.NewCachedEmissionRepo(externalapi.NewCarbonSutraRepository(httpClient))
	emissionUseCase := usecase.NewEmissionUseCase(emissionRepo, carbonSutraRepo)
	preferenceUseCase := usecase.NewPreferenceUseCase(preferenceRepo)
	emissionHandler := grpchandler.NewEmissionGRPCServer(emissionUseCase, preferenceUseCase)

	pb.RegisterEmissionServer(server, emissionHandler)

	port := os.Getenv("PORT") // Heroku injects $PORT
	if port == "" {
		port = os.Getenv("EMISSION_GRPC_PORT")
	}
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}
	log.Printf("count-emission-service gRPC server running on :%s", port)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
