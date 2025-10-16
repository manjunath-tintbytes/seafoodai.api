package models

import "time"

type Landing struct {
	ID            uint    `gorm:"primaryKey"`
	Year          int     `gorm:"not null" json:"year"`
	LandingPortID uint    `gorm:"not null" json:"landing_port_id"`
	LandingNameID uint    `gorm:"not null" json:"landing_name_id"`
	Pounds        float64 `gorm:"type:numeric(14,2)" json:"pounds"`
	Dollars       float64 `gorm:"type:numeric(14,2)" json:"dollars"`
	MetricTons    float64 `gorm:"type:numeric(14,2)" json:"metric_tons"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `gorm:"index"`

	// Associations
	LandingPort LandingPort `gorm:"foreignKey:LandingPortID"`
	LandingName LandingName `gorm:"foreignKey:LandingNameID"`
}
