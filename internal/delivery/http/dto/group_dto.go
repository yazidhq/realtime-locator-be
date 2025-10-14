package dto

import (
	"github.com/google/uuid"
)

type GroupCreateRequest struct {
	Name    string    `json:"name" binding:"required"`
	OwnerID uuid.UUID `json:"owner_id" binding:"required"`
}

type GroupUpdateRequest struct {
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id"`
}

type GroupResponse struct {
	ID   	string 	  `json:"id,omitempty"`
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id"`
}