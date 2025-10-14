package handler

import (
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthInterface interface {
	Register(user *dto.RegisterRequest) (*model.User, error)
	Login(email string, password string) (string, *model.User, error)
}

type AuthHandler struct {
	uc AuthInterface
}

func NewAuthHandler(uc AuthInterface) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.Password != req.ConfirmPassword {
		responses.Error(c, http.StatusBadRequest, "Password and Confirm Password do not match")
		return
	}

	registered, err := h.uc.Register(&req)
	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Failed to register user: "+err.Error())
		return
	}

	token, _, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		responses.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	response := dto.RegisterResponse{
		ID: registered.ID,
		Role: registered.Role,
		Name: registered.Name,
		Username: registered.Username,
		Email: registered.Email,
		PhoneNumber: registered.PhoneNumber,
		Token: token,
	}

	responses.Created(c, "Register successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	token, user, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		responses.Error(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	responses.Success(c, "Login successfully", gin.H{"userID": user.ID, "token": token })
}