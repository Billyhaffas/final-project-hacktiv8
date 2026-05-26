package main

import (
	"log"
	"net/http"
	"os"

	"convert-emission-service/config"
	handler "convert-emission-service/internal/delivery/http"
	"convert-emission-service/internal/repository"
	"convert-emission-service/internal/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "convert-emission-service/docs"
)

// @title           Convert Emission API
// @version         1.0
// @description     This is a microservice responsible for converting carbon emission into local currency.
func main() {
	// Set environment
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, relying on system environment variables")
	}

	if os.Getenv("JWT_SECRET_KEY") == "" {
		log.Fatal("Critical: JWT_SECRET_KEY is not set in the environment.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize databases
	mongoDBCol := config.ConnectMongo()

	// Inject dependencies
	priceRepo := repository.NewCarbonPriceRepository(mongoDBCol)
	convUsecase := usecase.NewConversionUsecase(priceRepo)
	convHandler := handler.NewConversionHandler(convUsecase)

	// Setup HTTP server for health check and documentation
	go func() {
		// Health Check
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "Ok"}`))
		})

		// Swagger UI
		http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

		log.Println("convert-emission service: HTTP server (health check and swagger) started on :8082")
		log.Fatal(http.ListenAndServe(":8082", nil))
	}()

	// Initialize echo
	e := echo.New()

	// Setup Middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("%s %s → %d (%s)", v.Method, v.URI, v.Status, v.Latency)
			return nil
		},
	}))
	e.Use(middleware.Recover())

	// Register routes
	e.POST("/api/v1/convert/daily", convHandler.HandleDaily)
	e.POST("/api/v1/convert/monthly", convHandler.HandleMonthly)
	e.POST("/api/v1/convert/yearly", convHandler.HandleYearly)

	// Start echo server
	log.Printf("convert-emission service: HTTP server starting on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
