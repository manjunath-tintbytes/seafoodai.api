package models

import "time"

type Seafood struct {
	ID          uint      `gorm:"primaryKey"`
	SpeciesID   uint      `gorm:"not null;index"`
	Species     Species   `gorm:"foreignKey:SpeciesID"`
	RegionID    uint      `gorm:"not null;index"`
	Region      Region    `gorm:"foreignKey:RegionID"`
	SubRegionID *uint     `gorm:"index"`
	SubRegion   SubRegion `gorm:"foreignKey:SubRegionID"`
	CategoryID  uint      `gorm:"not null;index"`
	Category    Category  `gorm:"foreignKey:CategoryID"`
	VolumeMT    float64   `gorm:"type:numeric(12,2)"`
	PriceUnit   string    `gorm:"size:10"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}
