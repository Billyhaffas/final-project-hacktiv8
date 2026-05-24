package main

import (
	"context"
	"count-emission-service/config"
	"count-emission-service/internal/repository/mongodb"
	"count-emission-service/internal/repository/remote"
	"count-emission-service/internal/usecase"
	"log"
	"time"

	"github.com/joho/godotenv"
)

const csvURL = "https://ourworldindata.org/grapher/co-emissions-per-capita.csv?v=1&csvType=full&useColumnShortNames=false"

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env variables")
	}

	// Connect to MongoDB
	dbCol := config.ConnectMongo()

	// Set up a timeout context for the operations
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// Inject dependencies
	repo := mongodb.NewMongoEmissionRepository(dbCol)
	provider := remote.NewCSVProvider(csvURL)
	seederUseCase := usecase.NewSeedEmissionUseCase(provider, repo)

	// Run the data seeding pipeline
	if err := seederUseCase.Execute(ctx); err != nil {
		log.Fatalf("count-emission-service: seeding failed: %v", err)
	}
}
