package handler

import (
	ws "TeamTrackerBE/internal/delivery/event/websocket"
	"TeamTrackerBE/internal/delivery/http/dto"
	"TeamTrackerBE/internal/domain/model"
	"TeamTrackerBE/internal/utils/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthInterface interface {
    Register(user *dto.RegisterRequest) (*model.User, error)
    Login(email string, password string) (string, string, *model.User, error)
    RefreshToken(refreshToken string) (string, string, error)
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

    accessToken, refreshToken, _, err := h.uc.Login(req.Email, req.Password)
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
        AccessToken: accessToken,
        RefreshToken: refreshToken,
    }

    responses.Created(c, "Register successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }

    accessToken, refreshToken, user, err := h.uc.Login(req.Email, req.Password)
    if err != nil {
        responses.Error(c, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    response := dto.LoginResponse{
        ID: user.ID,
        Name: user.Name,
        Email: user.Email,
        PhoneNumber: user.PhoneNumber,
        AccessToken: accessToken,
        RefreshToken: refreshToken,
    }

    responses.Success(c, "Login successfully", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
    var req dto.RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        responses.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }

	if req.RefreshToken == "" {
        responses.Error(c, http.StatusBadRequest, "Refresh token is required")
        return
    }

    newAccessToken, newRefreshToken, err := h.uc.RefreshToken(req.RefreshToken)
    if err != nil {
        responses.Error(c, http.StatusUnauthorized, "Invalid or expired refresh token: "+err.Error())
        return
    }

    response := dto.RefreshTokenResponse{
        AccessToken: newAccessToken,
        RefreshToken: newRefreshToken,
    }

    responses.Success(c, "Token refreshed successfully", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
    uidVal, exists := c.Get("userID")
    if !exists {
        responses.Error(c, http.StatusUnauthorized, "token user not found")
        return
    }

    userID, ok := uidVal.(uuid.UUID)
    if !ok {
        responses.Error(c, http.StatusUnauthorized, "invalid user id in token context")
        return
    }

    ws.GetHub().ForceOffline(userID)
    responses.Success(c, "Logout successfully", map[string]any{"user_id": userID.String()})
}