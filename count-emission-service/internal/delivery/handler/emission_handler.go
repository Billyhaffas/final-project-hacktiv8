package handler

import (
	"count-emission-service/helper"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/emission"
	"net/http"

	"github.com/labstack/echo/v5"
)

type EmissionHandler struct {
	EmissionUseCase domain.EmissionUseCase
}

func NewEmissionHandler(EmissionUC domain.EmissionUseCase) domain.EmissionHandler {
	return &EmissionHandler{EmissionUseCase: EmissionUC}
}

func (uh *EmissionHandler) CreateUserEmission(c *echo.Context) error {
	var requestEmission *emission.EmissionBody
	err := c.Bind(&requestEmission)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.Respon{
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
	}
	err = uh.EmissionUseCase.CreateUserEmission(c.Request().Context(), requestEmission)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, helper.Respon{
		Status:  http.StatusText(http.StatusCreated),
		Message: "Emission has been created",
	})
}
