package models

import "time"

type Species struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
