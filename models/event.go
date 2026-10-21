package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	ImageID     string    `json:"imageId"`
	Image       string    `json:"image"`
	UserId      int       `json:"userId"`
	User        User      `gorm:"foreignKey:userId" json:"-"` //untuk relasi ke table / model user
	DateTIME    time.Time `json:"datetime" binding:"required"`
}
