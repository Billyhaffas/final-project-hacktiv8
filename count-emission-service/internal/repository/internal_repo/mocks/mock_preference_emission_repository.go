package mocks

import (
	"context"

	"count-emission-service/internal/model/preference"

	"github.com/stretchr/testify/mock"
)

type MockPreferenceRepository struct {
	mock.Mock
}

func (m *MockPreferenceRepository) GetUserPreferences(ctx context.Context, userID int32) (*preference.UserEmissionPreference, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*preference.UserEmissionPreference), args.Error(1)
}

func (m *MockPreferenceRepository) UpsertUserPreferences(ctx context.Context, pref preference.UserEmissionPreference) (*preference.UserEmissionPreference, error) {
	args := m.Called(ctx, pref)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*preference.UserEmissionPreference), args.Error(1)
}
