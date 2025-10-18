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
	repoUser := _repo.NewUserRepository(db)

	repo := _repo.NewGroupRepository(db)
	uc := _uc.NewGroupUsecase(repo, repoUser)
	handler := _handler.NewGroupHandler(uc)

	groupRoutes := r.Group("/api/group")
	groupRoutes.Use(middleware.AuthMiddleware())
	
	groupRoutes.GET("/", handler.FindAll)
	groupRoutes.GET("/:id", handler.FindById)

	groupRoutes.POST("/", handler.Create)
	groupRoutes.PATCH("/:id", handler.Update)
	groupRoutes.DELETE("/:id", handler.Delete)
	groupRoutes.POST("/:id/invite", handler.Invite)

	groupRoutes.Use(middleware.RestrictTo("superadmin"))
	groupRoutes.DELETE("/truncate", handler.Truncate)
}
