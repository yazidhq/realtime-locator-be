package dto

import (
	"TeamTrackerBE/internal/domain/model"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name            string `json:"name" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegisterResponse struct {
	ID       	 uuid.UUID  `json:"id"`
	Role     	 model.Role `json:"role"`
	Name     	 string     `json:"name"`
	Username 	 string     `json:"username"`
	Email    	 string     `json:"email"`
	PhoneNumber	 string     `json:"phone_number"`
	Token    	 string     `json:"token"`
	RefreshToken string     `json:"refresh_token"`
}

type LoginResponse struct {
    ID       	 uuid.UUID  `json:"id"`
    Token        string     `json:"token"`
    RefreshToken string     `json:"refresh_token"`
}

type RefreshTokenResponse struct {
    Token        string `json:"token"`
    RefreshToken string `json:"refresh_token"`
}