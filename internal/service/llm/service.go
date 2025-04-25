package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ChatMessage 定义聊天消息结构
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 定义聊天请求结构
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	TopP        float64       `json:"top_p,omitempty"`
}

// ChatResponse 定义LLM API响应结构
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

// Service LLM服务
type Service struct {
	Config *Config
}

// NewService 创建新的LLM服务
func NewService(config *Config) *Service {
	return &Service{
		Config: config,
	}
}

// Chat 处理与LLM的聊天
func (s *Service) Chat(userMessage string) (string, error) {
	if s.Config.APIKey == "" {
		return "", fmt.Errorf("未配置API密钥，请在配置文件中设置api_key或通过环境变量LLM_API_KEY设置")
	}

	// 构建请求
	chatReq := ChatRequest{
		Model: s.Config.Model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		MaxTokens:   s.Config.MaxTokens,
		Temperature: s.Config.Temperature,
		TopP:        s.Config.TopP,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("无法序列化聊天请求: %v", err)
	}

	// 发送请求到LLM API
	req, err := http.NewRequest("POST", s.Config.APIEndpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Config.APIKey)

	logrus.Debugf("发送请求到LLM API: %s, 模型: %s", s.Config.APIEndpoint, s.Config.Model)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求到LLM API失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API返回错误: %s", string(respBody))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("解析LLM响应失败: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("LLM未返回任何响应")
	}

	// 返回响应
	return chatResp.Choices[0].Message.Content, nil
}
