package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_handler "TeamTrackerBE/internal/delivery/http/handler"
	"TeamTrackerBE/internal/delivery/http/middleware"
	_repo "TeamTrackerBE/internal/domain/repository"
	_uc "TeamTrackerBE/internal/usecase"
)

func InitGroupRoutes(r *gin.Engine, db *gorm.DB) {
	repo := _repo.NewGroupRepository(db)
	uc := _uc.NewGroupUsecase(repo)
	handler := _handler.NewGroupHandler(uc)

	groupRoutes := r.Group("/api/group")
	groupRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("admin", "customer"))
	
	groupRoutes.GET("/", handler.FindAll)
	groupRoutes.GET("/:id", handler.FindById)

	groupRoutes.POST("/", handler.Create)
	groupRoutes.PATCH("/:id", handler.Update)
	
	groupRoutes.Use(middleware.AuthMiddleware(), middleware.RestrictTo("superadmin"))
	groupRoutes.DELETE("/:id", handler.Delete)
	groupRoutes.DELETE("/truncate", handler.Truncate)
}
