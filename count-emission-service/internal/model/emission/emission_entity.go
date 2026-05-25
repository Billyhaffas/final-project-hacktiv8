package emission

import "time"

type Emission struct {
	Id            uint64
	UserId        int8
	VehicleType   string
	FuelType      string
	DistanceKm    float64
	EmissionKgCo2 float64
	RecordedAt    time.Time
	CreatedAt     time.Time
}

type UserDailyEmission struct {
	UserId             int8
	Date               string
	TotalEmissionKgCo2 float64
}

type UserMonthlyEmissionDetail struct {
	UserId             int8
	Month              string
	TotalEmissionKgCo2 float64
}

type UserMonthlyEmission struct {
	UserId                    int8
	DailyEmissions            []UserDailyEmission
	TotalEmissionMonthlyKgCo2 float64
}

type UserYearlyEmission struct {
	UserId                   int8
	MonthlyEmissions         []UserMonthlyEmissionDetail
	TotalEmissionYearlyKgCo2 float64
}
