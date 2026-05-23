package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

	col := client.Database("carbon_prices_db").Collection("carbon_prices")

	idx := mongo.IndexModel{
		Keys:    bson.D{{Key: "Id", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("idx_carbon_prices_id"),
	}

	if _, err := col.Indexes().CreateOne(ctx, idx); err != nil {
		log.Println("convert-emission-service: index warning:", err)
	}

	log.Println("convert-emission-service: connected to carbon_prices_db")
	return col
}
