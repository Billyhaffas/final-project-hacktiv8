package repository

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/preference"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PreferenceRepository struct {
	DB *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) domain.PreferenceRepository {
	return &PreferenceRepository{DB: db}
}

func (r *PreferenceRepository) GetUserPreferences(ctx context.Context, userID int32) (*preference.UserEmissionPreference, error) {
	var pref preference.UserEmissionPreference
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&pref).Error
	if err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *PreferenceRepository) UpsertUserPreferences(ctx context.Context, pref preference.UserEmissionPreference) (*preference.UserEmissionPreference, error) {
	err := r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"country_code", "custom_daily_limit_kg_co2", "updated_at"}),
		}).
		Create(&pref).Error
	if err != nil {
		return nil, err
	}
	return &pref, nil
}
