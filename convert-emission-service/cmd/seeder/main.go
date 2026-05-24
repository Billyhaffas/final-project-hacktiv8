package main

import (
	"context"
	"log"
	"time"

	"convert-emission-service/config"
	"convert-emission-service/internal/repository/mongodb"
	"convert-emission-service/internal/repository/remote"
	"convert-emission-service/internal/usecase"

	"github.com/joho/godotenv"
)

const csvURL = "https://ourworldindata.org/grapher/weighted-carbon-price-ets.csv?v=1&csvType=full&useColumnShortNames=false"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file detected, fallback configurations applied")
	}

	dbCol := config.ConnectMongo()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	repo := mongodb.NewMongoCarbonPricesRepository(dbCol)
	provider := remote.NewCSVPriceProvider(csvURL)
	seederUseCase := usecase.NewSeedCarbonPricesUseCase(provider, repo)

	if err := seederUseCase.Execute(ctx); err != nil {
		log.Fatalf("convert-emission-service: seeding failed: %v", err)
	}
}
