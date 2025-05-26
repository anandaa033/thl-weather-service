package models

import "time"

type UserCity struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"index:user_city_idx,unique;not null"`
	City      string    `gorm:"index:user_city_idx,unique;not null"`
	CityOrder int       `gorm:"not null;default:0" json:"city_order"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt *time.Time
}

func (UserCity) TableName() string {
	return "cities_users"
}
