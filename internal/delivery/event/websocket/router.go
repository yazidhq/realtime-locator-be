package websocket

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"TeamTrackerBE/internal/delivery/http/middleware"
	_repo "TeamTrackerBE/internal/domain/repository"
)

func InitWSRoutes(r *gin.Engine, db *gorm.DB) {
	groupRepo := _repo.NewGroupRepository(db)
	userRepo := _repo.NewUserRepository(db)
	userGroupParticipant := _repo.NewGroupParticipantRepository(db)
	liveTrackHandler := NewLiveTrackHandler(groupRepo, userRepo, userGroupParticipant)

	r.Use(middleware.AuthMiddleware())
	r.GET("/api/live_track", liveTrackHandler.LiveTrack)
}
