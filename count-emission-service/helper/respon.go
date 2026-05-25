package helper

import (
	"count-emission-service/internal/model/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type Respon struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetUserTypeRespon struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Data    *user.UserSelectModel `json:"data"`
}

type GetAllUserTypeRespon struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []user.UserSelectModel `json:"data"`
}

type GetUserId struct {
	Id string `json:"id"`
}

type DeleteUserTypeRespon struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    *mongo.DeleteResult `json:"data"`
}
