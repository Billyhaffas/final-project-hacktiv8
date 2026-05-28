package preference

import "time"

type UserEmissionPreference struct {
	Id                   uint64     `gorm:"column:id;primaryKey"`
	UserId               int32      `gorm:"column:user_id;uniqueIndex"`
	CountryCode          string     `gorm:"column:country_code"`
	CustomDailyLimitKgCo2 *float64  `gorm:"column:custom_daily_limit_kg_co2"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at"`
}

func (UserEmissionPreference) TableName() string {
	return "user_emission_preferences"
}
