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
	llmModel, err := llm.NewLLM(cfg.LLM)
	if err != nil {
		return nil, fmt.Errorf("初始化LLM失败: %v", err)
	}

	return &Server{
		config:   cfg,
		logger:   logger,
		llmModel: llmModel,
	}, nil
}

// Start 启动Web服务器
func (s *Server) Start() error {
	// 1. 静态文件服务
	fs := http.FileServer(http.Dir(filepath.Join("internal", "web", "static")))
	http.Handle("/", fs)

	// 2. API路由
	http.HandleFunc("/api/chat", s.handleChat)

	// 3. 启动服务器
	addr := fmt.Sprintf(":%d", s.config.Web.Port)
	s.logger.Info("Web服务器正在启动，地址 http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleChat 处理聊天API请求
func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只接受POST请求", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	// 检查消息是否为空
	if req.Message == "" {
		http.Error(w, "消息不能为空", http.StatusBadRequest)
		return
	}

	// 准备消息
	messages := []llm.Message{
		{
			Role:    "user",
			Content: req.Message,
		},
	}

	// 调用LLM
	resp, err := s.llmModel.Chat(messages)
	if err != nil {
		s.logger.Error("LLM调用失败: %v", err)
		writeJSON(w, ChatResponse{Error: fmt.Sprintf("LLM调用失败: %v", err)})
		return
	}

	// 返回响应
	writeJSON(w, ChatResponse{Message: resp.Text})
}

// writeJSON 将响应写入HTTP响应
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
