package handler

import (
	"net/http"
	"strconv"
	"time"

	"api-gateway/helper"
	pbnotif "api-gateway/proto/notification"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc/metadata"
)

type AlertHandler struct {
	client pbnotif.NotificationClient
}

func NewAlertHandler(client pbnotif.NotificationClient) *AlertHandler {
	return &AlertHandler{client: client}
}

func (h *AlertHandler) CheckAlert(c *echo.Context) error {
	userID := c.Get("user_id").(int)
	date := time.Now().Format("2006-01-02")

	md := metadata.Pairs("user-id", strconv.Itoa(userID))
	ctx := metadata.NewOutgoingContext(c.Request().Context(), md)

	resp, err := h.client.CheckDailyAlert(ctx, &pbnotif.DailyAlertRequest{
		UserId: int32(userID),
		Date:   date,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.Fail("INTERNAL_ERROR", "failed to check daily alert"))
	}

	return c.JSON(http.StatusOK, helper.Success(map[string]interface{}{
		"is_exceeded":      resp.IsExceeded,
		"daily_total_kg":   resp.DailyTotalKg,
		"daily_limit_kg":   resp.DailyLimitKg,
		"threshold_source": resp.ThresholdSource,
		"message":          resp.Message,
	}))
}
