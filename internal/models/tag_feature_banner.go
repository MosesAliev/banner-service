package models

import (
	"time"
)

type TagFeatureBanner struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
	TagID     int       `gorm:"primaryKey"`
	FeatureID int       `gorm:"primaryKey"`
	BannerID  int
}
