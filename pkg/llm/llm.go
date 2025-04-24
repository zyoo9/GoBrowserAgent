package llm

import (
	"GoBrowserAgent/internal/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

// GenericLLM 通用LLM实现
type GenericLLM struct {
	config config.LLMConfig
}

// ChatRequest 通用聊天请求
type ChatRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

// ChatResponse 通用聊天响应
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	// 其他可能的响应字段
	Output struct {
		Text string `json:"text"`
	} `json:"output"`
	Result  string `json:"result"`
	Content string `json:"content"`
}

// NewLLM 创建语言模型实例
func NewLLM(cfg config.LLMConfig) (LLM, error) {
	return &GenericLLM{config: cfg}, nil
}

// Chat 与模型对话
func (l *GenericLLM) Chat(messages []Message) (*Response, error) {
	if l.config.APIKey == "" {
		return nil, errors.New("未设置API密钥")
	}

	// 构建请求
	reqBody := ChatRequest{
		Model:     l.config.ModelName,
		Messages:  messages,
		MaxTokens: l.config.MaxTokens,
	}

	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 确保API路径正确
	apiUrl := l.config.APIUrl
	if !strings.HasSuffix(apiUrl, "/chat/completions") && !strings.HasSuffix(apiUrl, "/") {
		apiUrl += "/chat/completions"
	} else if strings.HasSuffix(apiUrl, "/") {
		apiUrl += "chat/completions"
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// 根据不同提供商设置认证头
	if strings.HasPrefix(l.config.APIKey, "Bearer ") {
		req.Header.Set("Authorization", l.config.APIKey)
	} else {
		req.Header.Set("Authorization", "Bearer "+l.config.APIKey)
	}

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
		return nil, fmt.Errorf("API错误(%d): %s", resp.StatusCode, respBody)
	}

	// 解析响应
	var apiResp ChatResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v, 原始响应: %s", err, respBody)
	}

	// 尝试从各种可能的字段中提取响应文本
	var responseText string

	// 1. 检查标准OpenAI格式
	if len(apiResp.Choices) > 0 && apiResp.Choices[0].Message.Content != "" {
		responseText = apiResp.Choices[0].Message.Content
	} else if apiResp.Output.Text != "" {
		// 2. 检查output.text格式
		responseText = apiResp.Output.Text
	} else if apiResp.Result != "" {
		// 3. 检查result字段
		responseText = apiResp.Result
	} else if apiResp.Content != "" {
		// 4. 检查content字段
		responseText = apiResp.Content
	} else {
		// 如果无法解析出标准格式，返回原始JSON
		responseText = string(respBody)
	}

	if responseText == "" {
		return nil, errors.New("模型未返回有效响应文本")
	}

	return &Response{
		Text:      responseText,
		RawOutput: apiResp,
	}, nil
}
