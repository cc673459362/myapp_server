package models

import (
	"time"
)

type User struct {
	ID                 uint   `gorm:"primaryKey"`
	Username           string `gorm:"uniqueIndex;size:50;not null"`
	Email              string `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash       string `gorm:"type:char(97);not null"`
	FailedLoginAttempt uint8  `gorm:"default:0"`
	LockedUntil        *time.Time
	PasswordChangedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
}
