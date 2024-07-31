package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	ID               int
	TagFeatureBanner []TagFeatureBanner `json:"-" gorm:"foreignKey:TagID"`
}
