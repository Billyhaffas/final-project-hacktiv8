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

	// Configure client options
	opts := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPI).
		SetTimeout(10 * time.Second)

	// Establish connection
	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal("notification-service: mongo connect error:", err)
	}

	// Use a fresh context for temporary lifecycle operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify the connection is alive
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("notification-service: mongo ping error:", err)
	}

	col := client.Database("climate_action_masterdata").Collection("emission_thresholds")

	log.Println("notification-service: connected to mongo database: climate_action_masterdata (collection: emission_thresholds)")
	return col
}
