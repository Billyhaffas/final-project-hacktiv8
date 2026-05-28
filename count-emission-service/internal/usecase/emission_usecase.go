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

var zeroEmissionTypes = map[string]bool{
	"Bicycle": true,
	"Walk":    true,
}

func (uc *EmissionUseCase) CreateUserEmission(ctx context.Context, request *emission.EmissionBody) error {
	var emissionKg float64

	if !zeroEmissionTypes[request.VehicleType] {
		payload := carbonsutra.CountEmisionBodyPayload{
			VehicleType:   request.VehicleType,
			FuelType:      request.FuelType,
			DistanceValue: request.DistanceKm,
			DistanceUnit:  "km",
			IncludeWtt:    "Y",
		}
		responTP, err := uc.CarbonSutraRepository.GetCarbonEmission(payload)
		if err != nil {
			return fmt.Errorf("CreateUserEmission: external API: %w", err)
		}
		emissionKg = responTP.Data.Co2eKg
	}

	insertOrigin := emission.EmissionOrigin{
		UserId:        request.UserId,
		VehicleType:   request.VehicleType,
		FuelType:      request.FuelType,
		DistanceKm:    request.DistanceKm,
		EmissionKgCo2: emissionKg,
		RecordedAt:    time.Now(),
		CreatedAt:     time.Now(),
	}
	if err := uc.EmissionRepository.CreateUserEmission(ctx, insertOrigin); err != nil {
		return fmt.Errorf("CreateUserEmission: %w", err)
	}
	return nil
}

func (uc *EmissionUseCase) GetUserDailyEmission(ctx context.Context, userId int32) (*emission.UserDailyEmission, error) {
	userEmission, err := uc.EmissionRepository.GetUserDailyEmission(ctx, userId)
	if err != nil {
		return nil, err
	}
	return userEmission, nil
}

func (uc *EmissionUseCase) GetUserMonthlyEmission(ctx context.Context, userId int32) (*emission.UserMonthlyEmission, error) {
	userMonthlyEmission, err := uc.EmissionRepository.GetUserMonthlyEmission(ctx, userId)
	if err != nil {
		return nil, err
	}
	return userMonthlyEmission, nil
}

func (uc *EmissionUseCase) GetUserYearlyEmission(ctx context.Context, userId int32) (*emission.UserYearlyEmission, error) {
	userYearlyEmission, err := uc.EmissionRepository.GetUserYearlyEmission(ctx, userId)
	if err != nil {
		return nil, err
	}
	return userYearlyEmission, nil
}
