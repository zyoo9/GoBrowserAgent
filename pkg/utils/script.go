package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Script 表示一个批处理脚本
type Script struct {
	Commands []string
	Filename string
}

// LoadScript 从文件加载脚本
func LoadScript(filename string) (*Script, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开脚本文件 %s: %v", filename, err)
	}
	defer file.Close()

	script := &Script{
		Commands: make([]string, 0),
		Filename: filename,
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		script.Commands = append(script.Commands, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取脚本文件 %s 时出错: %v", filename, err)
	}

	return script, nil
}

// ScriptExecutor 脚本执行器接口
type ScriptExecutor interface {
	// ExecuteCommand 执行单个命令
	ExecuteCommand(command string) error
}

// ExecuteScript 执行脚本
func ExecuteScript(script *Script, executor ScriptExecutor, verbose bool) error {
	if verbose {
		fmt.Printf("执行脚本: %s (%d 条命令)\n", script.Filename, len(script.Commands))
	}

	for i, cmd := range script.Commands {
		if verbose {
			fmt.Printf("[%d/%d] 执行: %s\n", i+1, len(script.Commands), cmd)
		}

		if err := executor.ExecuteCommand(cmd); err != nil {
			return fmt.Errorf("执行命令 '%s' 时出错: %v", cmd, err)
		}
	}

	if verbose {
		fmt.Println("脚本执行完成")
	}
	return nil
}
