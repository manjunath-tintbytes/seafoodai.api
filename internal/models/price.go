package models

import "time"

type Price struct {
	ID        uint      `gorm:"primaryKey"`
	SeafoodID uint      `gorm:"not null;index"`
	Date      time.Time `gorm:"default:now();index"`
	Price     float64   `gorm:"type:numeric(12,2);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
