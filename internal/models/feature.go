package models

import (
	"gorm.io/gorm"
)

type Feature struct {
	gorm.Model
	ID               int
	Banners          []Banner           `gorm:"foreignKey:FeatureID"`
	TagFeatureBanner []TagFeatureBanner `json:"-" gorm:"foreignKey:FeatureID"`
}
