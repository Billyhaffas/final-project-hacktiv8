package domain

import (
	"context"
	"time"
)

type CarbonPrice struct {
	Id             string    `bson:"_id" json:"id"`
	PricePerTonUsd float64   `bson:"price_per_ton_usd" json:"price_per_ton_usd"`
	UsdCurRate     float64   `bson:"usd_cur_rate" json:"usd_cur_rate"`
	Source         string    `bson:"source" json:"source"`
	FetchedAt      time.Time `bson:"fetched_at" json:"fetched_at"`
}

type CarbonPriceRepository interface {
	BulkInsert(ctx context.Context, carbonPrices []CarbonPrice) (bool, error)
}

type CarbonPriceProvider interface {
	FetchPrices(ctx context.Context) ([]CarbonPrice, error)
}
