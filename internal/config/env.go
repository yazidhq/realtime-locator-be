package config

import (
	"os"
)

type EnvConfig struct {
	App       App
	Database  Database
	JWT       JWT
}

type App struct {
	Port string `json:"port"`
}

type Database struct {
	Host 	 string `json:"host"`
	Port 	 string `json:"port"`
	User 	 string `json:"user"`
	Password string `json:"password"`
	Name 	 string `json:"name"`
}

type JWT struct {
	SecretKey 			  string `json:"secret_key"`
	ExpiresIn 			  string `json:"expires_in"`
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
}

var Env *EnvConfig
func LoadEnv() *EnvConfig {
	Env = &EnvConfig{
		App: App{
			Port: getEnv("APP_PORT", "3000"),
		},

		Database: Database{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "secret"),
			Name: getEnv("DB_NAME", "mydb"),
		},

		JWT: JWT{
			SecretKey: getEnv("JWT_SECRET_KEY", "supersecret"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
			RefreshTokenExpiresIn: getEnv("JWT_REFRESH_TOKEN_EXPIRES_IN", "7d"),
		},
	}

	return Env
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	
	return defaultValue
}
