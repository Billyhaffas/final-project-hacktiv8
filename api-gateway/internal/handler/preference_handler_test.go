package handler_test

import (
	"api-gateway/internal/handler"
	pb "api-gateway/proto/emission"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
)

func TestGetPreferences(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		client     *mockEmissionClient
		wantStatus int
	}{
		{
			name:       "200 — success",
			client:     &mockEmissionClient{},
			wantStatus: http.StatusOK,
		},
		{
			name: "500 — gRPC error",
			client: &mockEmissionClient{
				getUserPreferences: func(_ context.Context, _ *pb.Empty, _ ...grpc.CallOption) (*pb.UserPreferences, error) {
					return nil, errors.New("gRPC unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewPreferenceHandler(tt.client)
			c, rec := newCtxWithUser(e, http.MethodGet, "", 1)
			if err := h.GetPreferences(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestUpdatePreferences(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		body       string
		client     *mockEmissionClient
		wantStatus int
	}{
		{
			name:       "200 — success",
			body:       `{"country_code":"IDN","custom_daily_limit_kg_co2":5.0}`,
			client:     &mockEmissionClient{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "400 — malformed JSON",
			body:       `{invalid`,
			client:     &mockEmissionClient{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "500 — gRPC error",
			body: `{"country_code":"IDN"}`,
			client: &mockEmissionClient{
				setUserPreferences: func(_ context.Context, _ *pb.SetUserPreferencesBody, _ ...grpc.CallOption) (*pb.UserPreferences, error) {
					return nil, errors.New("gRPC unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewPreferenceHandler(tt.client)
			c, rec := newCtxWithUser(e, http.MethodPut, tt.body, 1)
			if err := h.UpdatePreferences(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}
