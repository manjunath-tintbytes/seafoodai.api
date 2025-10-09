package models

import "time"

type Region struct {
	ID        uint    `gorm:"primaryKey"`
	Region    string  `gorm:"not null;index"`
	Quota     float64 `gorm:"type:numeric(12,2)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
