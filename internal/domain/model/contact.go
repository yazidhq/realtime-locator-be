package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contact struct {
    ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
    UserID      uuid.UUID       `gorm:"not null" json:"user_id"`
    ContactID   uuid.UUID       `gorm:"not null" json:"contact_id"`
    Status      ContactStatus   `gorm:"not null;default:'pending'" json:"name"`
    gorm.Model                  `json:"-"`
}

func (u *Contact) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}