package parser

import (
	"strings"
	"time"
)

// Command 表示解析后的命令
type Command struct {
	Action  string
	Target  string
	Params  map[string]string
	RawText string
}

// Parser 指令解析器接口
type Parser interface {
	// Parse 解析指令文本
	Parse(text string) (*Command, error)
}

// SimpleParser 简单指令解析器
type SimpleParser struct{}

// NewSimpleParser 创建简单指令解析器
func NewSimpleParser() Parser {
	return &SimpleParser{}
}

// Parse 解析指令文本
func (p *SimpleParser) Parse(text string) (*Command, error) {
	// 创建命令对象
	cmd := &Command{
		RawText: text,
		Params:  make(map[string]string),
	}

	// 分割文本
	parts := strings.Split(text, " ")
	if len(parts) == 0 {
		return nil, nil
	}

	// 解析动作
	cmd.Action = strings.ToLower(parts[0])

	// 解析目标（如有）
	if len(parts) > 1 && !strings.Contains(parts[1], "=") {
		cmd.Target = parts[1]
	}

	// 解析参数
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				cmd.Params[kv[0]] = kv[1]
			}
		}
	}

	// 处理特殊参数
	p.processSpecialParams(cmd)

	return cmd, nil
}

// 处理特殊参数
func (p *SimpleParser) processSpecialParams(cmd *Command) {
	// 处理时间参数
	if timeStr, ok := cmd.Params["wait"]; ok {
		// 移除空格
		timeStr = strings.TrimSpace(timeStr)

		// 默认单位为秒
		if _, err := time.ParseDuration(timeStr); err != nil {
			if num := strings.TrimSuffix(timeStr, "s"); num == timeStr {
				cmd.Params["wait"] = timeStr + "s"
			}
		}
	}

	// 处理表单字段参数
	if cmd.Action == "form" {
		fields := make(map[string]string)

		for k, v := range cmd.Params {
			// 以field_前缀的参数都是表单字段
			if strings.HasPrefix(k, "field_") {
				fieldName := strings.TrimPrefix(k, "field_")
				fields[fieldName] = v
				delete(cmd.Params, k) // 从原始参数中移除
			}
		}

		// 将所有字段序列化为JSON格式存储在一个参数中
		if len(fields) > 0 {
			cmd.Params["fields"] = formatFieldsMap(fields)
		}
	}
}

// 格式化字段映射为字符串
func formatFieldsMap(fields map[string]string) string {
	parts := make([]string, 0, len(fields))
	for k, v := range fields {
		parts = append(parts, k+":"+v)
	}
	return strings.Join(parts, ",")
}

// CommandValidator 命令验证器
type CommandValidator struct {
	RequiredParams map[string][]string
}

// NewCommandValidator 创建命令验证器
func NewCommandValidator() *CommandValidator {
	validator := &CommandValidator{
		RequiredParams: make(map[string][]string),
	}

	// 初始化各命令的必要参数
	validator.RequiredParams["login"] = []string{"username", "password"}
	validator.RequiredParams["search"] = []string{"query"}
	validator.RequiredParams["form"] = []string{"url"}
	validator.RequiredParams["wait"] = []string{"duration"}

	return validator
}

// Validate 验证命令参数
func (v *CommandValidator) Validate(cmd *Command) []string {
	var missing []string

	// 检查该命令是否有必要参数定义
	requiredParams, exists := v.RequiredParams[cmd.Action]
	if !exists {
		return missing
	}

	// 检查必要参数是否都存在
	for _, param := range requiredParams {
		if _, ok := cmd.Params[param]; !ok && param != "url" {
			missing = append(missing, param)
		}
	}

	// 特殊处理url参数，可以是Target或Params["url"]
	if contains(requiredParams, "url") && cmd.Target == "" && cmd.Params["url"] == "" {
		missing = append(missing, "url")
	}

	return missing
}

// 检查slice中是否包含特定字符串
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// LLMParser 基于大型语言模型的解析器
type LLMParser struct {
	// TODO: 接入LLM API
	SimpleParser Parser
}

// NewLLMParser 创建LLM解析器
func NewLLMParser() Parser {
	return &LLMParser{
		SimpleParser: NewSimpleParser(),
	}
}

// Parse 解析指令文本
func (p *LLMParser) Parse(text string) (*Command, error) {
	// TODO: 实现LLM解析逻辑

	// 临时使用简单解析器
	return p.SimpleParser.Parse(text)
}
