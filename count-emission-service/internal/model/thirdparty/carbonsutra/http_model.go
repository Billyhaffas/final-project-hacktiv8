package carbonsutra

type CountEmisionBodyPayload struct {
	VehicleType   string  `bson:"vehicle_type" json:"vehicle_type"`
	FuelType      string  `bson:"fuel_type" json:"fuel_type"`
	DistanceValue float64 `bson:"distance_value" json:"distance_value"`
	DistanceUnit  string  `bson:"distance_unit" json:"distance_unit"`
	IncludeWtt    string  `bson:"include_wtt" json:"include_wtt"`
}

type CountEmisionThirdParty struct {
	Data    EmissionData `json:"data"`
	Success bool         `json:"success"`
	Status  int          `json:"status"`
}

type EmissionData struct {
	Type          string  `json:"type"`
	VehicleType   string  `json:"vehicle_type"`
	FuelType      string  `json:"fuel_type"`
	DistanceValue float64 `json:"distance_value"`
	DistanceUnit  string  `json:"distance_unit"`
	IncludeWtt    string  `json:"include_wtt"`
	Co2eGm        float64 `json:"co2e_gm"`
	Co2eKg        float64 `json:"co2e_kg"`
	Co2eMt        float64 `json:"co2e_mt"`
	Co2eLb        float64 `json:"co2e_lb"`
}
