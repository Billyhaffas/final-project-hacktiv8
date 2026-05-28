package room

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomInsertModel struct {
	Name     string  `bson:"Nama" json:"Nama"`
	TypeRoom string  `bson:"Tipe" json:"Tipe"`
	Price    float32 `bson:"Harga" json:"Harga"`
	Discount string  `bson:"Discount" json:"Discount"`
}

type RoomSelectModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"Nama" json:"Nama"`
	TypeRoom string             `bson:"Tipe" json:"Tipe"`
	Price    float32            `bson:"Harga" json:"Harga"`
	Discount bool               `bson:"Discount" json:"Discount"`
}
