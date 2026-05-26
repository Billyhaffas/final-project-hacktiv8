package main

import (
	"context"
	"log"
	"net/http"
	"notification-service/internal/config"
	handler "notification-service/internal/delivery/http"
	"notification-service/internal/repository"
	"notification-service/internal/usecase"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "notification-service/docs"
)

// @title           Notification API
// @version         1.0
// @description     This is a microservice responsible for validating user carbon emission thresholds and managing notifications.
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
	postgreDB, err := config.ConnectPostgre()
	if err != nil {
		panic(err)
	}
	mongoDBCol := config.ConnectMongo()

	// Inject dependencies
	userPreferenceRepo := repository.NewPreferenceRepository(postgreDB)
	masterRepo := repository.NewMasterLimitRepository(mongoDBCol)
	countEmissionServiceBaseURL := os.Getenv("COUNT_EMISSION_SERVICE_URL")
	if countEmissionServiceBaseURL == "" {
		countEmissionServiceBaseURL = "http://localhost:8081"
	}
	countEmissionClient := repository.NewEmissionClient(countEmissionServiceBaseURL)

	notificationUsecase := usecase.NewNotificationUsecase(userPreferenceRepo, masterRepo, countEmissionClient)
	userPreferenceUsecase := usecase.NewPreferenceUsecase(userPreferenceRepo)

	userPreferenceHandler := handler.NewPreferenceHandler(userPreferenceUsecase)

	// Setup HTTP server for health check and documentation
	go func() {
		// Health Check
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "Ok"}`))
		})

		// Swagger UI
		http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

		log.Println("notification service: HTTP server (health check and swagger) started on :8082")
		log.Fatal(http.ListenAndServe(":8082", nil))
	}()

	// Initialize Echo
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
	e.GET("/api/v1/preferences", userPreferenceHandler.Get)
	e.POST("/api/v1/preferences", userPreferenceHandler.Save)
	e.DELETE("/api/v1/preferences", userPreferenceHandler.Delete)

	// Start the Scheduler asynchronous engine
	log.Println("notification-service: background scheduler engine active...")
	gocron.Every(1).Minutes().Do(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in Emission Cron Job: %v\n", r)
			}
		}()

		log.Println("Cron Job: Starting global emission check...")

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// 1. Fetch all user IDs who have configured preferences
		userIDs, err := userPreferenceRepo.GetAllUserIDs(ctx)
		if err != nil {
			log.Printf("Cron Job Error: Failed to fetch users: %v\n", err)
			return
		}

		// 2. Loop and process each user via the Usecase
		for _, userID := range userIDs {
			breached, message, err := notificationUsecase.CheckAndSendNotification(ctx, userID)
			if err != nil {
				log.Printf("Cron Job Error for User %d: %v\n", userID, err)
				continue
			}

			if breached {
				log.Printf("Cron Job Alert [User %d]: %s\n", userID, message)
				// Your internal dispatch logic (e.g. email, push notification) runs inside the usecase
			}
		}
		log.Println("Cron Job: Global emission check finished successfully.")
	})
	gocron.Start()

	// Start echo server
	log.Printf("notification service: HTTP server starting on :%s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
