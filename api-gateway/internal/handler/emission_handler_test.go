package handler_test

import (
	"api-gateway/internal/handler"
	pb "api-gateway/proto/emission"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"google.golang.org/grpc"
)

// --- mock gRPC client ---

type mockEmissionClient struct {
	createUserEmission     func(context.Context, *pb.EmissionBody, ...grpc.CallOption) (*pb.PostRespon, error)
	getUserDailyEmission   func(context.Context, *pb.Empty, ...grpc.CallOption) (*pb.UserDailyEmission, error)
	getUserMonthlyEmission func(context.Context, *pb.Empty, ...grpc.CallOption) (*pb.UserMonthlyEmission, error)
	getUserYearlyEmission  func(context.Context, *pb.Empty, ...grpc.CallOption) (*pb.UserYearlyEmission, error)
	getUserPreferences     func(context.Context, *pb.Empty, ...grpc.CallOption) (*pb.UserPreferences, error)
	setUserPreferences     func(context.Context, *pb.SetUserPreferencesBody, ...grpc.CallOption) (*pb.UserPreferences, error)
}

func (m *mockEmissionClient) CreateUserEmission(ctx context.Context, in *pb.EmissionBody, opts ...grpc.CallOption) (*pb.PostRespon, error) {
	if m.createUserEmission != nil {
		return m.createUserEmission(ctx, in, opts...)
	}
	return &pb.PostRespon{Message: "Emission has been created"}, nil
}

func (m *mockEmissionClient) GetUserDailyEmission(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.UserDailyEmission, error) {
	if m.getUserDailyEmission != nil {
		return m.getUserDailyEmission(ctx, in, opts...)
	}
	return &pb.UserDailyEmission{UserId: 1, Date: "2026-05-27", TotalEmissionKgCo2: 2.5}, nil
}

func (m *mockEmissionClient) GetUserMonthlyEmission(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.UserMonthlyEmission, error) {
	if m.getUserMonthlyEmission != nil {
		return m.getUserMonthlyEmission(ctx, in, opts...)
	}
	return &pb.UserMonthlyEmission{UserId: 1, TotalEmissionMonthlyKgCo2: 50.0}, nil
}

func (m *mockEmissionClient) GetUserYearlyEmission(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.UserYearlyEmission, error) {
	if m.getUserYearlyEmission != nil {
		return m.getUserYearlyEmission(ctx, in, opts...)
	}
	return &pb.UserYearlyEmission{UserId: 1, TotalEmissionYearlyKgCo2: 500.0}, nil
}

func (m *mockEmissionClient) GetUserPreferences(ctx context.Context, in *pb.Empty, opts ...grpc.CallOption) (*pb.UserPreferences, error) {
	if m.getUserPreferences != nil {
		return m.getUserPreferences(ctx, in, opts...)
	}
	return &pb.UserPreferences{UserId: 1, CountryCode: "IDN", CustomDailyLimitKgCo2: 6.3}, nil
}

func (m *mockEmissionClient) SetUserPreferences(ctx context.Context, in *pb.SetUserPreferencesBody, opts ...grpc.CallOption) (*pb.UserPreferences, error) {
	if m.setUserPreferences != nil {
		return m.setUserPreferences(ctx, in, opts...)
	}
	return &pb.UserPreferences{UserId: 1, CountryCode: in.CountryCode, CustomDailyLimitKgCo2: in.CustomDailyLimitKgCo2}, nil
}

// newCtxWithUser builds an Echo context with user_id already set (simulates JWT middleware).
func newCtxWithUser(e *echo.Echo, method, body string, userID int) (*echo.Context, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", userID)
	return c, rec
}

// --- LogEmission ---

func TestLogEmission(t *testing.T) {
	e := echo.New()

	tests := []struct {
		name       string
		body       string
		client     *mockEmissionClient
		wantStatus int
	}{
		{
			name:       "201 — success",
			body:       `{"vehicle_type":"Car-Size-Medium","fuel_type":"Petrol","distance_km":10.5}`,
			client:     &mockEmissionClient{},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "400 — malformed JSON",
			body:       `{invalid`,
			client:     &mockEmissionClient{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "400 — missing vehicle_type",
			body:       `{"distance_km":10.5}`,
			client:     &mockEmissionClient{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "400 — distance_km is zero",
			body:       `{"vehicle_type":"Car-Size-Medium","distance_km":0}`,
			client:     &mockEmissionClient{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "500 — gRPC error",
			body: `{"vehicle_type":"Car-Size-Medium","distance_km":10.5}`,
			client: &mockEmissionClient{
				createUserEmission: func(_ context.Context, _ *pb.EmissionBody, _ ...grpc.CallOption) (*pb.PostRespon, error) {
					return nil, errors.New("gRPC unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewEmissionHandler(tt.client)
			c, rec := newCtxWithUser(e, http.MethodPost, tt.body, 1)
			if err := h.LogEmission(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

// --- GetDailyTotal ---

func TestGetDailyTotal(t *testing.T) {
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
				getUserDailyEmission: func(_ context.Context, _ *pb.Empty, _ ...grpc.CallOption) (*pb.UserDailyEmission, error) {
					return nil, errors.New("gRPC unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewEmissionHandler(tt.client)
			c, rec := newCtxWithUser(e, http.MethodGet, "", 1)
			if err := h.GetDailyTotal(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

// --- GetMonthlyReport ---

func TestGetMonthlyReport(t *testing.T) {
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
				getUserMonthlyEmission: func(_ context.Context, _ *pb.Empty, _ ...grpc.CallOption) (*pb.UserMonthlyEmission, error) {
					return nil, errors.New("gRPC unavailable")
				},
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewEmissionHandler(tt.client)
			c, rec := newCtxWithUser(e, http.MethodGet, "", 1)
			if err := h.GetMonthlyReport(c); err != nil {
				t.Fatalf("handler error: %v", err)
			}
			if rec.Code != tt.wantStatus {
				t.Fatalf("want %d, got %d — body: %s", tt.wantStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

