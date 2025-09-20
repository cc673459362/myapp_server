package handlers

import (
	"net/http"

	"github.com/cc673459362/myapp_server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateRoomRequest struct {
	Roomname string `json:"roomname" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
}

type RoomResponse struct {
	RoomID   string `json:"room_id"`
	RoomName string `json:"room_name"`
	Creator  uint   `json:"creator"`
}

func CreateRoomHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateRoomRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. 创建房间号UUID
		roomID, _ := uuid.New().MarshalBinary()

		// 3. 创建房间（模拟数据库操作）
		room := models.Room{
			Name:    req.Roomname,
			OwnerID: req.UserID,
			UUID:    roomID,
		}
		if err := db.Create(&room).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to create room: " + err.Error(),
			})
			return
		}

		// 4. 返回响应
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"data": RoomResponse{
				RoomID:   uuid.Must(uuid.FromBytes(room.UUID)).String(),
				RoomName: room.Name,
				Creator:  req.UserID,
			},
		})
	}
}
