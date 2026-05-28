package usecase_test

import (
	"context"
	"errors"
	"testing"

	"convert-emission-service/internal/domain"
	"convert-emission-service/internal/usecase"
)

type mockRepo struct {
	getLatest                  func(context.Context) (*domain.CarbonPrice, error)
	getLatestRateByCountryCode func(context.Context, string) (*domain.CarbonPrice, error)
}

func (m *mockRepo) GetLatest(ctx context.Context) (*domain.CarbonPrice, error) {
	if m.getLatest != nil {
		return m.getLatest(ctx)
	}
	return &domain.CarbonPrice{PricePerTonUsd: 23.0, UsdCurRate: 16250.0}, nil
}

func (m *mockRepo) GetLatestRateByCountryCode(ctx context.Context, code string) (*domain.CarbonPrice, error) {
	if m.getLatestRateByCountryCode != nil {
		return m.getLatestRateByCountryCode(ctx, code)
	}
	return &domain.CarbonPrice{PricePerTonUsd: 23.0, UsdCurRate: 16250.0}, nil
}

func TestConvertToIDR(t *testing.T) {
	tests := []struct {
		name             string
		emissionKg       float64
		repo             *mockRepo
		wantPriceUsd     float64
		wantExchangeRate float64
		wantTotalIDR     float64
		wantErr          bool
	}{
		{
			name:       "success — 1000 kg",
			emissionKg: 1000.0,
			repo:       &mockRepo{},
			// 1000kg / 1000 = 1 ton × $23 × 16250 = 373,750 IDR
			wantPriceUsd:     23.0,
			wantExchangeRate: 16250.0,
			wantTotalIDR:     373750.0,
		},
		{
			name:       "success — 500 kg",
			emissionKg: 500.0,
			repo:       &mockRepo{},
			// 500/1000 = 0.5 ton × $23 × 16250 = 186,875 IDR
			wantPriceUsd:     23.0,
			wantExchangeRate: 16250.0,
			wantTotalIDR:     186875.0,
		},
		{
			name:             "success — zero emission",
			emissionKg:       0,
			repo:             &mockRepo{},
			wantPriceUsd:     23.0,
			wantExchangeRate: 16250.0,
			wantTotalIDR:     0,
		},
		{
			name:       "repo error",
			emissionKg: 100.0,
			repo: &mockRepo{
				getLatest: func(_ context.Context) (*domain.CarbonPrice, error) {
					return nil, errors.New("mongo unavailable")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := usecase.NewConversionUsecase(tt.repo)
			priceUsd, rate, totalIDR, err := uc.ConvertToIDR(context.Background(), tt.emissionKg)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if priceUsd != tt.wantPriceUsd {
				t.Errorf("priceUsd: want %.2f, got %.2f", tt.wantPriceUsd, priceUsd)
			}
			if rate != tt.wantExchangeRate {
				t.Errorf("exchangeRate: want %.2f, got %.2f", tt.wantExchangeRate, rate)
			}
			if totalIDR != tt.wantTotalIDR {
				t.Errorf("totalIDR: want %.2f, got %.2f", tt.wantTotalIDR, totalIDR)
			}
		})
	}
}
