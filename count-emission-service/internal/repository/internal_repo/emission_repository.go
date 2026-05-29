package repository

import (
	"context"
	"count-emission-service/internal/domain"
	"count-emission-service/internal/model/emission"
	"sync"
	"time"

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

func (emr *EmissionRepository) GetUserDailyEmission(ctx context.Context, userId int32) (*emission.UserDailyEmission, error) {
	var userEmission emission.UserDailyEmission
	userEmission.UserId = userId
	err := emr.DB.WithContext(ctx).
		Table("emissions").
		Where(
			"user_id = ? AND DATE(recorded_at) = ?",
			userId,
			time.Now().Format("2006-01-02"),
		).
		Select("COALESCE(SUM(emission_kg_co2), 0)").
		Scan(&userEmission.TotalEmissionKgCo2).Error
	userEmission.Date = time.Now().Format("2006-01-02")
	if err != nil {
		return nil, err
	}
	return &userEmission, nil
}

func (emr *EmissionRepository) GetDailyTotal(ctx context.Context, userId int32, date string) (float64, int32, error) {
	var total float64
	var cnt int64

	if err := emr.DB.WithContext(ctx).Table("emissions").
		Where("user_id = ? AND DATE(recorded_at) = ?", userId, date).
		Select("COALESCE(SUM(emission_kg_co2), 0)").
		Scan(&total).Error; err != nil {
		return 0, 0, err
	}
	if err := emr.DB.WithContext(ctx).Table("emissions").
		Where("user_id = ? AND DATE(recorded_at) = ?", userId, date).
		Count(&cnt).Error; err != nil {
		return 0, 0, err
	}
	return total, int32(cnt), nil
}

func (emr *EmissionRepository) GetUserMonthlyEmission(ctx context.Context, userId int32) (*emission.UserMonthlyEmission, error) {
	var (
		result              []emission.UserDailyEmission
		monthlyTotal        float64
		dailyErr            error
		monthlyErr          error
		userMonthlyEmission emission.UserMonthlyEmission
		wg                  sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()

		dailyErr = emr.DB.WithContext(ctx).
			Table("emissions").
			Select(`
			user_id,
			TO_CHAR(recorded_at, 'YYYY-MM-DD') as date,
			SUM(emission_kg_co2) as total_emission_kg_co2
		`).
			Where(`
			user_id = ?
			AND DATE_TRUNC('month', recorded_at) =
			    DATE_TRUNC('month', CURRENT_DATE)
		`, userId).
			Group(`
			user_id,
			DATE(recorded_at)
		`).
			Order("date ASC").
			Scan(&result).Error
	}()

	go func() {
		defer wg.Done()

		monthlyErr = emr.DB.WithContext(ctx).
			Table("emissions").
			Select("SUM(emission_kg_co2)").
			Where(`
			user_id = ?
			AND DATE_TRUNC('month', recorded_at) =
			    DATE_TRUNC('month', CURRENT_DATE)
		`, userId).
			Scan(&monthlyTotal).Error
	}()

	wg.Wait()

	if dailyErr != nil {
		return nil, dailyErr
	}

	if monthlyErr != nil {
		return nil, monthlyErr
	}

	userMonthlyEmission.UserId = userId
	userMonthlyEmission.DailyEmissions = result
	userMonthlyEmission.TotalEmissionMonthlyKgCo2 = monthlyTotal
	return &userMonthlyEmission, nil
}

func (emr *EmissionRepository) GetUserYearlyEmission(ctx context.Context, userId int32) (*emission.UserYearlyEmission, error) {
	var (
		result             []emission.UserMonthlyEmissionDetail
		yearlyTotal        float64
		monthlyErr         error
		yearlyErr          error
		userYearlyEmission emission.UserYearlyEmission
		wg                 sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()

		monthlyErr = emr.DB.WithContext(ctx).
			Table("emissions").
			Select(`
				user_id,
				TRIM(TO_CHAR(DATE_TRUNC('month', recorded_at), 'Month')) as month,
				SUM(emission_kg_co2) as total_emission_kg_co2
			`).
			Where(`
				user_id = ?
				AND DATE_TRUNC('year', recorded_at) =
				    DATE_TRUNC('year', CURRENT_DATE)
			`, userId).
			Group(`
				user_id,
				DATE_TRUNC('month', recorded_at)
			`).
			Order(`
				DATE_TRUNC('month', recorded_at) ASC
			`).
			Scan(&result).Error
	}()

	go func() {
		defer wg.Done()

		yearlyErr = emr.DB.WithContext(ctx).
			Table("emissions").
			Select("SUM(emission_kg_co2)").
			Where(`
				user_id = ?
				AND DATE_TRUNC('year', recorded_at) =
				    DATE_TRUNC('year', CURRENT_DATE)
			`, userId).
			Scan(&yearlyTotal).Error
	}()

	wg.Wait()

	if monthlyErr != nil {
		return nil, monthlyErr
	}

	if yearlyErr != nil {
		return nil, yearlyErr
	}

	userYearlyEmission.UserId = userId
	userYearlyEmission.MonthlyEmissions = result
	userYearlyEmission.TotalEmissionYearlyKgCo2 = yearlyTotal

	return &userYearlyEmission, nil
}
