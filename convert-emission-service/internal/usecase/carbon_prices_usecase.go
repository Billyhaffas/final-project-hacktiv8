package usecase

import (
	"context"
	"convert-emission-service/internal/domain"
	"fmt"
)

type SeedCarbonPricesUseCase struct {
	provider   domain.CarbonPriceProvider
	repository domain.CarbonPriceRepository
}

func NewSeedCarbonPricesUseCase(p domain.CarbonPriceProvider, r domain.CarbonPriceRepository) *SeedCarbonPricesUseCase {
	return &SeedCarbonPricesUseCase{
		provider:   p,
		repository: r,
	}
}

func (u *SeedCarbonPricesUseCase) Execute(ctx context.Context) error {
	fmt.Println("Fetching carbon price data from remote source...")
	carbonPrices, err := u.provider.FetchPrices(ctx)
	if err != nil {
		return fmt.Errorf("error fetching data from provider: %w", err)
	}

	fmt.Println("Analyzing database status and applying carbon price records...")
	seeded, err := u.repository.BulkInsert(ctx, carbonPrices)
	if err != nil {
		return fmt.Errorf("error executing bulk storage: %w", err)
	}

	if !seeded {
		fmt.Println("Database collection already contains data. Seeding phase bypassed.")
		return nil
	}

	fmt.Printf("Successfully inserted %d records! Initial seed completed.\n", len(carbonPrices))
	return nil
}
