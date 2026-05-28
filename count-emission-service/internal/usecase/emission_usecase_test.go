package usecase_test

import (
	"context"
	"testing"

	"count-emission-service/internal/model/emission"
	"count-emission-service/internal/model/thirdparty/carbonsutra"
	exMocks "count-emission-service/internal/repository/external_api/mocks"
	inmocks "count-emission-service/internal/repository/internal_repo/mocks"
	"count-emission-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserEmission(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(mockEmissionRepo, mockCarbonRepo)
		req := &emission.EmissionBody{
			UserId:      1,
			VehicleType: "Car-Size-Medium",
			FuelType:    "Petrol",
			DistanceKm:  10,
		}

		carbonResponse := &carbonsutra.CountEmisionThirdParty{
			Data: carbonsutra.EmissionData{
				Co2eKg: 5.5,
			},
		}

		mockCarbonRepo.
			On("GetCarbonEmission", mock.Anything).
			Return(carbonResponse, nil)

		mockEmissionRepo.
			On(
				"CreateUserEmission",
				mock.Anything,
				mock.Anything,
			).
			Return(nil)

		err := uc.CreateUserEmission(context.Background(), req)

		assert.NoError(t, err)
	})

	t.Run("carbon sutra api error", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		req := &emission.EmissionBody{
			UserId:      1,
			VehicleType: "Car-Size-Medium",
			FuelType:    "Petrol",
			DistanceKm:  10,
		}

		mockCarbonRepo.
			On("GetCarbonEmission", mock.Anything).
			Return(nil, assert.AnError)

		err := uc.CreateUserEmission(context.Background(), req)

		assert.Error(t, err)
	})

	t.Run("repository insert error", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		req := &emission.EmissionBody{
			UserId:      1,
			VehicleType: "Car",
			FuelType:    "Petrol",
			DistanceKm:  10,
		}

		carbonResponse := &carbonsutra.CountEmisionThirdParty{
			Data: carbonsutra.EmissionData{
				Co2eKg: 5.5,
			},
		}

		mockCarbonRepo.
			On("GetCarbonEmission", mock.Anything).
			Return(carbonResponse, nil)

		mockEmissionRepo.
			On(
				"CreateUserEmission",
				mock.Anything,
				mock.Anything,
			).
			Return(assert.AnError)

		err := uc.CreateUserEmission(context.Background(), req)

		assert.Error(t, err)
	})
}

func TestGetUserDailyEmission(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		expected := &emission.UserDailyEmission{
			UserId: 1,
		}

		mockEmissionRepo.
			On(
				"GetUserDailyEmission",
				context.Background(),
				int32(1),
			).
			Return(expected, nil)

		result, err := uc.GetUserDailyEmission(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.UserId)

		mockEmissionRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		mockEmissionRepo.
			On(
				"GetUserDailyEmission",
				context.Background(),
				int32(1),
			).
			Return(nil, assert.AnError)

		result, err := uc.GetUserDailyEmission(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockEmissionRepo.AssertExpectations(t)
	})
}

func TestGetUserMonthlyEmission(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		expected := &emission.UserMonthlyEmission{
			UserId: 1,
		}

		mockEmissionRepo.
			On(
				"GetUserMonthlyEmission",
				context.Background(),
				int32(1),
			).
			Return(expected, nil)

		result, err := uc.GetUserMonthlyEmission(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.UserId)

		mockEmissionRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		mockEmissionRepo.
			On(
				"GetUserMonthlyEmission",
				context.Background(),
				int32(1),
			).
			Return(nil, assert.AnError)

		result, err := uc.GetUserMonthlyEmission(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockEmissionRepo.AssertExpectations(t)
	})
}

func TestGetUserYearlyEmission(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		expected := &emission.UserYearlyEmission{
			UserId: 1,
		}

		mockEmissionRepo.
			On(
				"GetUserYearlyEmission",
				context.Background(),
				int32(1),
			).
			Return(expected, nil)

		result, err := uc.GetUserYearlyEmission(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.UserId)

		mockEmissionRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {

		mockEmissionRepo := new(inmocks.MockEmissionRepository)
		mockCarbonRepo := new(exMocks.MockCarbonSutraRepository)

		uc := usecase.NewEmissionUseCase(
			mockEmissionRepo,
			mockCarbonRepo,
		)

		mockEmissionRepo.
			On(
				"GetUserYearlyEmission",
				context.Background(),
				int32(1),
			).
			Return(nil, assert.AnError)

		result, err := uc.GetUserYearlyEmission(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockEmissionRepo.AssertExpectations(t)
	})
}
