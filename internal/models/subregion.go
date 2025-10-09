package models

import "time"

type SubRegion struct {
	ID        uint   `gorm:"primaryKey"`
	RegionID  uint   `gorm:"not null;index"`
	Region    Region `gorm:"foreignKey:RegionID"`
	SubRegion string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}
