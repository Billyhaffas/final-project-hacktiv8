package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/internal/handler"
	pbconvert "api-gateway/proto/convert"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
)

// --- mock convert client ---

type mockConvertClient struct {
	convertToIDR func(context.Context, *pbconvert.ConvertToIDRRequest, ...grpc.CallOption) (*pbconvert.ConvertToIDRResponse, error)
}

func (m *mockConvertClient) ConvertToIDR(ctx context.Context, req *pbconvert.ConvertToIDRRequest, opts ...grpc.CallOption) (*pbconvert.ConvertToIDRResponse, error) {
	if m.convertToIDR != nil {
		return m.convertToIDR(ctx, req, opts...)
	}
	kg := req.EmissionKgCo2
	return &pbconvert.ConvertToIDRResponse{
		PricePerTonUsd:     23.0,
		ExchangeRateUsdIdr: 16250.0,
		TotalIdr:           (kg / 1000.0) * 23.0 * 16250.0,
	}, nil
}

// --- tests ---

func TestConvertToIDR_Handler(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		query      string
		client     *mockConvertClient
		wantStatus int
	}{
		{
			name:       "200 — valid kg",
			query:      "?kg=1000",
			client:     &mockConvertClient{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "200 — zero kg",
			query:      "?kg=0",
			client:     &mockConvertClient{},
			wantStatus: http.StatusOK,
		},
		{
			name:       "400 — missing kg param",
			query:      "",
			client:     &mockConvertClient{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "400 — negative kg",
			query:      "?kg=-5",
			client:     &mockConvertClient{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "400 — non-numeric kg",
			query:      "?kg=abc",
			client:     &mockConvertClient{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:  "500 — gRPC error",
			query: "?kg=100",
			client: &mockConvertClient{
				convertToIDR: func(_ context.Context, _ *pbconvert.ConvertToIDRRequest, _ ...grpc.CallOption) (*pbconvert.ConvertToIDRResponse, error) {
					return nil, errors.New("convert service unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewConvertHandler(tt.client)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/emissions/convert"+tt.query, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", 1)

			if err := h.ConvertToIDR(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}
