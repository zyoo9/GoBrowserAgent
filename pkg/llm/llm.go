package llm

import (
	"GoBrowserAgent/internal/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// LLM 大型语言模型接口
type LLM interface {
	// Chat 与模型进行对话
	Chat(messages []Message) (*Response, error)
}

// Message 表示对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response 模型响应
type Response struct {
	Text      string `json:"text"`
	RawOutput any    `json:"raw_output"`
}

// OpenAILLM OpenAI实现
type OpenAILLM struct {
	config config.LLMConfig
}

// OpenAIChatRequest OpenAI聊天请求
type OpenAIChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

// OpenAIChatResponse OpenAI聊天响应
type OpenAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewLLM 创建语言模型实例
func NewLLM(cfg config.LLMConfig) (LLM, error) {
	switch cfg.Provider {
	case "openai":
		return &OpenAILLM{config: cfg}, nil
	default:
		return nil, fmt.Errorf("不支持的LLM提供商: %s", cfg.Provider)
	}
}

// Chat 与OpenAI模型对话
func (l *OpenAILLM) Chat(messages []Message) (*Response, error) {
	if l.config.APIKey == "" {
		return nil, errors.New("未设置OpenAI API密钥")
	}

	// 构建请求
	reqBody := OpenAIChatRequest{
		Model:     l.config.ModelName,
		Messages:  messages,
		MaxTokens: l.config.MaxTokens,
	}

	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", l.config.APIUrl+"/chat/completions", bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.config.APIKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API错误: %s", respBody)
	}

	// 解析响应
	var apiResp OpenAIChatResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, err
	}

	// 提取结果
	if len(apiResp.Choices) == 0 {
		return nil, errors.New("模型未返回响应")
	}

	return &Response{
		Text:      apiResp.Choices[0].Message.Content,
		RawOutput: apiResp,
	}, nil
}
