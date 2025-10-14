package dto

import (
	"TeamTrackerBE/internal/domain/model"

	"github.com/google/uuid"
)

type ContactCreateRequest struct {
	UserID 	  uuid.UUID 			`json:"user_id" binding:"required"`
	ContactID uuid.UUID 			`json:"contact_id" binding:"required"`
	Status    model.ContactStatus   `json:"name" binding:"required"`
}

type ContactUpdateRequest struct {
	UserID 	  uuid.UUID 			`json:"user_id"`
	ContactID uuid.UUID 			`json:"contact_id"`
	Status    model.ContactStatus   `json:"name"`
}

type ContactResponse struct {
	ID   string `json:"id,omitempty"`
	UserID 	  uuid.UUID 			`json:"user_id"`
	ContactID uuid.UUID 			`json:"contact_id"`
	Status    model.ContactStatus   `json:"name"`
}