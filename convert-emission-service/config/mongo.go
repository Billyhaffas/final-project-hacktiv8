package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongo() *mongo.Collection {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	opts := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPI).
		SetTimeout(10 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("convert-emission-service: mongo connect error:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("convert-emission-service: mongo ping error:", err)
	}

	col := client.Database("climate_action_masterdata").Collection("carbon_prices")

	log.Println("convert-emission-service: connected to climate_action_masterdata")
	return col
}
