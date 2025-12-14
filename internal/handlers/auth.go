package handlers

import (
	"net/http"
	"time"

	"errors"

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

// Register godoc
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags auth
// @Accept json
// @Produce json
// @Param registerRequest body RegisterRequest true "注册信息"
// @Success 201 {object} map[string]interface{} "注册成功"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 409 {object} map[string]interface{} "用户名或邮箱已存在"
// @Failure 500 {object} map[string]interface{} "服务器内部错误"
// @Router /auth/register [post]
func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// check if username or email already exists
		var count int64
		db.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}

		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}

		user := models.User{
			Uin:          utils.GenerateID(),
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hash,
		}

		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "success to register"})
	}
}

/**
 * @description login request payload
 */
type LoginRequest struct {
	Identity string `json:"identity" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		result := db.Where("username = ? OR email = ?", req.Identity, req.Identity).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials, please check your identity and password."})
			return
		}

		// check if account is locked
		if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Sorry, the account is locked",
				"until":       user.LockedUntil.Format(time.RFC3339),
				"retry_after": time.Until(*user.LockedUntil).Seconds(),
			})
			return
		}

		// verify password
		if err := utils.VerifyPassword(user.PasswordHash, req.Password); err != nil {
			db.Model(&user).Updates(map[string]interface{}{
				"failed_login_attempt": gorm.Expr("failed_login_attempt + 1"),
			})
			// lock account if failed attempts exceed threshold
			if user.FailedLoginAttempt+1 >= 5 {
				lockTime := time.Now().Add(30 * time.Minute)
				db.Model(&user).Updates(map[string]interface{}{
					"locked_until": lockTime,
				})
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials, please check your identity and password."})
			return
		}

		// login success reset status
		db.Model(&user).Updates(map[string]interface{}{
			"failed_login_attempt": 0,
			"locked_until":         nil,
		})

		// generate JWT
		token, err := utils.GenerateJWT(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":      user.Uin,
			"token":   token,
			"expires": time.Now().Add(utils.JWTExpiration).Format(time.RFC3339),
		})
	}
}

func GetProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		var user models.User
		result := db.Where("uin = ?", userID).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
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
