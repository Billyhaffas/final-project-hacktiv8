package main

import (
	"api-gateway/internal/handler"
	"api-gateway/internal/router"
	pb "api-gateway/proto/emission"
	pbconvert "api-gateway/proto/convert"
	pbnotif "api-gateway/proto/notification"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}

	emissionConn, err := grpc.NewClient(
		os.Getenv("COUNT_EMISSION_SERVICE_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to count-emission-service: %v", err)
	}
	defer emissionConn.Close()

	notifConn, err := grpc.NewClient(
		os.Getenv("NOTIFICATION_SERVICE_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to notification-service: %v", err)
	}
	defer notifConn.Close()

	convertConn, err := grpc.NewClient(
		os.Getenv("CONVERT_EMISSION_SERVICE_ADDR"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to convert-emission-service: %v", err)
	}
	defer convertConn.Close()

	emissionClient := pb.NewEmissionClient(emissionConn)
	notifClient := pbnotif.NewNotificationClient(notifConn)
	convertClient := pbconvert.NewConvertClient(convertConn)

	authHandler := handler.NewAuthHandler()
	emissionHandler := handler.NewEmissionHandler(emissionClient)
	prefHandler := handler.NewPreferenceHandler(emissionClient)
	alertHandler := handler.NewAlertHandler(notifClient)
	convertHandler := handler.NewConvertHandler(convertClient)

	e := echo.New()
	router.Setup(e, authHandler, emissionHandler, prefHandler, alertHandler, convertHandler)

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil {
		e.Logger.Error("server stopped", "error", err)
	}
}
