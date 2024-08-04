package models

import (
	"time"
)

type Banner struct {
	CreatedAt        time.Time          `json:"-"`
	UpdatedAt        time.Time          `json:"-"`
	DeletedAt        time.Time          `json:"-"`
	ID               int                `json:"banner_id,omitempty"`
	TagIDs           []int              `json:"tag_ids" gorm:"-"`
	Tags             []Tag              `json:"-" gorm:"many2many:tag_ids"`
	FeatureID        int                `json:"feature_id"`
	Content          Content            `gorm:"embedded" json:"content"`
	IsActive         bool               `json:"is_active"`
	TagFeatureBanner []TagFeatureBanner `json:"-" gorm:"foreignKey:BannerID"`
}
