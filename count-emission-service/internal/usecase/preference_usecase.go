package usecase

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/preference"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PreferenceUseCase struct {
	PreferenceRepository domain.PreferenceRepository
}

func NewPreferenceUseCase(repo domain.PreferenceRepository) domain.PreferenceUseCase {
	return &PreferenceUseCase{PreferenceRepository: repo}
}

func (uc *PreferenceUseCase) GetUserPreferences(ctx context.Context, userID int32) (*preference.UserEmissionPreference, error) {
	pref, err := uc.PreferenceRepository.GetUserPreferences(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &preference.UserEmissionPreference{
				UserId:      userID,
				CountryCode: "IDN",
			}, nil
		}
		return nil, err
	}
	return pref, nil
}

func (uc *PreferenceUseCase) SetUserPreferences(ctx context.Context, userID int32, countryCode string, customLimit *float64) (*preference.UserEmissionPreference, error) {
	if countryCode == "" {
		countryCode = "IDN"
	}
	pref := preference.UserEmissionPreference{
		UserId:                userID,
		CountryCode:           countryCode,
		CustomDailyLimitKgCo2: customLimit,
		UpdatedAt:             time.Now(),
	}
	result, err := uc.PreferenceRepository.UpsertUserPreferences(ctx, pref)
	if err != nil {
		return nil, err
	}
	return result, nil
}
