package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
    "github.com/gin-gonic/gin"
)

const (
    JWTExpiration = 12 * time.Hour
    SecretKey     = "your-256-bit-secret" // 替换为真实的随机密钥
)

func GenerateJWT(userID uint, username string) (string, error) {
    claims := jwt.MapClaims{
        "sub":  userID,
        "name": username,
        "exp":  time.Now().Add(JWTExpiration).Unix(),
        "iat":  time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(SecretKey))
}

func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if len(tokenString) < 8 { // "Bearer " 长度校验
            c.AbortWithStatusJSON(401, gin.H{"error": "未授权的访问"})
            return
        }

        token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
            return []byte(SecretKey), nil
        })

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Set("userID", claims["sub"])
            c.Set("username", claims["name"])
            c.Next()
        } else {
            c.AbortWithStatusJSON(401, gin.H{"error": "无效令牌", "details": err.Error()})
        }
    }
}
