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
}
type EmissionUseCase interface {
	CreateUserEmission(ctx context.Context, request *emission.EmissionBody) error
}
type EmissionHandler interface {
	CreateUserEmission(c *echo.Context) error
}
