package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    Role        Role      `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
    Name        string    `gorm:"size:100;not null" json:"name"`
    Username    string    `gorm:"size:100;not null" json:"username"`
    Email       string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
    PhoneNumber string    `gorm:"uniqueIndex;size:100;not null" json:"phone_number"`
    Password    string    `gorm:"size:255;not null" json:"password"`
    gorm.Model            `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}