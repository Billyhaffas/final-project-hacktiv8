package handler

import (
	"net/http"
	"notification-service/internal/domain"
	"notification-service/internal/helper"

	"github.com/labstack/echo/v5"
)

type PreferenceHandler struct {
	usecase domain.PreferenceUsecase
}

func NewPreferenceHandler(u domain.PreferenceUsecase) *PreferenceHandler {
	return &PreferenceHandler{usecase: u}
}

// Get godoc
// @Summary      Get user preference configurations
// @Tags         preferences
// @Param        Authorization  header    string  true  "Bearer <JWT_TOKEN>"
// @Success      200            {object}  domain.UserEmissionPreference
// @Router       /api/v1/preferences [get]
func (h *PreferenceHandler) Get(e *echo.Context) error {
	userID, err := helper.GetUserIDFromJWT(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	pref, err := h.usecase.GetPreference(e.Request().Context(), userID)
	if err != nil {
		return e.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return e.JSON(http.StatusOK, pref)
}

// Save godoc
// @Summary      Create or Update user preference profiles
// @Tags         preferences
// @Param        Authorization  header    string  true  "Bearer <JWT_TOKEN>"
// @Param        body           body      domain.PreferenceUpsertInput true "Preference Data Payload"
// @Success      200            {object}  domain.UserEmissionPreference
// @Router       /api/v1/preferences [post]
func (h *PreferenceHandler) Save(e *echo.Context) error {
	userID, err := helper.GetUserIDFromJWT(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	var input domain.PreferenceUpsertInput
	if err := e.Bind(&input); err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body format"})
	}
	input.UserID = userID // Force authenticated token user ID onto payload

	pref, err := h.usecase.SavePreference(e.Request().Context(), input)
	if err != nil {
		return e.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	return e.JSON(http.StatusOK, pref)
}

// Delete godoc
// @Summary      Reset and drop custom user preference configs
// @Tags         preferences
// @Param        Authorization  header    string  true  "Bearer <JWT_TOKEN>"
// @Success      200            {object}  map[string]string
// @Router       /api/v1/preferences [delete]
func (h *PreferenceHandler) Delete(e *echo.Context) error {
	userID, err := helper.GetUserIDFromJWT(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	err = h.usecase.DeletePreference(e.Request().Context(), userID)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return e.JSON(http.StatusOK, map[string]string{"message": "preferences dropped successfully, system will default back to regional master limits"})
}
