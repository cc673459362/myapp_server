package db

import (
	"fmt"
	"log"
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
	if err := godotenv.Load(".env"); err == nil {
		log.Println("✅ 从当前目录加载 .env")
	} else {
		// 2. 尝试上级目录（项目根目录）
		if err := godotenv.Load("../.env"); err == nil {
			log.Println("✅ 从上级目录加载 .env")
		} else {
			// 3. 尝试上上级目录
			if err := godotenv.Load("../../.env"); err == nil {
				log.Println("✅ 从上上级目录加载 .env")
			} else {
				log.Println("⚠️  未找到 .env 文件，使用环境变量")
			}
		}
	}

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
