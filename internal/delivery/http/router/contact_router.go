package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_handler "TeamTrackerBE/internal/delivery/http/handler"
	"TeamTrackerBE/internal/delivery/http/middleware"
	_repo "TeamTrackerBE/internal/domain/repository"
	_uc "TeamTrackerBE/internal/usecase"
)

func InitContactRoutes(r *gin.Engine, db *gorm.DB) {
	repo := _repo.NewContactRepository(db)
	uc := _uc.NewContactUsecase(repo)
	handler := _handler.NewContactHandler(uc)

	contactRoutes := r.Group("/api/contact")
	contactRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("admin", "member"))
	
	contactRoutes.GET("/", handler.FindAll)
	contactRoutes.GET("/:id", handler.FindById)

	contactRoutes.POST("/", handler.Create)
	contactRoutes.PATCH("/:id", handler.Update)
	contactRoutes.DELETE("/:id", handler.Delete)
	
	contactRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("superadmin"))
	contactRoutes.DELETE("/truncate", handler.Truncate)
}
