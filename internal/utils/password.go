package utils

import "golang.org/x/crypto/bcrypt"

const (
	bcryptCost = 12 // 安全建议值，每+1计算时间翻倍
)

// 生成密码哈希
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcryptCost,
	)
	return string(hashed), err
}

// 验证密码 (防时序攻击版)
func VerifyPassword(hashedPassword, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(inputPassword),
	)
	return err
}
