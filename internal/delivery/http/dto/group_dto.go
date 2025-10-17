package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type GroupCreateRequest struct {
	Name       string    		`json:"name" binding:"required"`
	OwnerID    uuid.UUID 		`json:"owner_id" binding:"required"`
	RadiusArea json.RawMessage  `json:"radius_area"`
}

type GroupUpdateRequest struct {
	Name       string    		`json:"name"`
	OwnerID    uuid.UUID 		`json:"owner_id"`
	RadiusArea json.RawMessage  `json:"radius_area"`
}

type GroupResponse struct {
	ID   	   string    	   `json:"id,omitempty"`
	Name       string    	   `json:"name"`
	OwnerID    uuid.UUID 	   `json:"owner_id"`
	RadiusArea json.RawMessage `json:"radius_area"`
}