package domain

import "context"

type Emission struct {
	Id          string  `bson:"Id" json:"id"`
	Entity      string  `bson:"entity" json:"entity"`
	Code        string  `bson:"code" json:"code"`
	Year        int     `bson:"year" json:"year"`
	CarbonPrice float64 `bson:"carbon_price" json:"carbon_price"`
}

type EmissionRepository interface {
	BulkInsert(ctx context.Context, emissions []Emission) (bool, error)
}

type EmissionProvider interface {
	FetchEmissions(ctx context.Context) ([]Emission, error)
}
