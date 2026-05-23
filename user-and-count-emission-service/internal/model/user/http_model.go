package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInsertModel struct {
	Name    string `bson:"name" json:"name"`
	Age     int    `bson:"age" json:"age"`
	Address string `bson:"address" json:"address"`
	Phone   string `bson:"phone" json:"phone"`
	Email   string `bson:"email" json:"email"`
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
