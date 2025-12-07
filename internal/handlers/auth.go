package handlers

import (
	"net/http"
	"time"

	"github.com/cc673459362/myapp_server/internal/models"
	"github.com/cc673459362/myapp_server/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,alphanum,min=3,max=50"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查唯一性
		var count int64
		db.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名或邮箱已被注册"})
			return
		}

		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码处理失败"})
			return
		}

		user := models.User{
			UIN:          utils.GenerateID(),
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hash,
		}

		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "用户注册成功"})
	}
}

type LoginRequest struct {
	Identity string `json:"identity" binding:"required"` // 用户名或邮箱
	Password string `json:"password" binding:"required"`
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}

		var user models.User
		result := db.Where("username = ? OR email = ?", req.Identity, req.Identity).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效凭证"})
			return
		}

		// 检查账户锁定状态
		if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "账户已锁定",
				"until":       user.LockedUntil.Format(time.RFC3339),
				"retry_after": time.Until(*user.LockedUntil).Seconds(),
			})
			return
		}

		// 验证密码
		if err := utils.VerifyPassword(user.PasswordHash, req.Password); err != nil {
			// 记录失败尝试
			db.Model(&user).Updates(map[string]interface{}{
				"failed_login_attempt": gorm.Expr("failed_login_attempt + 1"),
			})

			// 检查是否需锁定账户（例如连续失败5次）
			if user.FailedLoginAttempt+1 >= 5 {
				lockTime := time.Now().Add(30 * time.Minute)
				db.Model(&user).Updates(map[string]interface{}{
					"locked_until": lockTime,
				})
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
			return
		}

		// 登录成功重置状态
		db.Model(&user).Updates(map[string]interface{}{
			"failed_login_attempt": 0,
			"locked_until":         nil,
		})

		// 生成JWT
		token, err := utils.GenerateJWT(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "会话创建失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      user.UIN,
			"token":   token,
			"expires": time.Now().Add(utils.JWTExpiration).Format(time.RFC3339),
		})
	}
}
