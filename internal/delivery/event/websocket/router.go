package websocket

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitWSRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/api/live_track", LiveTrack)
}
