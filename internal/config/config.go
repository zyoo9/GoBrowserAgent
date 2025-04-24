package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 存储全局配置信息
type Config struct {
	Browser BrowserConfig `json:"browser"`
	LLM     LLMConfig     `json:"llm"`
	Log     LogConfig     `json:"log"`
	Web     WebConfig     `json:"web"`
}

// BrowserConfig 浏览器配置
type BrowserConfig struct {
	Headless      bool   `json:"headless"`       // 是否隐藏浏览器界面
	UserDataDir   string `json:"user_data_dir"`  // 浏览器用户数据目录
	DefaultWidth  int    `json:"default_width"`  // 浏览器默认宽度
	DefaultHeight int    `json:"default_height"` // 浏览器默认高度
	Timeout       int    `json:"timeout"`        // 操作超时时间(秒)
}

// LLMConfig 大语言模型配置
type LLMConfig struct {
	Provider  string `json:"provider"`   // 提供商名称
	APIKey    string `json:"api_key"`    // API密钥
	APIUrl    string `json:"api_url"`    // API地址
	ModelName string `json:"model_name"` // 模型名称
	MaxTokens int    `json:"max_tokens"` // 最大token数
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `json:"level"`       // 日志级别
	OutputPath string `json:"output_path"` // 日志输出路径
}

// WebConfig Web服务器配置
type WebConfig struct {
	Port int `json:"port"` // 服务端口号
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Browser: BrowserConfig{
			Headless:      false,
			UserDataDir:   "./user_data",
			DefaultWidth:  1280,
			DefaultHeight: 800,
			Timeout:       30,
		},
		LLM: LLMConfig{
			Provider:  "openai",
			APIUrl:    "https://api.openai.com/v1",
			ModelName: "gpt-4",
			MaxTokens: 4096,
		},
		Log: LogConfig{
			Level:      "info",
			OutputPath: "./logs",
		},
		Web: WebConfig{
			Port: 8080,
		},
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	// 确保文件存在
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件路径 %s: %v", path, err)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("配置文件 %s 不存在", absPath)
		}
		return nil, fmt.Errorf("无法访问配置文件 %s: %v", absPath, err)
	}

	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s 是一个目录，不是配置文件", absPath)
	}

	// 读取文件内容
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件 %s 失败: %v", absPath, err)
	}

	// 解析JSON
	config := DefaultConfig()
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件 %s 失败: %v", absPath, err)
	}

	return config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, path string) error {
	// 创建目录（如果不存在）
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("创建配置目录 %s 失败: %v", dir, err)
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("写入配置文件 %s 失败: %v", path, err)
	}

	return nil
}
