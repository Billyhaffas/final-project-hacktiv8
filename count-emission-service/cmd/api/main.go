package main

import (
	"count-emission-service/internal/delivery/handler"
	"count-emission-service/internal/infrastructure/database"
	externalapi "count-emission-service/internal/repository/external_api"
	repository "count-emission-service/internal/repository/internal_repo"
	"count-emission-service/internal/usecase"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

func main() {
	db, err := database.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	emissionRepo := repository.NewEmissionCollection(db)
	carbonSutraRepo := externalapi.NewCarbonSutraRepository(httpClient)

	emissionUseCase := usecase.NewEmissionUseCase(emissionRepo, carbonSutraRepo)

	emissionHandler := handler.NewEmissionHandler(emissionUseCase)

	e := echo.New()
	e.POST("/emissions", emissionHandler.CreateUserEmission)

	e.Start(":8080")
}
