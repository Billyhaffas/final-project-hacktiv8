package handler

import (
	"api-gateway/helper"
	pb "api-gateway/proto/emission"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/metadata"
)

type EmissionHandler struct {
	client pb.EmissionClient
}

func NewEmissionHandler(client pb.EmissionClient) *EmissionHandler {
	return &EmissionHandler{client: client}
}

func (h *EmissionHandler) LogEmission(c *echo.Context) error {
	userID := c.Get("user_id").(int)

	var req struct {
		VehicleType string  `json:"vehicle_type"`
		FuelType    string  `json:"fuel_type"`
		DistanceKm  float64 `json:"distance_km"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.Fail("INVALID_REQUEST", "invalid request body"))
	}
	if req.VehicleType == "" || req.DistanceKm <= 0 {
		return c.JSON(http.StatusBadRequest, helper.Fail("VALIDATION_ERROR", "vehicle_type and distance_km are required"))
	}

	md := metadata.Pairs("user-id", strconv.Itoa(userID))
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

	resp, err := h.client.CreateUserEmission(ctx, &pb.EmissionBody{
		UserId:      int32(userID),
		VehicleType: req.VehicleType,
		FuelType:    req.FuelType,
		DistanceKm:  req.DistanceKm,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "failed to log emission"))
	}

	return c.JSON(http.StatusCreated, helper.Success(map[string]string{"message": resp.Message}))
}

func (h *EmissionHandler) GetDailyTotal(c *echo.Context) error {
	userID := c.Get("user_id").(int)

	md := metadata.Pairs("user-id", strconv.Itoa(userID))
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

	resp, err := h.client.GetUserDailyEmission(ctx, &pb.Empty{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "failed to get daily emission"))
	}

	return c.JSON(http.StatusOK, helper.Success(map[string]interface{}{
		"user_id":               resp.UserId,
		"date":                  resp.Date,
		"total_emission_kg_co2": resp.TotalEmissionKgCo2,
	}))
}

func (h *EmissionHandler) GetMonthlyReport(c *echo.Context) error {
	userID := c.Get("user_id").(int)

	md := metadata.Pairs("user-id", strconv.Itoa(userID))
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

	resp, err := h.client.GetUserMonthlyEmission(ctx, &pb.Empty{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "failed to get monthly report"))
	}

	return c.JSON(http.StatusOK, helper.Success(resp))
}

func (h *EmissionHandler) GetAlert(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, helper.Fail("NOT_IMPLEMENTED", "notification service not ready"))
}

func (h *EmissionHandler) ConvertToIDR(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, helper.Fail("NOT_IMPLEMENTED", "convert service not ready"))
}
