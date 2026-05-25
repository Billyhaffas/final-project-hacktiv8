package repository

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/emission"

	"gorm.io/gorm"
)

type EmissionRepository struct {
	DB *gorm.DB
}

func NewEmissionCollection(db *gorm.DB) domain.EmissionRepository {
	return &EmissionRepository{DB: db}
}

func (emr *EmissionRepository) CreateUserEmission(ctx context.Context, req emission.EmissionOrigin) error {
	err := emr.DB.WithContext(ctx).
		Table("emissions").
		Create(&req).Error
	if err != nil {
		return err
	}
	return nil
}
