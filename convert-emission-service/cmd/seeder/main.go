package main

import (
	"context"
	"log"
	"time"

	"convert-emission-service/config"
	"convert-emission-service/internal/repository/mongodb"
	"convert-emission-service/internal/repository/remote"
	"convert-emission-service/internal/usecase"
)

const csvURL = "https://ourworldindata.org/grapher/weighted-carbon-price-ets.csv?v=1&csvType=full&useColumnShortNames=false"

func main() {
	dbCol := config.ConnectMongo()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	repo := mongodb.NewMongoEmissionRepository(dbCol)
	provider := remote.NewCSVProvider(csvURL)
	seederUseCase := usecase.NewSeedEmissionUseCase(provider, repo)

	if err := seederUseCase.Execute(ctx); err != nil {
		log.Fatalf("convert-emission-service: seeding failed: %v", err)
	}
}
