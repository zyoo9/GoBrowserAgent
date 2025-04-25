package llm

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// Config 存储LLM配置信息
type Config struct {
	APIEndpoint string  `json:"api_endpoint"`
	Model       string  `json:"model"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	APIKey      string  `json:"api_key"`
}

// LoadConfig 从配置文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		logrus.Errorf("读取配置文件失败: %v", err)
		return nil, err
	}

	// 解析JSON配置
	var configData struct {
		LLM Config `json:"llm"`
	}

	if err := json.Unmarshal(data, &configData); err != nil {
		logrus.Errorf("解析配置文件失败: %v", err)
		return nil, err
	}

	config := configData.LLM

	// 检查环境变量中是否有API密钥
	if apiKey := os.Getenv("LLM_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	return &config, nil
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	// 默认配置路径
	return filepath.Join("internal", "web", "config.json")
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		APIEndpoint: "https://api.openai.com/v1/chat/completions",
		Model:       "gpt-3.5-turbo",
		MaxTokens:   2000,
		Temperature: 0.7,
		TopP:        1.0,
		APIKey:      os.Getenv("LLM_API_KEY"), // 尝试从环境变量获取
	}
}
