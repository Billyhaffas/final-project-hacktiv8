// @title          Climate Action API
// @version        1.0
// @description    Personal Carbon Emission Tracker — log commutes, monitor limits, get monthly reports.
// @host           165.245.178.118:8080
// @BasePath       /api/v1
// @schemes        http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Enter your JWT as: Bearer <token>
package main

import (
	"api-gateway/internal/handler"
	"api-gateway/internal/router"
	pbconvert "api-gateway/proto/convert"
	pb "api-gateway/proto/emission"
	pbnotif "api-gateway/proto/notification"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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

	emissionConn, err := grpc.NewClient(
		os.Getenv("COUNT_EMISSION_SERVICE_ADDR"),
		grpcCreds(),
	)
	if err != nil {
		log.Fatalf("failed to connect to count-emission-service: %v", err)
	}
	defer emissionConn.Close()

	notifConn, err := grpc.NewClient(
		os.Getenv("NOTIFICATION_SERVICE_ADDR"),
		grpcCreds(),
	)
	if err != nil {
		log.Fatalf("failed to connect to notification-service: %v", err)
	}
	defer notifConn.Close()

	convertConn, err := grpc.NewClient(
		os.Getenv("CONVERT_EMISSION_SERVICE_ADDR"),
		grpcCreds(),
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

	port := os.Getenv("PORT") // Heroku injects $PORT
	if port == "" {
		port = os.Getenv("GATEWAY_PORT")
	}
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil {
		e.Logger.Error("server stopped", "error", err)
	}
}
