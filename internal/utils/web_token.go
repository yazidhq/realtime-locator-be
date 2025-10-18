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
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.JWT.SecretKey))
}

func GenerateRefreshToken(userID uuid.UUID, role model.Role) (string, error) {
	expStr := config.Env.JWT.RefreshTokenExpiresIn
	duration, err := time.ParseDuration(expStr)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_REFRESH_TOKEN_EXPIRES_IN format: %v", err)
	}
    
    claims := jwt.MapClaims{
        "userID": userID.String(),
		"role": role,
        "exp": time.Now().Add(duration).Unix(),
        "iat": time.Now().Unix(),
        "type": "refresh",
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.Env.JWT.SecretKey))
}

func ValidateRefreshToken(tokenString string) (uuid.UUID, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(config.Env.JWT.SecretKey), nil
    })

    if err != nil || !token.Valid {
        return uuid.Nil, fmt.Errorf("invalid or expired refresh token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, fmt.Errorf("invalid token claims")
    }

    tokenType, ok := claims["type"].(string)
    if !ok || tokenType != "refresh" {
        return uuid.Nil, fmt.Errorf("invalid token type")
    }

    if exp, ok := claims["exp"].(float64); ok {
        if time.Unix(int64(exp), 0).Before(time.Now()) {
            return uuid.Nil, fmt.Errorf("refresh token expired")
        }
    }

    userIDStr, ok := claims["userID"].(string)
    if !ok {
        return uuid.Nil, fmt.Errorf("invalid userID in token")
    }

    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid UUID format")
    }

    return userID, nil
}

func GenerateGroupInviteToken(groupID uuid.UUID) (string, error) {
    expStr := config.Env.JWT.ExpiresIn
	duration, err := time.ParseDuration(expStr)
    if err != nil {
        return "", fmt.Errorf("invalid invite expires duration: %v", err)
    }

    claims := jwt.MapClaims{
        "groupID": groupID.String(),
        "type":    "group_invite",
        "exp":     time.Now().Add(duration).Unix(),
        "iat":     time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(config.Env.JWT.SecretKey))
}

func ValidateGroupInviteToken(tokenString string) (uuid.UUID, error) {
    token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(config.Env.JWT.SecretKey), nil
    })

    if err != nil || !token.Valid {
        return uuid.Nil, fmt.Errorf("invalid or expired invite token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return uuid.Nil, fmt.Errorf("invalid token claims")
    }

    tokenType, ok := claims["type"].(string)
    if !ok || tokenType != "group_invite" {
        return uuid.Nil, fmt.Errorf("invalid token type")
    }

    if exp, ok := claims["exp"].(float64); ok {
        if time.Unix(int64(exp), 0).Before(time.Now()) {
            return uuid.Nil, fmt.Errorf("invite token expired")
        }
    }

    groupIDStr, ok := claims["groupID"].(string)
    if !ok {
        return uuid.Nil, fmt.Errorf("invalid groupID in token")
    }

    groupID, err := uuid.Parse(groupIDStr)
    if err != nil {
        return uuid.Nil, fmt.Errorf("invalid UUID format")
    }

    return groupID, nil
}