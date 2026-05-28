package usecase_test

import (
	"context"
	"errors"
	"testing"

	"notification-service/internal/domain"
	"notification-service/internal/usecase"
)

// --- mock EmissionClient ---

type mockEmissionClient struct {
	getDailyTotal      func(context.Context, int32, string) (float64, error)
	getUserPreferences func(context.Context, int32) (string, float64, error)
}

func (m *mockEmissionClient) GetDailyTotal(ctx context.Context, userID int32, date string) (float64, error) {
	if m.getDailyTotal != nil {
		return m.getDailyTotal(ctx, userID, date)
	}
	return 0, nil
}

func (m *mockEmissionClient) GetUserPreferences(ctx context.Context, userID int32) (string, float64, error) {
	if m.getUserPreferences != nil {
		return m.getUserPreferences(ctx, userID)
	}
	return "IDN", 0, nil
}

var _ domain.EmissionClient = (*mockEmissionClient)(nil)

// --- tests ---

func TestCheckDailyAlert(t *testing.T) {
	tests := []struct {
		name            string
		client          *mockEmissionClient
		userID          int32
		date            string
		wantExceeded    bool
		wantSource      string
		wantDailyTotal  float64
		wantDailyLimit  float64
		wantErr         bool
	}{
		{
			name: "custom limit — not exceeded",
			client: &mockEmissionClient{
				getDailyTotal:      func(_ context.Context, _ int32, _ string) (float64, error) { return 3.0, nil },
				getUserPreferences: func(_ context.Context, _ int32) (string, float64, error) { return "IDN", 5.0, nil },
			},
			wantExceeded:   false,
			wantSource:     "user",
			wantDailyTotal: 3.0,
			wantDailyLimit: 5.0,
		},
		{
			name: "custom limit — exceeded",
			client: &mockEmissionClient{
				getDailyTotal:      func(_ context.Context, _ int32, _ string) (float64, error) { return 7.0, nil },
				getUserPreferences: func(_ context.Context, _ int32) (string, float64, error) { return "IDN", 5.0, nil },
			},
			wantExceeded:   true,
			wantSource:     "user",
			wantDailyTotal: 7.0,
			wantDailyLimit: 5.0,
		},
		{
			name: "country default IDN — not exceeded",
			client: &mockEmissionClient{
				getDailyTotal:      func(_ context.Context, _ int32, _ string) (float64, error) { return 4.0, nil },
				getUserPreferences: func(_ context.Context, _ int32) (string, float64, error) { return "IDN", 0, nil },
			},
			wantExceeded:   false,
			wantSource:     "country",
			wantDailyTotal: 4.0,
			wantDailyLimit: 6.3,
		},
		{
			name: "country default IDN — exceeded",
			client: &mockEmissionClient{
				getDailyTotal:      func(_ context.Context, _ int32, _ string) (float64, error) { return 8.0, nil },
				getUserPreferences: func(_ context.Context, _ int32) (string, float64, error) { return "IDN", 0, nil },
			},
			wantExceeded:   true,
			wantSource:     "country",
			wantDailyTotal: 8.0,
			wantDailyLimit: 6.3,
		},
		{
			name: "unknown country — fallback limit 6.3",
			client: &mockEmissionClient{
				getDailyTotal:      func(_ context.Context, _ int32, _ string) (float64, error) { return 5.0, nil },
				getUserPreferences: func(_ context.Context, _ int32) (string, float64, error) { return "SGP", 0, nil },
			},
			wantExceeded:   false,
			wantSource:     "country",
			wantDailyTotal: 5.0,
			wantDailyLimit: 6.3,
		},
		{
			name: "GetDailyTotal error",
			client: &mockEmissionClient{
				getDailyTotal: func(_ context.Context, _ int32, _ string) (float64, error) {
					return 0, errors.New("count-emission-service unavailable")
				},
			},
			wantErr: true,
		},
		{
			name: "GetUserPreferences error",
			client: &mockEmissionClient{
				getDailyTotal: func(_ context.Context, _ int32, _ string) (float64, error) { return 3.0, nil },
				getUserPreferences: func(_ context.Context, _ int32) (string, float64, error) {
					return "", 0, errors.New("preferences unavailable")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := usecase.NewNotificationUsecase(tt.client)
			isExceeded, dailyTotal, dailyLimit, source, _, err := uc.CheckDailyAlert(context.Background(), tt.userID, tt.date)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if isExceeded != tt.wantExceeded {
				t.Errorf("isExceeded: want %v, got %v", tt.wantExceeded, isExceeded)
			}
			if source != tt.wantSource {
				t.Errorf("thresholdSource: want %q, got %q", tt.wantSource, source)
			}
			if dailyTotal != tt.wantDailyTotal {
				t.Errorf("dailyTotal: want %.2f, got %.2f", tt.wantDailyTotal, dailyTotal)
			}
			if dailyLimit != tt.wantDailyLimit {
				t.Errorf("dailyLimit: want %.2f, got %.2f", tt.wantDailyLimit, dailyLimit)
			}
		})
	}
}
