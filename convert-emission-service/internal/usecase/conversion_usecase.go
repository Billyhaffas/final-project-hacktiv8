package usecase

import (
	"context"
	"fmt"

	"convert-emission-service/internal/domain"
)

type conversionUsecase struct {
	repo domain.CarbonPriceRepository
}

func NewConversionUsecase(r domain.CarbonPriceRepository) domain.ConversionUsecase {
	return &conversionUsecase{repo: r}
}

// Internal Helper to evaluate calculations uniformly
func (u *conversionUsecase) computeValuation(kgCo2 float64, rate *domain.CarbonPrice) domain.CarbonCostValuation {
	metricTons := kgCo2 / 1000.0
	costUsd := metricTons * rate.PricePerTonUsd
	costLocal := costUsd * rate.UsdCurRate

	return domain.CarbonCostValuation{
		TotalEmissionKgCo2: kgCo2,
		TotalCostUsd:       costUsd,
		TotalCostLocalCur:  costLocal,
	}
}

func (u *conversionUsecase) ConvertDailyEmission(ctx context.Context, countryCode string, emission domain.UserDailyEmission) (*domain.UserDailyCostResponse, error) {
	rate, err := u.repo.GetLatestRateByCountryCode(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate parameters: %w", err)
	}

	return &domain.UserDailyCostResponse{
		UserId:    emission.UserId,
		Date:      emission.Date,
		Valuation: u.computeValuation(emission.TotalEmissionKgCo2, rate),
	}, nil
}

func (u *conversionUsecase) ConvertMonthlyEmission(ctx context.Context, countryCode string, emission domain.UserMonthlyEmission) (*domain.UserMonthlyCostResponse, error) {
	rate, err := u.repo.GetLatestRateByCountryCode(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate parameters: %w", err)
	}

	dailyCosts := make([]domain.UserDailyCostResponse, len(emission.DailyEmissions))
	for i, d := range emission.DailyEmissions {
		dailyCosts[i] = domain.UserDailyCostResponse{
			UserId:    d.UserId,
			Date:      d.Date,
			Valuation: u.computeValuation(d.TotalEmissionKgCo2, rate),
		}
	}

	return &domain.UserMonthlyCostResponse{
		UserId:          emission.UserId,
		DailyCostDetail: dailyCosts,
		TotalValuation:  u.computeValuation(emission.TotalEmissionMonthlyKgCo2, rate),
	}, nil
}

func (u *conversionUsecase) ConvertToIDR(ctx context.Context, emissionKgCo2 float64) (pricePerTonUsd, exchangeRateUsdIdr, totalIdr float64, err error) {
	rate, err := u.repo.GetLatest(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to fetch carbon price: %w", err)
	}
	metricTons := emissionKgCo2 / 1000.0
	totalIdr = metricTons * rate.PricePerTonUsd * rate.UsdCurRate
	return rate.PricePerTonUsd, rate.UsdCurRate, totalIdr, nil
}

func (u *conversionUsecase) ConvertYearlyEmission(ctx context.Context, countryCode string, emission domain.UserYearlyEmission) (*domain.UserYearlyCostResponse, error) {
	rate, err := u.repo.GetLatestRateByCountryCode(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rate parameters: %w", err)
	}

	monthlyCosts := make([]domain.UserMonthlyCostResponse, len(emission.MonthlyEmissions))
	for i, m := range emission.MonthlyEmissions {
		monthlyCosts[i] = domain.UserMonthlyCostResponse{
			UserId:         m.UserId,
			TotalValuation: u.computeValuation(m.TotalEmissionKgCo2, rate),
		}
	}

	return &domain.UserYearlyCostResponse{
		UserId:            emission.UserId,
		MonthlyCostDetail: monthlyCosts,
		TotalValuation:    u.computeValuation(emission.TotalEmissionYearlyKgCo2, rate),
	}, nil
}
