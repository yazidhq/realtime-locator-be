package websocket

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"TeamTrackerBE/internal/delivery/http/middleware"
	_repo "TeamTrackerBE/internal/domain/repository"
)

func InitWSRoutes(r *gin.Engine, db *gorm.DB) {
	locationRepo := _repo.NewLocationRepository(db)
	LocationRepository(locationRepo)

	userRepo := _repo.NewUserRepository(db)
	liveTrackHandler := NewLiveTrackHandler(userRepo)
	
	r.Use(middleware.AuthMiddleware())
	r.GET("/api/live_track", liveTrackHandler.LiveTrack)
}
