package domain

import (
	"context"
	"time"
)

type Emission struct {
	Id              string    `bson:"_id" json:"id"`
	CountryCode     string    `bson:"country_code" json:"country_code"`
	DailyLimitKgCo2 float64   `bson:"daily_limit_kg_co2" json:"daily_limit_kg_co2"`
	SourceUrl       string    `bson:"source_url" json:"source_url"`
	UpdatedAt       time.Time `bson:"updated_at" json:"updated_at"`
}

type EmissionRepository interface {
	BulkInsert(ctx context.Context, emissions []Emission) (bool, error)
}

type EmissionProvider interface {
	FetchEmissions(ctx context.Context) ([]Emission, error)
}
