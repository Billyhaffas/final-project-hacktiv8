package helper

import (
	"count-emission-service/internal/model/emission"
	"count-emission-service/internal/model/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type Respon struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type GetDailyUserEmissionTypeRespon struct {
	Status  string                      `json:"status"`
	Message string                      `json:"message"`
	Data    *emission.UserDailyEmission `json:"data"`
}

type GetMonthlyUserEmissionTypeRespon struct {
	Status  string                        `json:"status"`
	Message string                        `json:"message"`
	Data    *emission.UserMonthlyEmission `json:"data"`
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
