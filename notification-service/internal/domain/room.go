package domain

import (
	"context"
	"p3-lc01-billyhaffas/internal/model/room"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomRepository interface {
	GetRoomById(ctx context.Context, roomId primitive.ObjectID) (*room.RoomSelectModel, error)
}
