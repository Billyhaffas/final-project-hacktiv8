package main

import (
	grpchandler "count-emission-service/internal/delivery/gRPC"
	"count-emission-service/internal/infrastructure/database"
	externalapi "count-emission-service/internal/repository/external_api"
	repository "count-emission-service/internal/repository/internal_repo"
	"count-emission-service/internal/usecase"
	pb "count-emission-service/proto/generated"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

func main() {
	db, err := database.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	// server := grpc.NewServer(
	// 	grpc.UnaryInterceptor(jwtMiddleware.Unary()),
	// )

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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	if err := server.Serve(lis); err != nil {
		panic(err)
	}

}
