package domain

import (
	"context"
	"time"
)

type PreferenceUpsertInput struct {
	UserID                int     `json:"user_id" validate:"required,gt=0"`
	CountryCode           string  `json:"country_code" validate:"required,len=3"`
	CustomDailyLimitKgCo2 float64 `json:"custom_daily_limit_kg_co2" validate:"required,gt=0"`
}

type UserEmissionPreference struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                int       `gorm:"column:user_id;not null" json:"user_id"`
	CountryCode           string    `gorm:"column:country_code;type:varchar(3);not null" json:"country_code"`
	CustomDailyLimitKgCo2 float64   `gorm:"column:custom_daily_limit_kg_co2;type:numeric" json:"custom_daily_limit_kg_co2"`
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

type PreferenceRepository interface {
	GetByUserID(ctx context.Context, userID int) (*UserEmissionPreference, error)

	Create(ctx context.Context, pref *UserEmissionPreference) error
	Update(ctx context.Context, pref *UserEmissionPreference) error
	Delete(ctx context.Context, userID int) error

	GetAllUserIDs(ctx context.Context) ([]int, error)
}

type MasterLimitRepository interface {
	GetDefaultLimitByCountry(ctx context.Context, countryCode string) (float64, error)
}

type PreferenceUsecase interface {
	GetPreference(ctx context.Context, userID int) (*UserEmissionPreference, error)
	SavePreference(ctx context.Context, input PreferenceUpsertInput) (*UserEmissionPreference, error)
	DeletePreference(ctx context.Context, userID int) error
}
