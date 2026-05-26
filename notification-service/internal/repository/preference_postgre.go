package repository

import (
	"context"
	"notification-service/internal/domain"

	"gorm.io/gorm"
)

type preferenceRepo struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) domain.PreferenceRepository {
	return &preferenceRepo{db: db}
}

func (r *preferenceRepo) GetAllUserIDs(ctx context.Context) ([]int, error) {
	var userIDs []int
	err := r.db.WithContext(ctx).Model(&domain.UserEmissionPreference{}).Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (r *preferenceRepo) GetByUserID(ctx context.Context, userID int) (*domain.UserEmissionPreference, error) {
	var pref domain.UserEmissionPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	if err != nil {
		return nil, err
	}
	return &pref, nil
}

func (r *preferenceRepo) Create(ctx context.Context, pref *domain.UserEmissionPreference) error {
	return r.db.WithContext(ctx).Create(pref).Error
}

func (r *preferenceRepo) Update(ctx context.Context, pref *domain.UserEmissionPreference) error {
	// Updates all fields except primary key and creation time
	return r.db.WithContext(ctx).Model(pref).Updates(map[string]interface{}{
		"country_code":              pref.CountryCode,
		"custom_daily_limit_kg_co2": pref.CustomDailyLimitKgCo2,
	}).Error
}

func (r *preferenceRepo) Delete(ctx context.Context, userID int) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&domain.UserEmissionPreference{}).Error
}
