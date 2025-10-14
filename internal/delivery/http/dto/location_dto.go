package dto

import (
	"github.com/google/uuid"
)

type LocationCreateRequest struct {
	GroupID 	uuid.UUID	`json:"group_id" binding:"required"`
	UserID		uuid.UUID	`json:"user_id" binding:"required"`
	Latitude	float64		`json:"latitude" binding:"required"`
	Longitude	float64		`json:"longitude" binding:"required"`
}

type LocationUpdateRequest struct {
	GroupID		uuid.UUID	`json:"group_id"`
	UserID		uuid.UUID	`json:"user_id"`
	Latitude    float64   	`json:"latitude"`
	Longitude   float64   	`json:"longitude"`
}

type LocationResponse struct {
	ID   		string 		`json:"id,omitempty"`
	GroupID 	uuid.UUID	`json:"group_id"`
	UserID 	  	uuid.UUID	`json:"user_id"`
	Latitude    float64   	`json:"latitude"`
	Longitude   float64   	`json:"longitude"`
}