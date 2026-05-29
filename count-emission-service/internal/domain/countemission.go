package domain

import (
	"context"
	"count-emission-service/internal/model/emission"
	"count-emission-service/internal/model/preference"
	"count-emission-service/internal/model/thirdparty/carbonsutra"

	"github.com/labstack/echo/v5"
)

type CarbonSutraRepository interface {
	GetCarbonEmission(payload carbonsutra.CountEmisionBodyPayload) (*carbonsutra.CountEmisionThirdParty, error)
}
type EmissionRepository interface {
	CreateUserEmission(ctx context.Context, req emission.EmissionOrigin) error
	GetUserDailyEmission(ctx context.Context, userId int32) (*emission.UserDailyEmission, error)
	GetUserMonthlyEmission(ctx context.Context, userId int32) (*emission.UserMonthlyEmission, error)
	GetUserYearlyEmission(ctx context.Context, userId int32) (*emission.UserYearlyEmission, error)
	GetDailyTotal(ctx context.Context, userId int32, date string) (float64, int32, error)
}
type EmissionUseCase interface {
	CreateUserEmission(ctx context.Context, request *emission.EmissionBody) error
	GetUserDailyEmission(ctx context.Context, userId int32) (*emission.UserDailyEmission, error)
	GetUserMonthlyEmission(ctx context.Context, userId int32) (*emission.UserMonthlyEmission, error)
	GetUserYearlyEmission(ctx context.Context, userId int32) (*emission.UserYearlyEmission, error)
	GetDailyTotal(ctx context.Context, userId int32, date string) (float64, int32, error)
}
type EmissionHandler interface {
	CreateUserEmission(c *echo.Context) error
	GetUserDailyEmission(c *echo.Context) error
	GetUserMonthlyEmission(c *echo.Context) error
	GetUserYearlyEmission(c *echo.Context) error
}

type PreferenceRepository interface {
	GetUserPreferences(ctx context.Context, userID int32) (*preference.UserEmissionPreference, error)
	UpsertUserPreferences(ctx context.Context, pref preference.UserEmissionPreference) (*preference.UserEmissionPreference, error)
}

type PreferenceUseCase interface {
	GetUserPreferences(ctx context.Context, userID int32) (*preference.UserEmissionPreference, error)
	SetUserPreferences(ctx context.Context, userID int32, countryCode string, customLimit *float64) (*preference.UserEmissionPreference, error)
}
