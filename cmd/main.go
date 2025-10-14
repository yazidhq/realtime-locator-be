package main

import (
	"TeamTrackerBE/internal/config"
	"TeamTrackerBE/internal/delivery/event/websocket"
	"TeamTrackerBE/internal/delivery/http/middleware"
	"TeamTrackerBE/internal/delivery/http/router"
	"TeamTrackerBE/internal/domain/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil { 
		log.Println("No .env file found, using environment variables") 
	}

	config.LoadEnv()

	db := config.SetupDatabase()
	db.AutoMigrate(model.Models...)

	r := gin.Default()
	r.MaxMultipartMemory = 15 << 30

	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.SecurityHeaders())
	
	r.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, "/") })
	router.InitRoutes(r, db)
	websocket.InitWSRoutes(r, db)

	port := config.Env.App.Port
	r.Run(":" + port)
}