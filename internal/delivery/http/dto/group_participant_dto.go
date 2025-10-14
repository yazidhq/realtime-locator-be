package dto

import (
	"github.com/google/uuid"
)

type GroupParticipantCreateRequest struct {
	GroupID uuid.UUID `json:"group_id" binding:"required"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
}

type GroupParticipantUpdateRequest struct {
	GroupID uuid.UUID `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}

type GroupParticipantResponse struct {
	ID   	string 	  `json:"id,omitempty"`
	GroupID uuid.UUID `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}