package models

import "time"

type LandingPort struct {
	ID         uint   `gorm:"primaryKey"`
	RegionName string `gorm:"type:varchar(100);not null" json:"region_name"`
	Port       string `gorm:"type:varchar(100)" json:"port"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}
