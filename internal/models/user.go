package models

import (
	"time"
)

type User struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement"`
	Uin                uint64 `gorm:"uniqueIndex;not null;comment:全局唯一标识号"`
	Username           string `gorm:"uniqueIndex;size:50;not null"`
	Email              string `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash       string `gorm:"type:char(97);not null"`
	FailedLoginAttempt uint8  `gorm:"default:0"`
	LockedUntil        *time.Time
	PasswordChangedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}

type Room struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`             // 内部自增ID（可选）
	UUID    []byte `gorm:"type:binary(16);uniqueIndex;not null"` // 二进制UUID + 唯一索引
	Name    string `gorm:"size:100;not null"`
	OwnerID uint   `gorm:"not null"`
}
