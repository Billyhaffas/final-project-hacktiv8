package mongodb

import (
	"context"
	"count-emission-service/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoEmissionRepository struct {
	collection *mongo.Collection
}

func NewMongoEmissionRepository(collection *mongo.Collection) domain.EmissionRepository {
	return &MongoEmissionRepository{
		collection: collection,
	}
}

func (r *MongoEmissionRepository) BulkInsert(ctx context.Context, emissions []domain.Emission) (bool, error) {
	// Check if records already exist
	count, err := r.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil // Data exists, safely skip insertion
	}

	if len(emissions) == 0 {
		return false, nil
	}

	// Prepare documents for bulk write
	var documents []interface{}
	for _, e := range emissions {
		documents = append(documents, e)
	}

	// Perform the insertion
	_, err = r.collection.InsertMany(ctx, documents)
	if err != nil {
		return false, err
	}

	return true, nil // Successfully seeded
}
