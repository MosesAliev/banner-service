package models

import "time"

type Tag struct {
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
	DeletedAt        time.Time `json:"-"`
	ID               int
	TagFeatureBanner []TagFeatureBanner `json:"-" gorm:"foreignKey:TagID"`
}
