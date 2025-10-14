package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupParticipant struct {
    ID       uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
    GroupID  uuid.UUID  `gorm:"not null" json:"group_id"`
    UserID   uuid.UUID  `gorm:"not null" json:"user_id"`
    gorm.Model          `json:"-"`
}

func (u *GroupParticipant) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}