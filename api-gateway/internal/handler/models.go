package handler

// --- Request bodies ---

type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
	Name     string `json:"name" example:"John Doe"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" example:"user@example.com"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" example:"reset-token-here"`
	Password string `json:"password" example:"newpassword123"`
}

type LogEmissionRequest struct {
	VehicleType string  `json:"vehicle_type" example:"Car-Size-Medium"`
	FuelType    string  `json:"fuel_type" example:"Petrol"`
	DistanceKm  float64 `json:"distance_km" example:"15.5"`
}

type UpdatePreferencesRequest struct {
	CountryCode           string  `json:"country_code" example:"IDN"`
	CustomDailyLimitKgCo2 float64 `json:"custom_daily_limit_kg_co2" example:"5.0"`
}

// --- Response envelopes ---

type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorInfo `json:"error"`
}

type ErrorInfo struct {
	Code    string `json:"code" example:"VALIDATION_ERROR"`
	Message string `json:"message" example:"vehicle_type is required"`
}

// --- Response data shapes ---

type EmissionData struct {
	Message string `json:"message" example:"Emission has been created"`
}

type DailyTotalData struct {
	UserId             int32   `json:"user_id" example:"1"`
	Date               string  `json:"date" example:"2026-05-28"`
	TotalEmissionKgCo2 float64 `json:"total_emission_kg_co2" example:"3.5"`
}

type AlertData struct {
	IsExceeded      bool    `json:"is_exceeded" example:"false"`
	DailyTotalKg    float32 `json:"daily_total_kg" example:"3.5"`
	DailyLimitKg    float32 `json:"daily_limit_kg" example:"6.3"`
	ThresholdSource string  `json:"threshold_source" example:"country"`
	Message         string  `json:"message" example:"You are within today's limit"`
}

type ConvertData struct {
	EmissionKgCo2     float64 `json:"emission_kg_co2" example:"100.0"`
	PricePerTonUsd    float64 `json:"price_per_ton_usd" example:"23.0"`
	ExchangeRateUsdIdr float64 `json:"exchange_rate_usd_idr" example:"16250.0"`
	TotalIdr          float64 `json:"total_idr" example:"37375.0"`
}

type PreferencesData struct {
	UserId                int32   `json:"user_id" example:"1"`
	CountryCode           string  `json:"country_code" example:"IDN"`
	CustomDailyLimitKgCo2 float32 `json:"custom_daily_limit_kg_co2" example:"5.0"`
}
