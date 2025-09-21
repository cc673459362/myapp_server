package handlers

import (
	"net/http"

	"github.com/cc673459362/myapp_server/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateRoomRequest struct {
	Roomname string `json:"room_name" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
}

type JoinRoomRequest struct {
	RoomID string `json:"room_id" binding:"required"`
	UserID uint   `json:"user_id" binding:"required"`
}

type RoomResponse struct {
	RoomID     string `json:"room_id"`
	RoomName   string `json:"room_name"`
	Creator    uint   `json:"creator,omitempty"`
	ServerIp   string `json:"server_ip,omitempty"`
	ServerPort string `json:"server_port,omitempty"`
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

func JoinRoomHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req JoinRoomRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// 2. 查找房间
		roomUUID, err := uuid.Parse(req.RoomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		var room models.Room
		uuidBytes, err := roomUUID.MarshalBinary()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		if err := db.Where("uuid = ?", uuidBytes).First(&room).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		// 3. 模拟用户加入房间逻辑（如更新数据库等）
		// 这里什么都不用干，因为房间转发逻辑在nginx中处理
		// 4. 返回响应
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": RoomResponse{
				RoomID:     req.RoomID,
				RoomName:   room.Name,
				ServerIp:   "42.194.195.181",
				ServerPort: "20101",
			},
		})
	}
}
