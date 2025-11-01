package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
    ID          uuid.UUID   `gorm:"type:uuid;primaryKey" json:"id"`
    UserID      uuid.UUID   `gorm:"not null" json:"user_id"`
    Latitude    float64     `gorm:"not null" json:"latitude"`
    Longitude   float64     `gorm:"not null" json:"longitude"`
    gorm.Model              `json:"-"`
}

func (u *Location) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}