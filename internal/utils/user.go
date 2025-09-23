package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) uint {
	fmt.Print("c: ", c)
	if userIDRaw, exists := c.Get("userID"); exists {
		fmt.Print("userIDRaw: ", userIDRaw)
		if userID, ok := userIDRaw.(uint); ok {
			fmt.Print("userID: ", userID)
			return userID
		}
	}
	fmt.Print("return 0")
	return 0
}
