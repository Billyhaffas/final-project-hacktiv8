package domain

import (
	"context"
	"time"
)

// Request model

type UserDailyEmission struct {
	UserId             int8    `json:"user_id"`
	Date               string  `json:"date"`
	TotalEmissionKgCo2 float64 `json:"total_emission_kg_co2"`
}

type UserMonthlyEmissionDetail struct {
	UserId             int8    `json:"user_id"`
	Month              string  `json:"month"`
	TotalEmissionKgCo2 float64 `json:"total_emission_kg_co2"`
}

type UserMonthlyEmission struct {
	UserId                    int8                `json:"user_id"`
	DailyEmissions            []UserDailyEmission `json:"daily_emissions"`
	TotalEmissionMonthlyKgCo2 float64             `json:"total_emission_monthly_kg_co2"`
}

type UserYearlyEmission struct {
	UserId                   int8                        `json:"user_id"`
	MonthlyEmissions         []UserMonthlyEmissionDetail `json:"monthly_emissions"`
	TotalEmissionYearlyKgCo2 float64                     `json:"total_emission_yearly_kg_co2"`
}

// MongoDB Model

type CarbonPrice struct {
	ID             string    `json:"id"`
	PricePerTonUsd float64   `json:"price_per_ton_usd"`
	UsdCurRate     float64   `json:"usd_cur_rate"`
	Source         string    `json:"source"`
	FetchedAt      time.Time `json:"fetched_at"`
}

// Calculation

type CarbonCostValuation struct {
	TotalEmissionKgCo2 float64 `json:"total_emission_kg_co2"`
	TotalCostUsd       float64 `json:"total_cost_usd"`
	TotalCostLocalCur  float64 `json:"total_cost_local_cur"`
}

// Response

type UserDailyCostResponse struct {
	UserId    int8                `json:"user_id"`
	Date      string              `json:"date"`
	Valuation CarbonCostValuation `json:"valuation"`
}

type UserMonthlyCostResponse struct {
	UserId          int8                    `json:"user_id"`
	DailyCostDetail []UserDailyCostResponse `json:"daily_cost_detail"`
	TotalValuation  CarbonCostValuation     `json:"total_valuation"`
}

type UserYearlyCostResponse struct {
	UserId            int8                      `json:"user_id"`
	MonthlyCostDetail []UserMonthlyCostResponse `json:"monthly_cost_detail"`
	TotalValuation    CarbonCostValuation       `json:"total_valuation"`
}

// CarbonPriceRepository manages querying conversion rate metrics
type CarbonPriceRepository interface {
	GetLatestRateByCountryCode(ctx context.Context, countryCode string) (*CarbonPrice, error)
}

// ConversionUsecase encapsulates the calculation orchestrations
type ConversionUsecase interface {
	ConvertDailyEmission(ctx context.Context, countryCode string, emission UserDailyEmission) (*UserDailyCostResponse, error)
	ConvertMonthlyEmission(ctx context.Context, countryCode string, emission UserMonthlyEmission) (*UserMonthlyCostResponse, error)
	ConvertYearlyEmission(ctx context.Context, countryCode string, emission UserYearlyEmission) (*UserYearlyCostResponse, error)
}
