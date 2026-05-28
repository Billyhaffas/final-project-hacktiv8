package usecase_test

import (
	"context"
	"testing"

	"count-emission-service/internal/model/preference"
	inmocks "count-emission-service/internal/repository/internal_repo/mocks"
	"count-emission-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetUserPreferences(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mockPreferenceRepo := new(inmocks.MockPreferenceRepository)

		uc := usecase.NewPreferenceUseCase(
			mockPreferenceRepo,
		)

		expected := &preference.UserEmissionPreference{
			UserId:      1,
			CountryCode: "IDN",
		}

		mockPreferenceRepo.
			On(
				"GetUserPreferences",
				context.Background(),
				int32(1),
			).
			Return(expected, nil)

		result, err := uc.GetUserPreferences(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.UserId)
		assert.Equal(t, "IDN", result.CountryCode)

		mockPreferenceRepo.AssertExpectations(t)
	})

	t.Run("record not found should return default preference", func(t *testing.T) {

		mockPreferenceRepo := new(inmocks.MockPreferenceRepository)

		uc := usecase.NewPreferenceUseCase(
			mockPreferenceRepo,
		)

		mockPreferenceRepo.
			On(
				"GetUserPreferences",
				context.Background(),
				int32(1),
			).
			Return(nil, gorm.ErrRecordNotFound)

		result, err := uc.GetUserPreferences(
			context.Background(),
			1,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.UserId)
		assert.Equal(t, "IDN", result.CountryCode)

		mockPreferenceRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {

		mockPreferenceRepo := new(inmocks.MockPreferenceRepository)

		uc := usecase.NewPreferenceUseCase(
			mockPreferenceRepo,
		)

		mockPreferenceRepo.
			On(
				"GetUserPreferences",
				context.Background(),
				int32(1),
			).
			Return(nil, assert.AnError)

		result, err := uc.GetUserPreferences(
			context.Background(),
			1,
		)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockPreferenceRepo.AssertExpectations(t)
	})
}

func TestSetUserPreferences(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mockPreferenceRepo := new(inmocks.MockPreferenceRepository)

		uc := usecase.NewPreferenceUseCase(
			mockPreferenceRepo,
		)

		customLimit := 10.5

		expected := &preference.UserEmissionPreference{
			UserId:                1,
			CountryCode:           "IDN",
			CustomDailyLimitKgCo2: &customLimit,
		}

		mockPreferenceRepo.
			On(
				"UpsertUserPreferences",
				context.Background(),
				mock.Anything,
			).
			Return(expected, nil)

		result, err := uc.SetUserPreferences(
			context.Background(),
			1,
			"IDN",
			&customLimit,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int32(1), result.UserId)
		assert.Equal(t, "IDN", result.CountryCode)

		mockPreferenceRepo.AssertExpectations(t)
	})

	t.Run("empty country code should default to IDN", func(t *testing.T) {

		mockPreferenceRepo := new(inmocks.MockPreferenceRepository)

		uc := usecase.NewPreferenceUseCase(
			mockPreferenceRepo,
		)

		customLimit := 10.5

		expected := &preference.UserEmissionPreference{
			UserId:                1,
			CountryCode:           "IDN",
			CustomDailyLimitKgCo2: &customLimit,
		}

		mockPreferenceRepo.
			On(
				"UpsertUserPreferences",
				context.Background(),
				mock.MatchedBy(func(p preference.UserEmissionPreference) bool {
					return p.CountryCode == "IDN"
				}),
			).
			Return(expected, nil)

		result, err := uc.SetUserPreferences(
			context.Background(),
			1,
			"",
			&customLimit,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "IDN", result.CountryCode)

		mockPreferenceRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {

		mockPreferenceRepo := new(inmocks.MockPreferenceRepository)

		uc := usecase.NewPreferenceUseCase(
			mockPreferenceRepo,
		)

		customLimit := 10.5

		mockPreferenceRepo.
			On(
				"UpsertUserPreferences",
				context.Background(),
				mock.Anything,
			).
			Return(nil, assert.AnError)

		result, err := uc.SetUserPreferences(
			context.Background(),
			1,
			"IDN",
			&customLimit,
		)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockPreferenceRepo.AssertExpectations(t)
	})
}
