package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_handler "TeamTrackerBE/internal/delivery/http/handler"
	"TeamTrackerBE/internal/delivery/http/middleware"
	_repo "TeamTrackerBE/internal/domain/repository"
	_uc "TeamTrackerBE/internal/usecase"
)

func InitLocationRoutes(r *gin.Engine, db *gorm.DB) {
	repoUser := _repo.NewUserRepository(db)

	repo := _repo.NewLocationRepository(db)
	uc := _uc.NewLocationUsecase(repo, repoUser)
	handler := _handler.NewLocationHandler(uc)

	locationRoutes := r.Group("/api/location")
	locationRoutes.Use(middleware.AuthMiddleware())
	
	locationRoutes.GET("/", handler.FindAll)
	locationRoutes.GET("/history", handler.HistoryByUser)
	locationRoutes.GET("/:id", handler.FindById)

	locationRoutes.POST("/", handler.Create)
	locationRoutes.PATCH("/:id", handler.Update)
	locationRoutes.DELETE("/:id", handler.Delete)
	
	locationRoutes.Use(middleware.RestrictTo("superadmin"))
	locationRoutes.DELETE("/truncate", handler.Truncate)
}
