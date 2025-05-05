package db

import (
	"fmt"

	"github.com/lvsoso/tools/backend/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

// 初始化PostgreSQL数据库连接
func InitDB() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&User{},
		&Conversation{},
		&Message{},
	); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	DB = db
	return nil
}

// 初始化Redis连接
func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort),
		Password: config.AppConfig.RedisPassword,
		DB:       config.AppConfig.RedisDB,
	})

	return nil
}

// User 用户模型
type User struct {
	gorm.Model
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
	// 关联
	Conversations []Conversation
}

// Conversation 对话模型
type Conversation struct {
	gorm.Model
	Title    string
	UserID   uint
	User     User
	Messages []Message
}

// Message 消息模型
type Message struct {
	gorm.Model
	Role           string // 'user' 或 'assistant'
	Content        string
	ModelName      string // 使用的AI模型
	ConversationID uint
	Conversation   Conversation
}
