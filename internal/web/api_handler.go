package web

import (
	"encoding/json"
	"net/http"

	"GoBrowserAgent/internal/service/llm"

	"github.com/sirupsen/logrus"
)

// UserChatRequest 定义用户请求结构
type UserChatRequest struct {
	Message string `json:"message"`
}

// UserChatResponse 定义响应给用户的结构
type UserChatResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// APIHandler 处理API请求
type APIHandler struct {
	LLMService *llm.Service
}

// NewAPIHandler 创建新的API处理程序
func NewAPIHandler(llmService *llm.Service) *APIHandler {
	return &APIHandler{
		LLMService: llmService,
	}
}

// RegisterHandlers 注册HTTP处理程序
func (h *APIHandler) RegisterHandlers() {
	http.HandleFunc("/api/chat", h.handleChat)
}

// handleChat 处理聊天请求
func (h *APIHandler) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST请求", http.StatusMethodNotAllowed)
		return
	}

	var req UserChatRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		logrus.Errorf("解析请求体失败: %v", err)
		http.Error(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "消息不能为空", http.StatusBadRequest)
		return
	}

	// 处理聊天请求
	message, err := h.LLMService.Chat(req.Message)
	var resp UserChatResponse

	if err != nil {
		logrus.Errorf("处理聊天请求失败: %v", err)
		resp = UserChatResponse{
			Error: err.Error(),
		}
	} else {
		resp = UserChatResponse{
			Message: message,
		}
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		logrus.Errorf("编码响应失败: %v", err)
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
}
