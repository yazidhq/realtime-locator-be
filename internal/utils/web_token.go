package utils

import (
	"TeamTrackerBE/internal/config"
	"TeamTrackerBE/internal/domain/model"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(userID uuid.UUID, role model.Role) (string, error) {
	expStr := config.Env.JWT.ExpiresIn
	duration, err := time.ParseDuration(expStr)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRES_IN format: %v", err)
	}

	claims := jwt.MapClaims{
		"userID": userID,
		"role": role,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.JWT.SecretKey))
}