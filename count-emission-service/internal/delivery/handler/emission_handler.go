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

func (uh *EmissionHandler) GetUserDailyEmission(c *echo.Context) error {
	var userId int32
	userId = 1
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, helper.Respon{
	// 		Status:  http.StatusText(http.StatusBadRequest),
	// 		Message: "Invalid user ID",
	// 	})
	// }
	userEmission, err := uh.EmissionUseCase.GetUserDailyEmission(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, helper.GetDailyUserEmissionTypeRespon{
		Status:  http.StatusText(http.StatusOK),
		Message: "User daily emission retrieved successfully",
		Data:    userEmission,
	})
}

func (uh *EmissionHandler) GetUserMonthlyEmission(c *echo.Context) error {
	var userId int32
	userId = 1
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, helper.Respon{
	// 		Status:  http.StatusText(http.StatusBadRequest),
	// 		Message: "Invalid user ID",
	// 	})
	// }
	userMonthlyEmission, err := uh.EmissionUseCase.GetUserMonthlyEmission(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, helper.GetMonthlyUserEmissionTypeRespon{
		Status:  http.StatusText(http.StatusOK),
		Message: "User monthly emission retrieved successfully",
		Data:    userMonthlyEmission,
	})
}

func (uh *EmissionHandler) GetUserYearlyEmission(c *echo.Context) error {
	var userId int32
	userId = 1
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, helper.Respon{
	// 		Status:  http.StatusText(http.StatusBadRequest),
	// 		Message: "Invalid user ID",
	// 	})
	// }
	userYearlyEmission, err := uh.EmissionUseCase.GetUserYearlyEmission(c.Request().Context(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Respon{
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, helper.GetYearlyUserEmissionTypeRespon{
		Status:  http.StatusText(http.StatusOK),
		Message: "User yearly emission retrieved successfully",
		Data:    userYearlyEmission,
	})
}
