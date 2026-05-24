package mongodb

import (
	"context"
	"convert-emission-service/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoCarbonPricesRepository struct {
	collection *mongo.Collection
}

// Renamed from NewMongoEmissionRepository
func NewMongoCarbonPricesRepository(collection *mongo.Collection) domain.CarbonPriceRepository {
	return &MongoCarbonPricesRepository{
		collection: collection,
	}
}

func (r *MongoCarbonPricesRepository) BulkInsert(ctx context.Context, carbonPrices []domain.CarbonPrice) (bool, error) {
	// Check if records already exist
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil // Data exists, safely skip insertion
	}

	if len(carbonPrices) == 0 {
		return false, nil
	}

	// Prepare documents for bulk write
	var documents []interface{}
	for _, cp := range carbonPrices {
		documents = append(documents, cp)
	}

	// Perform the insertion
	_, err = r.collection.InsertMany(ctx, documents)
	if err != nil {
		return false, err
	}

	return true, nil // Successfully seeded
}
