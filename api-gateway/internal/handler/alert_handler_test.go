package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/internal/handler"
	pbnotif "api-gateway/proto/notification"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
)

// --- mock notification client ---

type mockNotificationClient struct {
	checkDailyAlert func(context.Context, *pbnotif.DailyAlertRequest, ...grpc.CallOption) (*pbnotif.DailyAlertResponse, error)
}

func (m *mockNotificationClient) CheckDailyAlert(ctx context.Context, req *pbnotif.DailyAlertRequest, opts ...grpc.CallOption) (*pbnotif.DailyAlertResponse, error) {
	if m.checkDailyAlert != nil {
		return m.checkDailyAlert(ctx, req, opts...)
	}
	return &pbnotif.DailyAlertResponse{
		IsExceeded:      false,
		DailyTotalKg:    3.0,
		DailyLimitKg:    6.3,
		ThresholdSource: "country",
		Message:         "Your emission is within safe limits.",
	}, nil
}

// --- tests ---

func TestCheckAlert(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		client     *mockNotificationClient
		wantStatus int
	}{
		{
			name:       "200 — not exceeded",
			client:     &mockNotificationClient{},
			wantStatus: http.StatusOK,
		},
		{
			name: "200 — exceeded",
			client: &mockNotificationClient{
				checkDailyAlert: func(_ context.Context, _ *pbnotif.DailyAlertRequest, _ ...grpc.CallOption) (*pbnotif.DailyAlertResponse, error) {
					return &pbnotif.DailyAlertResponse{
						IsExceeded:      true,
						DailyTotalKg:    8.5,
						DailyLimitKg:    6.3,
						ThresholdSource: "country",
						Message:         "Alert! ...",
					}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "500 — gRPC error",
			client: &mockNotificationClient{
				checkDailyAlert: func(_ context.Context, _ *pbnotif.DailyAlertRequest, _ ...grpc.CallOption) (*pbnotif.DailyAlertResponse, error) {
					return nil, errors.New("notification service unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewAlertHandler(tt.client)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/emissions/alert", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", 1)

			if err := h.CheckAlert(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}
