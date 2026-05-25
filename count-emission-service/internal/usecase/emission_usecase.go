package usecase

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/emission"
	"count-emission-service/internal/model/thirdparty/carbonsutra"
	"fmt"
	"time"
)

type EmissionUseCase struct {
	EmissionRepository    domain.EmissionRepository
	CarbonSutraRepository domain.CarbonSutraRepository
}

func NewEmissionUseCase(EmissionRepo domain.EmissionRepository, CarbonSutraRepo domain.CarbonSutraRepository) domain.EmissionUseCase {
	return &EmissionUseCase{EmissionRepository: EmissionRepo, CarbonSutraRepository: CarbonSutraRepo}
}

func (uc *EmissionUseCase) CreateUserEmission(ctx context.Context, request *emission.EmissionBody) error {
	payload := carbonsutra.CountEmisionBodyPayload{
		VehicleType:   request.VehicleType,
		FuelType:      request.FuelType,
		DistanceValue: request.DistanceKm,
		DistanceUnit:  "km",
		IncludeWtt:    "Y",
	}
	responTP, err := uc.CarbonSutraRepository.GetCarbonEmission(payload)
	if err != nil {
		return err
	}
	insertOrigin := emission.EmissionOrigin{
		UserId:      request.UserId,
		VehicleType: request.VehicleType,
		FuelType:    request.FuelType,
		DistanceKm:  request.DistanceKm,
	}
	insertOrigin.EmissionKgCo2 = responTP.Data.Co2eKg
	insertOrigin.RecordedAt = time.Now()
	insertOrigin.CreatedAt = time.Now()
	fmt.Println("EmissionKgCo2", insertOrigin.EmissionKgCo2)
	err = uc.EmissionRepository.CreateUserEmission(ctx, insertOrigin)
	if err != nil {
		return err
	}
	return nil
}
