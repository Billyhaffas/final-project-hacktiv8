package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {

	// if err := godotenv.Load(); err != nil {
	// 	log.Println(".env file not found, using system environment")
	// }

	// uri := os.Getenv("MONGO_URI")
	// dbName := os.Getenv("DB_NAME")
	fmt.Println("MONGO_URI =", os.Getenv("MONGO_URI"))
	fmt.Println("DB_NAME =", os.Getenv("DB_NAME"))
	// uri := os.Getenv("mongodb://mongo:27017")
	// dbName := os.Getenv("learn-microservice-hacktiv8")

	// if uri == "" {
	// 	log.Fatal("MONGO_URI is empty")
	// }

	// if dbName == "" {
	// 	log.Fatal("DB_NAME is empty")
	// }

	// fmt.Println("Mongo URI:", uri)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	opts := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = client.Database("lc01-billyhaffas-p3")

	log.Println("mongodb connected")
}
