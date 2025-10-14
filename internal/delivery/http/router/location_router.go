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
	repo := _repo.NewLocationRepository(db)
	uc := _uc.NewLocationUsecase(repo)
	handler := _handler.NewLocationHandler(uc)

	locationRoutes := r.Group("/api/location")
	locationRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("admin", "member"))
	
	locationRoutes.GET("/", handler.FindAll)
	locationRoutes.GET("/:id", handler.FindById)

	locationRoutes.POST("/", handler.Create)
	locationRoutes.PATCH("/:id", handler.Update)
	locationRoutes.DELETE("/:id", handler.Delete)
	
	locationRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("superadmin"))
	locationRoutes.DELETE("/truncate", handler.Truncate)
}
