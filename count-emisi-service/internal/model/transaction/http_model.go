package transaction

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionBodyPayload struct {
	UserId    string    `bson:"user_id" json:"user_id"`
	RoomId    string    `bson:"room_id" json:"room_id"`
	Qty       int       `bson:"quantity" json:"quantity"`
	OrderDate time.Time `bson:"order_date" json:"order_date"`
}

type TransactionInsertModel struct {
	UserId    string    `bson:"user_id" json:"user_id"`
	RoomId    string    `bson:"room_id" json:"room_id"`
	Qty       int       `bson:"quantity" json:"quantity"`
	OrderDate time.Time `bson:"order_date" json:"order_date"`
	Subtotal  float32   `bson:"subtotal" json:"subtotal"`
	Total     float32   `bson:"total" json:"total"`
}

type UserSelectModel struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Age     int                `bson:"age" json:"age"`
	Address string             `bson:"address" json:"address"`
	Phone   string             `bson:"phone" json:"phone"`
	Email   string             `bson:"email" json:"email"`
}

type UserUpdateModel struct {
	Name    string `bson:"name" json:"name"`
	Age     int    `bson:"age" json:"age"`
	Address string `bson:"address" json:"address"`
	Phone   string `bson:"phone" json:"phone"`
	Email   string `bson:"email" json:"email"`
}
