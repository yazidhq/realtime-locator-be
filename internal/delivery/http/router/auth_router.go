package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_handler "TeamTrackerBE/internal/delivery/http/handler"
	_repo "TeamTrackerBE/internal/domain/repository"
	_uc "TeamTrackerBE/internal/usecase"
)

func InitAuthRoutes(r *gin.Engine, db *gorm.DB) {
	repo := _repo.NewUserRepository(db)
	uc := _uc.NewAuthUsecase(repo)
	handler := _handler.NewAuthHandler(uc)

	authRoutes := r.Group("/api/auth")
	authRoutes.POST("/register", handler.Register)
	authRoutes.POST("/login", handler.Login)
	authRoutes.POST("/refresh_token", handler.RefreshToken)
}
