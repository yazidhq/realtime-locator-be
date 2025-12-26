package dto

import (
	"time"

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

type LocationHistoryItemResponse struct {
	ID        string    `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
}

type LocationHistoryGroupResponse struct {
	Date      string                        `json:"date"`
	Locations []LocationHistoryItemResponse `json:"locations"`
}