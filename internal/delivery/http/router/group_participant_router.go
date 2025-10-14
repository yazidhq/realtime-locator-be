package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_handler "TeamTrackerBE/internal/delivery/http/handler"
	"TeamTrackerBE/internal/delivery/http/middleware"
	_repo "TeamTrackerBE/internal/domain/repository"
	_uc "TeamTrackerBE/internal/usecase"
)

func InitGroupParticipantRoutes(r *gin.Engine, db *gorm.DB) {
	repo := _repo.NewGroupParticipantRepository(db)
	uc := _uc.NewGroupParticipantUsecase(repo)
	handler := _handler.NewGroupParticipantHandler(uc)

	groupParticipantRoutes := r.Group("/api/group_participant")
	groupParticipantRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("admin", "member"))
	
	groupParticipantRoutes.GET("/", handler.FindAll)
	groupParticipantRoutes.GET("/:id", handler.FindById)

	groupParticipantRoutes.POST("/", handler.Create)
	groupParticipantRoutes.PATCH("/:id", handler.Update)
	groupParticipantRoutes.DELETE("/:id", handler.Delete)
	
	groupParticipantRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("superadmin"))
	groupParticipantRoutes.DELETE("/truncate", handler.Truncate)
}
