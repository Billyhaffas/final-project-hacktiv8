package repository

import (
	"context"
	"notification-service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type masterLimitRepo struct {
	collection *mongo.Collection
}

func NewMasterLimitRepository(col *mongo.Collection) domain.MasterLimitRepository {
	return &masterLimitRepo{collection: col}
}

// MongoEmissionThreshold represents your exact MongoDB schema match
type MongoEmissionThreshold struct {
	ID              string    `bson:"_id"`
	CountryCode     string    `bson:"country_code"`
	DailyLimitKgCo2 float64   `bson:"daily_limit_kg_co2"`
	SourceURL       string    `bson:"source_url"`
	UpdatedAt       time.Time `bson:"updated_at"`
}

func (r *masterLimitRepo) GetDefaultLimitByCountry(ctx context.Context, countryCode string) (float64, error) {
	var result MongoEmissionThreshold

	// Match the 3-letter country code
	filter := bson.M{"country_code": countryCode}

	// Sort by updated_at descending (-1) or _id descending (-1)
	// So that "IDN-2024" comes before "IDN-2023" if timestamps match
	opts := options.FindOne().SetSort(bson.D{
		{Key: "updated_at", Value: -1},
		{Key: "_id", Value: -1},
	})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.DailyLimitKgCo2, nil
}
