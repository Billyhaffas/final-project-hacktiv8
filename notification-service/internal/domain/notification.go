package domain

import "context"

type EmissionClient interface {
	GetDailyTotal(ctx context.Context, userID int32, date string) (float64, error)
	GetUserPreferences(ctx context.Context, userID int32) (countryCode string, customDailyLimitKg float64, err error)
}

type NotificationUsecase interface {
	CheckDailyAlert(ctx context.Context, userID int32, date string) (isExceeded bool, dailyTotalKg, dailyLimitKg float64, thresholdSource, message string, err error)
}
