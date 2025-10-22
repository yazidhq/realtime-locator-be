package dto

import (
	"github.com/google/uuid"
)

type LocationCreateRequest struct {
	UserID		uuid.UUID	`json:"user_id" binding:"required"`
	Latitude	float64		`json:"latitude" binding:"required"`
	Longitude	float64		`json:"longitude" binding:"required"`
}

type LocationUpdateRequest struct {
	UserID		uuid.UUID	`json:"user_id"`
	Latitude    float64   	`json:"latitude"`
	Longitude   float64   	`json:"longitude"`
}

type LocationResponse struct {
	ID   		string 		`json:"id,omitempty"`
	UserID 	  	uuid.UUID	`json:"user_id"`
	Latitude    float64   	`json:"latitude"`
	Longitude   float64   	`json:"longitude"`
}