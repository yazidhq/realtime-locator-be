package dto

type UserCreateRequest struct {
	Role        string `json:"role" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Role        string `json:"role"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserResponse struct {
	ID          string `json:"id,omitempty"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}