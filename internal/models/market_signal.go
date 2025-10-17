package models

import (
	"time"
)

type MarketSignal struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_title_published" json:"title"`
	PublishedDate time.Time `gorm:"uniqueIndex:idx_title_published" json:"published_date"`
	Author        string    `gorm:"type:varchar(255)" json:"author"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `gorm:"index"`
}
