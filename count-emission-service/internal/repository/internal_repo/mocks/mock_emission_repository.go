package mocks

import (
	"context"

	"count-emission-service/internal/model/emission"

	"github.com/stretchr/testify/mock"
)

type MockEmissionRepository struct {
	mock.Mock
}

func (m *MockEmissionRepository) CreateUserEmission(ctx context.Context, req emission.EmissionOrigin) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockEmissionRepository) GetUserDailyEmission(ctx context.Context, userId int32) (*emission.UserDailyEmission, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*emission.UserDailyEmission), args.Error(1)
}

func (m *MockEmissionRepository) GetUserMonthlyEmission(ctx context.Context, userId int32) (*emission.UserMonthlyEmission, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*emission.UserMonthlyEmission), args.Error(1)
}

func (m *MockEmissionRepository) GetUserYearlyEmission(ctx context.Context, userId int32) (*emission.UserYearlyEmission, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*emission.UserYearlyEmission), args.Error(1)
}

func (m *MockEmissionRepository) GetDailyTotal(ctx context.Context, userId int32, date string) (float64, int32, error) {
	args := m.Called(ctx, userId, date)
	return args.Get(0).(float64), args.Get(1).(int32), args.Error(2)
}
