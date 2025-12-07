package handlers

import (
	"encoding/binary"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/cc673459362/myapp_server/internal/models"
	"github.com/cc673459362/myapp_server/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateRoomRequest struct {
	Roomname string `json:"room_name" binding:"required"`
}

type JoinRoomRequest struct {
	RoomID string `json:"room_id" binding:"required"`
}

type RoomResponse struct {
	RoomID     string `json:"room_id"`
	RoomName   string `json:"room_name"`
	Creator    uint   `json:"creator,omitempty"`
	ServerIp   string `json:"server_ip,omitempty"`
	ServerPort string `json:"server_port,omitempty"`
}

var roomIDCounter uint32 = 0 // 初始值

func GenerateRoomID() uint32 {
	return atomic.AddUint32(&roomIDCounter, 1) // 线程安全递增
}

func Uint32ToBinary16(value uint32) []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint32(buf[12:], value) // 高地址存有效数据，低地址补零
	return buf
}

func Uint64ToBinary16(value uint64) []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[8:], value)
	return buf
}

func CreateRoomHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateRoomRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 1. 验证用户身份
		userId := utils.GetUserID(c)
		if userId == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			return
		}

		// 2. 创建房间号UUID
		roomUUID := utils.GenerateID()

		// 3. 创建房间（模拟数据库操作）
		room := models.Room{
			Name:    req.Roomname,
			OwnerID: userId,
			UUID:    Uint64ToBinary16(roomUUID),
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
			"room_id":   strconv.FormatUint(uint64(roomUUID), 10), // 转为字符串
			"room_name": room.Name,
			"creator":   userId,
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

		// 1. 验证用户身份
		userId := utils.GetUserID(c)
		if userId == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
			return
		}

		// 2. 查找房间
		roomID, err := strconv.ParseUint(req.RoomID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		uuidBytes := Uint64ToBinary16((uint64(roomID)))
		var room models.Room
		if err := db.Where("uuid = ?", uuidBytes).First(&room).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}
		// 3. 这里可以添加用户加入房间的逻辑，比如更新数据库等，也可以直接在nginx做

		// 4. 返回响应
		c.JSON(http.StatusOK, gin.H{
			"room_id":     req.RoomID,
			"room_name":   room.Name,
			"server_ip":   "42.194.195.181",
			"server_port": "10086",
		})
	}
}
