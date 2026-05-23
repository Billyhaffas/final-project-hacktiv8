package repository

import (
	"context"
	"p3-lc01-billyhaffas/internal/domain"
	"p3-lc01-billyhaffas/internal/model/room"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository struct {
	Collection *mongo.Collection
}

func NewRoomCollection(collectionRepo *mongo.Collection) domain.RoomRepository {
	return &RoomRepository{Collection: collectionRepo}
}

func (cp *RoomRepository) GetRoomById(ctx context.Context, roomId primitive.ObjectID) (*room.RoomSelectModel, error) {
	var responRoom room.RoomSelectModel
	err := cp.Collection.FindOne(ctx, bson.M{
		"_id": roomId,
	}).Decode(&responRoom)
	if err != nil {
		return nil, err
	}
	return &responRoom, nil
}
