package usecase

import (
	"context"
	"count-emission-service/internal/domain"
	"fmt"
)

type SeedEmissionUseCase struct {
	provider   domain.EmissionProvider
	repository domain.EmissionRepository
}

func NewSeedEmissionUseCase(p domain.EmissionProvider, r domain.EmissionRepository) *SeedEmissionUseCase {
	return &SeedEmissionUseCase{
		provider:   p,
		repository: r,
	}
}

func (u *SeedEmissionUseCase) Execute(ctx context.Context) error {
	fmt.Println("Fetching emission data from remote source...")
	emissions, err := u.provider.FetchEmissions(ctx)
	if err != nil {
		return fmt.Errorf("error fetching data from provider: %w", err)
	}

	fmt.Println("Analyzing database status and applying emission records...")
	seeded, err := u.repository.BulkInsert(ctx, emissions)
	if err != nil {
		return fmt.Errorf("error executing bulk storage: %w", err)
	}

	if !seeded {
		fmt.Println("Database collection already contains data. Seeding phase bypassed.")
		return nil
	}

	fmt.Printf("Successfully inserted %d records! Initial seed completed.\n", len(emissions))
	return nil
}
