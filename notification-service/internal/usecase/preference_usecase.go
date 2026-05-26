package usecase

import (
	"context"
	"errors"
	"fmt"
	"notification-service/internal/domain"

	"gorm.io/gorm"
)

type preferenceUsecase struct {
	repo domain.PreferenceRepository
}

func NewPreferenceUsecase(r domain.PreferenceRepository) domain.PreferenceUsecase {
	return &preferenceUsecase{repo: r}
}

func (u *preferenceUsecase) GetPreference(ctx context.Context, userID int) (*domain.UserEmissionPreference, error) {
	pref, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("preferences not found for user %d", userID)
		}
		return nil, err
	}
	return pref, nil
}

func (u *preferenceUsecase) SavePreference(ctx context.Context, input domain.PreferenceUpsertInput) (*domain.UserEmissionPreference, error) {
	// Business validation: verify country code is exactly 3 characters
	if len(input.CountryCode) != 3 {
		return nil, errors.New("invalid country code: must be a 3-character ISO code")
	}

	existing, err := u.repo.GetByUserID(ctx, input.UserID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed checking existence: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// CRUD: CREATE
		newPref := &domain.UserEmissionPreference{
			UserID:                input.UserID,
			CountryCode:           input.CountryCode,
			CustomDailyLimitKgCo2: input.CustomDailyLimitKgCo2,
		}
		if err := u.repo.Create(ctx, newPref); err != nil {
			return nil, fmt.Errorf("failed creating preference: %w", err)
		}
		return newPref, nil
	}

	// CRUD: UPDATE
	existing.CountryCode = input.CountryCode
	existing.CustomDailyLimitKgCo2 = input.CustomDailyLimitKgCo2

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed updating preference: %w", err)
	}

	return existing, nil
}

func (u *preferenceUsecase) DeletePreference(ctx context.Context, userID int) error {
	// Verify it exists before trying to delete it
	_, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cannot delete: preferences do not exist for user %d", userID)
		}
		return err
	}
	return u.repo.Delete(ctx, userID)
}
