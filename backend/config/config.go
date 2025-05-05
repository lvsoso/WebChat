package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// 服务器配置
	ServerPort string

	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Redis配置
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// JWT配置
	JWTSecret         string
	JWTExpireDuration time.Duration

	// AI模型配置
	OpenAIAPIKey   string
	DeepseekAPIKey string
}

var AppConfig Config

func LoadConfig() {
	godotenv.Load()

	// 从环境变量加载配置
	AppConfig = Config{
		// 服务器配置
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// 数据库配置
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "webchat"),

		// Redis配置
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,

		// JWT配置
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpireDuration: 24 * time.Hour,

		// AI模型配置
		OpenAIAPIKey:   getEnv("OPENAI_API_KEY", ""),
		DeepseekAPIKey: getEnv("DEEPSEEK_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
