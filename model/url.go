package model

import "gorm.io/gorm"

type Url struct {
	gorm.Model
	ShortUrl string `gorm:"unique;not null" json:"short_url"`
	LongUrl  string `gorm:"not null" json:"long_url"`
}
