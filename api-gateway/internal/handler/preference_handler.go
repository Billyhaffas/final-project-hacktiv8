package handler

import (
	"api-gateway/helper"
	pb "api-gateway/proto/emission"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/metadata"
)

type PreferenceHandler struct {
	client pb.EmissionClient
}

func NewPreferenceHandler(client pb.EmissionClient) *PreferenceHandler {
	return &PreferenceHandler{client: client}
}

// GetPreferences godoc
// @Summary      Get user emission preferences
// @Tags         preferences
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  SuccessResponse{data=PreferencesData}
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /preferences [get]
func (h *PreferenceHandler) GetPreferences(c *echo.Context) error {
	userID := c.Get("user_id").(int)

	md := metadata.Pairs("user-id", strconv.Itoa(userID))
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

	resp, err := h.client.GetUserPreferences(ctx, &pb.Empty{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "failed to get preferences"))
	}

	return c.JSON(http.StatusOK, helper.Success(map[string]interface{}{
		"user_id":                   resp.UserId,
		"country_code":              resp.CountryCode,
		"custom_daily_limit_kg_co2": resp.CustomDailyLimitKgCo2,
	}))
}

// UpdatePreferences godoc
// @Summary      Update user emission preferences
// @Tags         preferences
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      UpdatePreferencesRequest  true  "Preferences to update"
// @Success      200   {object}  SuccessResponse{data=PreferencesData}
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /preferences [put]
func (h *PreferenceHandler) UpdatePreferences(c *echo.Context) error {
	userID := c.Get("user_id").(int)

	var req struct {
		CountryCode           string  `json:"country_code"`
		CustomDailyLimitKgCo2 float64 `json:"custom_daily_limit_kg_co2"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}

	md := metadata.Pairs("user-id", strconv.Itoa(userID))
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

	resp, err := h.client.SetUserPreferences(ctx, &pb.SetUserPreferencesBody{
		CountryCode:           req.CountryCode,
		CustomDailyLimitKgCo2: req.CustomDailyLimitKgCo2,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "failed to update preferences"))
	}

	return c.JSON(http.StatusOK, helper.Success(map[string]interface{}{
		"user_id":                   resp.UserId,
		"country_code":              resp.CountryCode,
		"custom_daily_limit_kg_co2": resp.CustomDailyLimitKgCo2,
	}))
}
