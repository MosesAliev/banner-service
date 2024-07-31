package models

import (
	"gorm.io/gorm"
)

type Banner struct {
	gorm.Model
	TagIDs           []int              `json:"tag_ids" gorm:"-"`
	Tags             []Tag              `json:"-" gorm:"many2many:tag_ids"`
	FeatureID        int                `json:"feature_id"`
	Content          Content            `gorm:"embedded" json:"content"`
	IsActive         bool               `json:"is_active"`
	TagFeatureBanner []TagFeatureBanner `json:"-" gorm:"foreignKey:BannerID"`
}
