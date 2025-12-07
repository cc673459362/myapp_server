package handlers

import (
	"github.com/cc673459362/myapp_server/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/**
 * @description Gin router setup
 */
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	apiGroup := router.Group("/myapp_server/api")
	authGroup := apiGroup.Group("/auth")
	{
		authGroup.POST("/register", RegisterHandler(db))
		authGroup.POST("/login", LoginHandler(db))
	}

	profileGroup := apiGroup.Group("/profile")
	profileGroup.Use(utils.JWTMiddleware())
	{
		profileGroup.GET("/:id", GetProfileHandler(db))
	}

	voiceRoomGroup := apiGroup.Group("/voiceroom")
	voiceRoomGroup.Use(utils.JWTMiddleware())
	{
		voiceRoomGroup.POST("/createroom", CreateRoomHandler(db))
		voiceRoomGroup.POST("/joinroom", JoinRoomHandler(db))
	}
}
