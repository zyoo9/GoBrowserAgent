package main

import (
	"GoBrowserAgent/internal/service/llm"
	"GoBrowserAgent/internal/web"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("开始启动web服务")

	// 加载LLM配置
	configPath := llm.GetConfigPath()
	llmConfig, err := llm.LoadConfig(configPath)
	if err != nil {
		logrus.Warnf("加载LLM配置失败: %v, 将使用默认配置", err)
		// 使用默认配置
		llmConfig = llm.GetDefaultConfig()
	}

	// 创建LLM服务
	llmService := llm.NewService(llmConfig)

	// 创建API处理程序
	apiHandler := web.NewAPIHandler(llmService)
	apiHandler.RegisterHandlers()

	// 设置静态文件服务
	fs := http.FileServer(http.Dir(filepath.Join("internal", "web", "static")))
	http.Handle("/", fs)

	// 启动HTTP服务器
	port := "8080"
	logrus.Infof("Web服务器正在运行: http://localhost:%s", port)
	logrus.Infof("配置的LLM模型: %s", llmConfig.Model)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Errorf("HTTP服务器启动失败: %v", err)
		os.Exit(1)
	}
}
