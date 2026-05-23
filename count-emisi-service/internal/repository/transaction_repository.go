package repository

import (
	"context"
	"p3-lc01-billyhaffas/internal/domain"
	"p3-lc01-billyhaffas/internal/model/transaction"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	Collection *mongo.Collection
}

func NewTransactionCollection(collectionRepo *mongo.Collection) domain.TransactionRepository {
	return &TransactionRepository{Collection: collectionRepo}
}

func (cp *TransactionRepository) CreateTransaction(ctx context.Context, user *transaction.TransactionInsertModel) error {
	_, err := cp.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (cp *TransactionRepository) DeleteTransaction(ctx context.Context) error {
	_, err := cp.Collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}
