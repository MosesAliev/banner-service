package models

import "time"

type Feature struct {
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
	DeletedAt        time.Time `json:"-"`
	ID               int
	Banners          []Banner           `gorm:"foreignKey:FeatureID"`
	TagFeatureBanner []TagFeatureBanner `json:"-" gorm:"foreignKey:FeatureID"`
}
