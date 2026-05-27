package domain

import (
	"context"
	"count-emission-service/internal/model/emission"
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
}
type EmissionUseCase interface {
	CreateUserEmission(ctx context.Context, request *emission.EmissionBody) error
	GetUserDailyEmission(ctx context.Context, userId int32) (*emission.UserDailyEmission, error)
	GetUserMonthlyEmission(ctx context.Context, userId int32) (*emission.UserMonthlyEmission, error)
	GetUserYearlyEmission(ctx context.Context, userId int32) (*emission.UserYearlyEmission, error)
}
type EmissionHandler interface {
	CreateUserEmission(c *echo.Context) error
	GetUserDailyEmission(c *echo.Context) error
	GetUserMonthlyEmission(c *echo.Context) error
	GetUserYearlyEmission(c *echo.Context) error
}
