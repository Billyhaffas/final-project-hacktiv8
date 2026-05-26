package usecase

import (
	"context"
	"errors"
	"fmt"
	"notification-service/internal/domain"
	"notification-service/internal/repository"

	"gorm.io/gorm"
)

type notificationUsecase struct {
	prefRepo       domain.PreferenceRepository  // PostgreDB
	masterRepo     domain.MasterLimitRepository // MongoDB
	emissionClient *repository.EmissionClient   // count-service
}

func NewNotificationUsecase(
	p domain.PreferenceRepository,
	m domain.MasterLimitRepository,
	e *repository.EmissionClient,
) domain.NotificationUsecase {
	return &notificationUsecase{
		prefRepo:       p,
		masterRepo:     m,
		emissionClient: e,
	}
}

func (u *notificationUsecase) CheckAndSendNotification(ctx context.Context, userID int) (bool, string, error) {
	// 1. Fetch current daily emissions from the external API
	currentEmission, err := u.emissionClient.GetDailyEmission(ctx, userID)
	if err != nil {
		return false, "", fmt.Errorf("failed to fetch current emissions: %w", err)
	}

	// 2. Fetch user threshold preferences from Postgres
	pref, err := u.prefRepo.GetByUserID(ctx, userID)

	// Separate true DB connection errors from a missing preferences row record
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, "", fmt.Errorf("failed to fetch user preferences: %w", err)
	}

	var allowedLimit float64
	var countryCode string

	// 3. Fallback evaluation logic
	// If preference record is missing completely, or custom limit isn't valid/configured (<= 0)
	if errors.Is(err, gorm.ErrRecordNotFound) || pref.CustomDailyLimitKgCo2 <= 0 {

		// Fallback to a default country code if preference row does not exist entirely
		countryCode = "IDN"
		if pref != nil && pref.CountryCode != "" {
			countryCode = pref.CountryCode
		}

		// Fetch the latest master baseline threshold from MongoDB for this country code
		fallbackLimit, err := u.masterRepo.GetDefaultLimitByCountry(ctx, countryCode)
		if err != nil {
			return false, "", fmt.Errorf("failed to fetch master fallback limit for country %s: %w", countryCode, err)
		}
		allowedLimit = fallbackLimit
	} else {
		// If custom preference is present and explicitly filled out, use it
		allowedLimit = pref.CustomDailyLimitKgCo2
		countryCode = pref.CountryCode
	}

	// 4. Check if threshold limit is breached
	if currentEmission.TotalEmissionKgCo2 > allowedLimit {
		msg := fmt.Sprintf(
			"Alert! Your daily carbon emission (%.2f kg CO2) has exceeded your allowed limit of %.2f kg CO2.",
			currentEmission.TotalEmissionKgCo2,
			allowedLimit,
		)
		return true, msg, nil
	}

	return false, "Emission within safe limits", nil
}
