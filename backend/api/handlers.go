package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvsoso/tools/backend/auth"
	"github.com/lvsoso/tools/backend/db"
	"github.com/lvsoso/tools/backend/services"
	"golang.org/x/crypto/bcrypt"
)

// 用户认证相关处理函数
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 查找用户
	var user db.User
	result := db.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 检查邮箱是否已存在
	var existingUser db.User
	if result := db.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新用户
	user := db.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if result := db.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user db.User
	if result := db.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func Logout(c *gin.Context) {
	// 实际的注销逻辑可以在这里实现
	// 例如：将令牌加入黑名单等
	c.JSON(http.StatusOK, gin.H{"message": "注销成功"})
}

// 聊天相关处理函数
type SendMessageRequest struct {
	Model   string `json:"model" binding:"required"`
	Message string `json:"message" binding:"required"`
}

// 限制发送给模型的消息数量
const MAX_CONTEXT_MESSAGES = 10

func SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 创建或获取当前会话
	var conversation db.Conversation
	result := db.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&conversation)
	if result.Error != nil {
		// 如果没有会话，创建新会话
		conversation = db.Conversation{
			UserID: userID.(uint),
			Title:  "新对话",
		}
		db.DB.Create(&conversation)
	}

	// 保存用户消息
	userMessage := db.Message{
		Role:           "user",
		Content:        req.Message,
		ModelName:      req.Model,
		TokenCount:     0, // 用户消息的token计数初始为0
		ConversationID: conversation.ID,
	}
	db.DB.Create(&userMessage)

	// 获取当前会话的历史消息
	var messages []db.Message
	db.DB.Where("conversation_id = ?", conversation.ID).Order("created_at ASC").Find(&messages)

	// 转换为AI服务需要的格式
	var chatMessages []services.ChatMessage
	for _, msg := range messages {
		chatMessages = append(chatMessages, services.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 限制上下文窗口大小
	if len(chatMessages) > MAX_CONTEXT_MESSAGES {
		// 只保留最近的N条消息
		chatMessages = chatMessages[len(chatMessages)-MAX_CONTEXT_MESSAGES:]
	}

	// 调用AI服务获取响应
	response, err := services.GetAIResponse(req.Model, chatMessages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取AI响应失败"})
		return
	}

	// 保存AI响应
	aiMessage := db.Message{
		Role:           "assistant",
		Content:        response.Content,
		ModelName:      req.Model,
		TokenCount:     response.TokenCount,
		ConversationID: conversation.ID,
	}
	db.DB.Create(&aiMessage)

	c.JSON(http.StatusOK, gin.H{
		"content":     response.Content,
		"token_count": response.TokenCount,
	})
}

func GetConversations(c *gin.Context) {
	userID, _ := c.Get("userID")

	var conversations []db.Conversation
	db.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&conversations)

	c.JSON(http.StatusOK, conversations)
}

func DeleteConversation(c *gin.Context) {
	userID, _ := c.Get("userID")
	conversationID := c.Param("id")

	// 验证会话所有权
	var conversation db.Conversation
	result := db.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	// 删除会话及其消息
	db.DB.Delete(&conversation)

	c.JSON(http.StatusOK, gin.H{"message": "会话已删除"})
}

func GetConversationMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	conversationID := c.Param("id")

	// 验证会话所有权
	var conversation db.Conversation
	result := db.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "会话不存在"})
		return
	}

	// 获取会话消息
	var messages []db.Message
	db.DB.Where("conversation_id = ?", conversation.ID).Order("created_at ASC").Find(&messages)

	c.JSON(http.StatusOK, messages)
}
