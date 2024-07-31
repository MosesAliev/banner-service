package models

import (
	"time"
)

type TagFeatureBanner struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	TagID     int `gorm:"primaryKey"`
	FeatureID int `gorm:"primaryKey"`
	BannerID  uint
}
