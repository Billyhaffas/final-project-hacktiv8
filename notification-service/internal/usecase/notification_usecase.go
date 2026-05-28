package usecase

import (
	"context"
	"fmt"
	"time"

	"notification-service/internal/domain"
)

// countryDefaultLimits holds the daily CO₂ limit (kg) per country code.
// Source: Our World in Data — 2.3 tons/year ÷ 365 for IDN.
var countryDefaultLimits = map[string]float64{
	"IDN": 6.3,
}

const fallbackDailyLimitKg = 6.3

type notificationUsecase struct {
	emissionClient domain.EmissionClient
}

func NewNotificationUsecase(ec domain.EmissionClient) domain.NotificationUsecase {
	return &notificationUsecase{emissionClient: ec}
}

func (u *notificationUsecase) CheckDailyAlert(
	ctx context.Context,
	userID int32,
	date string,
) (bool, float64, float64, string, string, error) {
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	dailyTotal, err := u.emissionClient.GetDailyTotal(ctx, userID, date)
	if err != nil {
		return false, 0, 0, "", "", fmt.Errorf("CheckDailyAlert: %w", err)
	}

	countryCode, customLimit, err := u.emissionClient.GetUserPreferences(ctx, userID)
	if err != nil {
		return false, 0, 0, "", "", fmt.Errorf("CheckDailyAlert: %w", err)
	}

	var dailyLimit float64
	var thresholdSource string

	if customLimit > 0 {
		dailyLimit = customLimit
		thresholdSource = "user"
	} else {
		if limit, ok := countryDefaultLimits[countryCode]; ok {
			dailyLimit = limit
		} else {
			dailyLimit = fallbackDailyLimitKg
		}
		thresholdSource = "country"
	}

	isExceeded := dailyTotal > dailyLimit

	var msg string
	if isExceeded {
		msg = fmt.Sprintf(
			"Alert! Your daily carbon emission (%.2f kg CO₂) has exceeded your limit of %.2f kg CO₂.",
			dailyTotal, dailyLimit,
		)
	} else {
		msg = "Your emission is within safe limits."
	}

	return isExceeded, dailyTotal, dailyLimit, thresholdSource, msg, nil
}
