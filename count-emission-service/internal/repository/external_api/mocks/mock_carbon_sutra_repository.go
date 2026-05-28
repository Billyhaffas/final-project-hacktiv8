package mocks

import (
	"count-emission-service/internal/model/thirdparty/carbonsutra"

	"github.com/stretchr/testify/mock"
)

type MockCarbonSutraRepository struct {
	mock.Mock
}

func (m *MockCarbonSutraRepository) GetCarbonEmission(
	payload carbonsutra.CountEmisionBodyPayload,
) (*carbonsutra.CountEmisionThirdParty, error) {

	args := m.Called(payload)

	var result *carbonsutra.CountEmisionThirdParty

	if args.Get(0) != nil {
		result = args.Get(0).(*carbonsutra.CountEmisionThirdParty)
	}

	return result, args.Error(1)
}
