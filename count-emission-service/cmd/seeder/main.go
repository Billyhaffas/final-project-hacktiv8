package main

import (
	"context"
	"count-emission-service/config"
	"count-emission-service/internal/repository/mongodb"
	"count-emission-service/internal/repository/remote"
	"count-emission-service/internal/usecase"
	"log"
	"time"
)

const csvURL = "https://ourworldindata.org/grapher/co-emissions-per-capita.csv?v=1&csvType=full&useColumnShortNames=false"

func main() {
	// 1. Fetch pre-configured collection directly
	dbCol := config.ConnectMongo()

	// 2. Set up a timeout context for the operations
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// 3. Inject dependencies using Clean Architecture patterns
	repo := mongodb.NewMongoEmissionRepository(dbCol)
	provider := remote.NewCSVProvider(csvURL)
	seederUseCase := usecase.NewSeedEmissionUseCase(provider, repo)

	// 4. Run the data seeding pipeline
	if err := seederUseCase.Execute(ctx); err != nil {
		log.Fatalf("count-emission-service: seeding failed: %v", err)
	}
}
