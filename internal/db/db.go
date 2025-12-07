package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() (*DBConfig, error) {
	// 本地开发时加载.env文件
	_ = godotenv.Load(".env")

	config := &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "myapp"),
	}

	if config.User == "" || config.Password == "" {
		return nil, fmt.Errorf("数据库用户名或密码未配置")
	}
	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
