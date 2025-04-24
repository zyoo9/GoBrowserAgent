package web

import (
	"GoBrowserAgent/internal/config"
	"GoBrowserAgent/internal/log"
	"GoBrowserAgent/pkg/llm"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
)

// Server 表示Web服务器
type Server struct {
	config   *config.Config
	logger   log.Logger
	llmModel llm.LLM
}

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Message string `json:"message"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// NewServer 创建一个新的Web服务器
func NewServer(cfg *config.Config, logger log.Logger) (*Server, error) {
	logger.Info("正在初始化LLM，提供商: %s, API地址: %s", cfg.LLM.Provider, cfg.LLM.APIUrl)

	if cfg.LLM.APIKey == "" {
		logger.Warn("警告: LLM API密钥未设置，可能无法正常工作")
	}

	llmModel, err := llm.NewLLM(cfg.LLM)
	if err != nil {
		logger.Error("初始化LLM失败: %v", err)
		return nil, fmt.Errorf("初始化LLM失败: %v", err)
	}

	logger.Info("LLM初始化成功")

	return &Server{
		config:   cfg,
		logger:   logger,
		llmModel: llmModel,
	}, nil
}

// Start 启动Web服务器
func (s *Server) Start() error {
	// 1. 检查静态文件目录
	staticPath := filepath.Join("internal", "web", "static")
	s.logger.Info("静态文件目录路径: %s", staticPath)

	// 2. 创建静态文件服务器
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/", fs)
	s.logger.Info("静态文件服务器已创建")

	// 3. API路由
	http.HandleFunc("/api/chat", s.handleChat)
	s.logger.Info("API路由已注册: /api/chat")

	// 4. 启动服务器
	addr := fmt.Sprintf(":%d", s.config.Web.Port)
	s.logger.Info("Web服务器正在启动，地址 http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleChat 处理聊天API请求
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.logger.Warn("拒绝非POST请求: %s %s", r.Method, r.URL.Path)
		http.Error(w, "只接受POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.logger.Error("解析请求失败: %v", err)
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	// 检查消息是否为空
	if req.Message == "" {
		s.logger.Warn("接收到空消息请求")
		http.Error(w, "消息不能为空", http.StatusBadRequest)
		return
	}

	s.logger.Info("收到用户消息: %s", req.Message)

	// 准备消息
	messages := []llm.Message{
		{
			Role:    "user",
			Content: req.Message,
		},
	}

	// 调用LLM
	s.logger.Info("正在调用LLM...")
	resp, err := s.llmModel.Chat(messages)
	if err != nil {
		s.logger.Error("LLM调用失败: %v", err)
		writeJSON(w, ChatResponse{Error: fmt.Sprintf("LLM调用失败: %v", err)})
		return
	}

	// 返回响应
	s.logger.Info("LLM响应成功, 长度: %d字符", len(resp.Text))
	writeJSON(w, ChatResponse{Message: resp.Text})
}

// writeJSON 将响应写入HTTP响应
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
