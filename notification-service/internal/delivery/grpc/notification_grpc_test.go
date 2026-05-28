package grpchandler_test

import (
	"context"
	"errors"
	"testing"

	grpchandler "notification-service/internal/delivery/grpc"
	pb "notification-service/proto/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- mock usecase ---

type mockUsecase struct {
	checkDailyAlert func(context.Context, int32, string) (bool, float64, float64, string, string, error)
}

func (m *mockUsecase) CheckDailyAlert(ctx context.Context, userID int32, date string) (bool, float64, float64, string, string, error) {
	if m.checkDailyAlert != nil {
		return m.checkDailyAlert(ctx, userID, date)
	}
	return false, 3.0, 6.3, "country", "Your emission is within safe limits.", nil
}

// --- tests ---

func TestCheckDailyAlert_Handler(t *testing.T) {
	tests := []struct {
		name         string
		req          *pb.DailyAlertRequest
		usecase      *mockUsecase
		wantCode     codes.Code
		wantExceeded bool
		wantSource   string
	}{
		{
			name:         "OK — not exceeded, country threshold",
			req:          &pb.DailyAlertRequest{UserId: 1, Date: "2026-05-28"},
			usecase:      &mockUsecase{},
			wantCode:     codes.OK,
			wantExceeded: false,
			wantSource:   "country",
		},
		{
			name: "OK — exceeded, user threshold",
			req:  &pb.DailyAlertRequest{UserId: 2, Date: "2026-05-28"},
			usecase: &mockUsecase{
				checkDailyAlert: func(_ context.Context, _ int32, _ string) (bool, float64, float64, string, string, error) {
					return true, 8.0, 5.0, "user", "Alert!", nil
				},
			},
			wantCode:     codes.OK,
			wantExceeded: true,
			wantSource:   "user",
		},
		{
			name:     "InvalidArgument — user_id = 0",
			req:      &pb.DailyAlertRequest{UserId: 0, Date: "2026-05-28"},
			usecase:  &mockUsecase{},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Internal — usecase error",
			req:  &pb.DailyAlertRequest{UserId: 1, Date: "2026-05-28"},
			usecase: &mockUsecase{
				checkDailyAlert: func(_ context.Context, _ int32, _ string) (bool, float64, float64, string, string, error) {
					return false, 0, 0, "", "", errors.New("count-emission-service unavailable")
				},
			},
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := grpchandler.NewNotificationGRPCServer(tt.usecase)
			resp, err := srv.CheckDailyAlert(context.Background(), tt.req)

			code := codes.OK
			if err != nil {
				code = status.Code(err)
			}
			if code != tt.wantCode {
				t.Fatalf("want code %v, got %v (err: %v)", tt.wantCode, code, err)
			}
			if tt.wantCode == codes.OK {
				if resp.IsExceeded != tt.wantExceeded {
					t.Errorf("isExceeded: want %v, got %v", tt.wantExceeded, resp.IsExceeded)
				}
				if resp.ThresholdSource != tt.wantSource {
					t.Errorf("thresholdSource: want %q, got %q", tt.wantSource, resp.ThresholdSource)
				}
			}
		})
	}
}
