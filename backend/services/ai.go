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
}

// GetAIResponse 根据选择的模型调用相应的AI服务
func GetAIResponse(model, message string) (string, error) {
	// 如果模型为空，默认使用deepseek
	if model == "" {
		model = "deepseek"
	}

	switch model {
	case "gpt-4", "gpt-3.5-turbo":
		// 检查OpenAI API密钥是否配置
		if config.AppConfig.OpenAIAPIKey == "" {
			return "", fmt.Errorf("未配置OpenAI API密钥")
		}
		return callOpenAI(model, message)
	case "deepseek":
		// 检查Deepseek API密钥是否配置
		if config.AppConfig.DeepseekAPIKey == "" {
			return "", fmt.Errorf("未配置Deepseek API密钥")
		}
		return callDeepseek(message)
	default:
		return "", fmt.Errorf("不支持的模型: %s", model)
	}
}

// callOpenAI 调用OpenAI API
func callOpenAI(model, message string) (string, error) {
	reqBody := OpenAIRequest{
		Model: model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AppConfig.OpenAIAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API返回错误状态码: %d", resp.StatusCode)
	}

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("OpenAI API返回空响应")
	}

	return response.Choices[0].Message.Content, nil
}

// callDeepseek 调用Deepseek API
func callDeepseek(message string) (string, error) {
	// Deepseek API请求结构与OpenAI类似
	reqBody := OpenAIRequest{
		Model: "deepseek-chat", // 使用Deepseek的模型名称
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: message,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Deepseek API端点
	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AppConfig.DeepseekAPIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Deepseek API返回错误状态码: %d", resp.StatusCode)
	}

	// 假设Deepseek的响应格式与OpenAI相同
	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("Deepseek API返回空响应")
	}

	return response.Choices[0].Message.Content, nil
}
