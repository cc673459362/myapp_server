package handlers

import (
	"errors"
	"net/http"

	"github.com/cc673459362/myapp_server/internal/models"
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
		profileGroup.GET("/:id", GetProfileHandler(db))
	}

	voiceRoomGroup := apiGroup.Group("/voiceroom")
	voiceRoomGroup.Use(utils.JWTMiddleware()) // JWT保护路由
	{
		voiceRoomGroup.POST("/createroom", CreateRoomHandler(db))
		//profileGroup.POST("/joinroom", JoinRoomHandler(db))
		//profileGroup.POST("/quitroom", QuitRoomHandler(db))
	}
}

func GetProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		var user models.User
		result := db.Where("id = ?", userID).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效凭证"})
			return
		}

		// 错误处理
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		})
	}
}
