package main

import (
	"fmt"
	"log"

	"github.com/cc673459362/myapp_server/internal/db"
	"github.com/cc673459362/myapp_server/internal/handlers"
	"github.com/cc673459362/myapp_server/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 配置 MySQL 连接（根据实际情况调整）
	dbConfig, err := db.LoadConfig()
	if err != nil {
		log.Fatal("数据库配置加载失败: ", err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移表结构（生产环境应使用迁移工具）
	db.AutoMigrate(&models.User{}, &models.Room{})

	// 初始化Gin
	router := gin.Default()

	// 注册路由
	handlers.SetupRoutes(router, db)

	// 启动服务
	router.Run(":8080")
}
