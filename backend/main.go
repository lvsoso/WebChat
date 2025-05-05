package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lvsoso/tools/backend/api"
	"github.com/lvsoso/tools/backend/auth"
	"github.com/lvsoso/tools/backend/config"
	"github.com/lvsoso/tools/backend/db"
)

func main() {
	// 加载配置
	config.LoadConfig()

	// 初始化数据库连接
	if err := db.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化Redis连接
	if err := db.InitRedis(); err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}

	// 创建Gin路由
	r := gin.Default()

	// 注册中间件
	r.Use(auth.CORSMiddleware())

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/auth/login", api.Login)
		public.POST("/auth/register", api.Register)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(auth.JWTMiddleware())
	{
		// 用户相关
		protected.GET("/auth/me", api.GetCurrentUser)
		protected.POST("/auth/logout", api.Logout)

		// 聊天相关
		protected.POST("/chat/send", api.SendMessage)
		protected.GET("/chat/conversations", api.GetConversations)
		protected.DELETE("/chat/conversations/:id", api.DeleteConversation)
		protected.GET("/chat/conversations/:id", api.GetConversation)
		protected.GET("/chat/conversations/:id/messages", api.GetConversationMessages)
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
