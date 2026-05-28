package grpchandler_test

import (
	"context"
	"errors"
	"testing"

	"convert-emission-service/internal/domain"
	grpchandler "convert-emission-service/internal/delivery/grpc"
	pb "convert-emission-service/proto/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// --- mock usecase ---

type mockUsecase struct {
	convertToIDR func(context.Context, float64) (float64, float64, float64, error)
}

func (m *mockUsecase) ConvertToIDR(ctx context.Context, kg float64) (float64, float64, float64, error) {
	if m.convertToIDR != nil {
		return m.convertToIDR(ctx, kg)
	}
	return 23.0, 16250.0, (kg / 1000.0) * 23.0 * 16250.0, nil
}

func (m *mockUsecase) ConvertDailyEmission(ctx context.Context, cc string, e domain.UserDailyEmission) (*domain.UserDailyCostResponse, error) {
	return nil, nil
}

func (m *mockUsecase) ConvertMonthlyEmission(ctx context.Context, cc string, e domain.UserMonthlyEmission) (*domain.UserMonthlyCostResponse, error) {
	return nil, nil
}

func (m *mockUsecase) ConvertYearlyEmission(ctx context.Context, cc string, e domain.UserYearlyEmission) (*domain.UserYearlyCostResponse, error) {
	return nil, nil
}

// --- tests ---

func TestConvertToIDR_Handler(t *testing.T) {
	tests := []struct {
		name       string
		req        *pb.ConvertToIDRRequest
		usecase    *mockUsecase
		wantCode   codes.Code
		wantTotIDR float64
	}{
		{
			name:       "OK — 1000 kg",
			req:        &pb.ConvertToIDRRequest{EmissionKgCo2: 1000},
			usecase:    &mockUsecase{},
			wantCode:   codes.OK,
			wantTotIDR: 373750.0,
		},
		{
			name:       "OK — zero emission",
			req:        &pb.ConvertToIDRRequest{EmissionKgCo2: 0},
			usecase:    &mockUsecase{},
			wantCode:   codes.OK,
			wantTotIDR: 0,
		},
		{
			name:     "InvalidArgument — negative emission",
			req:      &pb.ConvertToIDRRequest{EmissionKgCo2: -1},
			usecase:  &mockUsecase{},
			wantCode: codes.InvalidArgument,
		},
		{
			name: "Internal — usecase error",
			req:  &pb.ConvertToIDRRequest{EmissionKgCo2: 100},
			usecase: &mockUsecase{
				convertToIDR: func(_ context.Context, _ float64) (float64, float64, float64, error) {
					return 0, 0, 0, errors.New("mongo unavailable")
				},
			},
			wantCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := grpchandler.NewConvertGRPCServer(tt.usecase)
			resp, err := srv.ConvertToIDR(context.Background(), tt.req)

			code := codes.OK
			if err != nil {
				code = status.Code(err)
			}
			if code != tt.wantCode {
				t.Fatalf("want code %v, got %v (err: %v)", tt.wantCode, code, err)
			}
			if tt.wantCode == codes.OK {
				if resp.TotalIdr != tt.wantTotIDR {
					t.Errorf("totalIdr: want %.2f, got %.2f", tt.wantTotIDR, resp.TotalIdr)
				}
			}
		})
	}
}
