package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lvsoso/tools/backend/config"
)

// OpenAI API请求结构
type OpenAIRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// AIResponse 包含AI响应内容和token使用情况
type AIResponse struct {
	Content    string
	TokenCount int
}

// GetAIResponse 根据选择的模型调用相应的AI服务
func GetAIResponse(model string, messages []ChatMessage) (AIResponse, error) {
	// 如果模型为空，默认使用deepseek
	if model == "" {
		model = "deepseek"
	}

	switch model {
	case "gpt-4", "gpt-3.5-turbo":
		// 检查OpenAI API密钥是否配置
		if config.AppConfig.OpenAIAPIKey == "" {
			return AIResponse{}, fmt.Errorf("未配置OpenAI API密钥")
		}
		return callOpenAI(model, messages)
	case "deepseek":
		// 检查Deepseek API密钥是否配置
		if config.AppConfig.DeepseekAPIKey == "" {
			return AIResponse{}, fmt.Errorf("未配置Deepseek API密钥")
		}
		return callDeepseek(messages)
	default:
		return AIResponse{}, fmt.Errorf("不支持的模型: %s", model)
	}
}

// callOpenAI 调用OpenAI API
func callOpenAI(model string, messages []ChatMessage) (AIResponse, error) {
	reqBody := OpenAIRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return AIResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return AIResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AppConfig.OpenAIAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AIResponse{}, fmt.Errorf("OpenAI API返回错误状态码: %d", resp.StatusCode)
	}

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return AIResponse{}, err
	}

	if len(response.Choices) == 0 {
		return AIResponse{}, fmt.Errorf("OpenAI API返回空响应")
	}

	return AIResponse{
		Content:    response.Choices[0].Message.Content,
		TokenCount: response.Usage.TotalTokens,
	}, nil
}

// callDeepseek 调用Deepseek API
func callDeepseek(messages []ChatMessage) (AIResponse, error) {
	// Deepseek API请求结构与OpenAI类似
	reqBody := OpenAIRequest{
		Model:    "deepseek-chat", // 使用Deepseek的模型名称
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return AIResponse{}, err
	}

	// Deepseek API端点
	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return AIResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AppConfig.DeepseekAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AIResponse{}, fmt.Errorf("Deepseek API返回错误状态码: %d", resp.StatusCode)
	}

	// 假设Deepseek的响应格式与OpenAI相同
	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return AIResponse{}, err
	}

	if len(response.Choices) == 0 {
		return AIResponse{}, fmt.Errorf("Deepseek API返回空响应")
	}

	return AIResponse{
		Content:    response.Choices[0].Message.Content,
		TokenCount: response.Usage.TotalTokens,
	}, nil
}
