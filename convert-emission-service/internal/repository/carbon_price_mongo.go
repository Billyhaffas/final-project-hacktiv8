package repository

import (
	"context"
	"convert-emission-service/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CarbonPriceDBModel struct {
	ID             string    `bson:"_id"`
	PricePerTonUsd float64   `bson:"price_per_ton_usd"`
	UsdCurRate     float64   `bson:"usd_cur_rate"`
	Source         string    `bson:"source"`
	FetchedAt      time.Time `bson:"fetched_at"`
}

type carbonPriceRepository struct {
	collection *mongo.Collection
}

func NewCarbonPriceRepository(col *mongo.Collection) domain.CarbonPriceRepository {
	return &carbonPriceRepository{collection: col}
}

func (r *carbonPriceRepository) GetLatestRateByCountryCode(ctx context.Context, countryCode string) (*domain.CarbonPrice, error) {
	var result CarbonPriceDBModel

	// Match IDs that start with the 3-letter country code (e.g. "IDN-2025")
	filter := bson.M{"_id": bson.M{"$regex": "^" + countryCode}}

	// Tie-breaker sort order to always retrieve the newest data entry
	opts := options.FindOne().SetSort(bson.D{
		{Key: "fetched_at", Value: -1},
		{Key: "_id", Value: -1},
	})

	err := r.collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &domain.CarbonPrice{
		ID:             result.ID,
		PricePerTonUsd: result.PricePerTonUsd,
		UsdCurRate:     result.UsdCurRate,
		Source:         result.Source,
		FetchedAt:      result.FetchedAt,
	}, nil
}
