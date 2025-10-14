package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_handler "TeamTrackerBE/internal/delivery/http/handler"
	"TeamTrackerBE/internal/delivery/http/middleware"
	_repo "TeamTrackerBE/internal/domain/repository"
	_uc "TeamTrackerBE/internal/usecase"
)

func InitUserRoutes(r *gin.Engine, db *gorm.DB) {
	repo := _repo.NewUserRepository(db)
	uc := _uc.NewUserUsecase(repo)
	handler := _handler.NewUserHandler(uc)

	userRoutes := r.Group("/api/user")
	userRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("superadmin"))
	
	userRoutes.GET("/", handler.FindAll)
	userRoutes.GET("/:id", handler.FindById)

	userRoutes.POST("/", handler.Create)
	userRoutes.PATCH("/:id", handler.Update)
	
	userRoutes.DELETE("/:id", handler.Delete)
	userRoutes.DELETE("/truncate", handler.Truncate)
}
