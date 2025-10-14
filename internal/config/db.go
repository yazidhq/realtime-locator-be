package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() *gorm.DB {
	env := Env.Database

	var (
		host     = env.Host
		port     = env.Port
		user     = env.User
		password = env.Password
		dbName   = env.Name
	)

	if host == "" || port == "" || user == "" || password == "" || dbName == "" {
		panic("Database environment variables are not set properly")
	}

	sqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	var db *gorm.DB
	var err error

	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})
		if err == nil {
			log.Println("Connected to database successfully")
			return db
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i, maxRetries, err)
		time.Sleep(5 * time.Second)
	}

	panic(fmt.Sprintf("Failed to connect database after %d attempts: %v", maxRetries, err))
}
