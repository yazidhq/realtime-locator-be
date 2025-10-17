package model

import (
	"encoding/json"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
    ID          uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
    Name        string          `gorm:"not null" json:"name"`
    OwnerID     uuid.UUID       `gorm:"not null" json:"owner_id"`
    RadiusArea  json.RawMessage `gorm:"type:jsonb;null" json:"radius_area"`
    gorm.Model                  `json:"-"`
}

func (u *Group) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    return
}