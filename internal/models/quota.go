package models

import "time"

type Quota struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Date           time.Time `gorm:"not null" json:"date"`
	ProductName    string    `gorm:"type:varchar(255);not null" json:"product_name"`
	RemainingQuota float64   `gorm:"type:decimal(15,2);not null" json:"remaining_quota"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time `gorm:"index"`
}
