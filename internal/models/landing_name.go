package models

import "time"

type LandingName struct {
	ID             uint      `gorm:"primaryKey"`
	NMFSName       string    `gorm:"type:varchar(150);not null" json:"nmfs_name"`
	ScientificName string    `gorm:"type:varchar(150)" json:"scientific_name"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}
