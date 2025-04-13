package handlers

import (
	"github.com/cc673459362/myapp_server/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	apiGroup := router.Group("/myapp_server/api")
	authGroup := apiGroup.Group("/auth")
	{
		authGroup.POST("/register", RegisterHandler(db))
		authGroup.POST("/login", LoginHandler(db))
	}

	profileGroup := apiGroup.Group("/profile")
	profileGroup.Use(utils.JWTMiddleware()) // JWT保护路由
	{
		profileGroup.GET("", GetProfileHandler)
	}
}

func GetProfileHandler(c *gin.Context) {
	userID, _ := c.Get("userID")
	c.JSON(200, gin.H{"user_id": userID})
}
