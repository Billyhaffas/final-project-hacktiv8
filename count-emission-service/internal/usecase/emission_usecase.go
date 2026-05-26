package usecase

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/emission"
	"count-emission-service/internal/model/thirdparty/carbonsutra"
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
	err = uc.EmissionRepository.CreateUserEmission(ctx, insertOrigin)
	if err != nil {
		return err
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
