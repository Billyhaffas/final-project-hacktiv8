package emission

import (
	"time"
)

type EmissionOrigin struct {
	UserId        int32     `json:"user_id"`
	VehicleType   string    `json:"vehicle_type"`
	FuelType      string    `json:"fuel_type"`
	DistanceKm    float64   `json:"distance_km"`
	EmissionKgCo2 float64   `json:"emission_kg_co2"`
	RecordedAt    time.Time `json:"recorded_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type EmissionBody struct {
	UserId      int32   `json:"user_id"`
	VehicleType string  `json:"vehicle_type"`
	FuelType    string  `json:"fuel_type"`
	DistanceKm  float64 `json:"distance_km"`
}
