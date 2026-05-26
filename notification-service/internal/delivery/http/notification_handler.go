package handler

import (
	"net/http"
	"notification-service/internal/domain"
	"notification-service/internal/helper"

	"github.com/labstack/echo/v5"
)

type notificationHandler struct {
	usecase domain.NotificationUsecase
}

func NewNotificationHandler(u domain.NotificationUsecase) *notificationHandler {
	return &notificationHandler{usecase: u}
}

// TriggerCheck godoc
// @Summary      Check and trigger emission notifications
// @Description  Fetches recent emission counts, compares them against user preferences or MongoDB master fallbacks, and triggers a warning if breached.
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  true  "Insert Bearer JWT_TOKEN here"
// @Success      200            {object}  map[string]interface{} "Example: {'status':'success','limit_breached':true,'message':'Alert!...'}"
// @Failure      401            {object}  map[string]string      "Example: {'error':'missing or invalid token'}"
// @Failure      500            {object}  map[string]string      "Example: {'error':'internal database failure'}"
// @Router       /api/v1/notifications/check-emission [post]
func (h *notificationHandler) TriggerCheck(e *echo.Context) error {
	// 1. Extract user_id securely from the JWT token claims
	userID, err := helper.GetUserIDFromJWT(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	// 2. Safely fetch the request's context.Context
	ctx := e.Request().Context()

	// 3. Execute business logic using the secure identity
	breached, message, err := h.usecase.CheckAndSendNotification(ctx, userID)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return e.JSON(http.StatusOK, map[string]interface{}{
		"status":         "success",
		"limit_breached": breached,
		"message":        message,
	})
}
