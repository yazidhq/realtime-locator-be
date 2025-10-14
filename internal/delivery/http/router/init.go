package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	InitAuthRoutes(r, db)
	InitUserRoutes(r, db)
	InitContactRoutes(r, db)
	InitGroupRoutes(r, db)
	InitGroupParticipantRoutes(r, db)
}
