package main

import (
	"GoBrowserAgent/internal/config"
	"GoBrowserAgent/internal/log"
	"GoBrowserAgent/internal/web"
	"GoBrowserAgent/pkg/browser"
	"GoBrowserAgent/pkg/parser"
	"GoBrowserAgent/pkg/task"
	"GoBrowserAgent/pkg/utils"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// 全局变量
var (
	cfg    *config.Config
	logger log.Logger
)

// CommandExecutor 实现脚本执行器接口
type CommandExecutor struct {
	Parser      parser.Parser
	Browser     browser.Browser
	Validator   *parser.CommandValidator
	Logger      log.Logger
	Interactive bool
}

// ExecuteCommand 执行单个命令
func (ce *CommandExecutor) ExecuteCommand(command string) error {
	// 跳过空命令
	if strings.TrimSpace(command) == "" {
		return nil
	}

	// 解析命令
	cmd, err := ce.Parser.Parse(command)
	if err != nil {
		ce.Logger.Error("解析命令失败: %v", err)
		return fmt.Errorf("解析命令失败: %v", err)
	}

	// 处理命令
	switch cmd.Action {
	case "login":
		url := cmd.Target
		username := cmd.Params["username"]
		password := cmd.Params["password"]

		// 参数验证
		missing := ce.Validator.Validate(cmd)
		if len(missing) > 0 || url == "" {
			return fmt.Errorf("login命令缺少参数: url=%v, username=%v, password=%v", url == "", contains(missing, "username"), contains(missing, "password"))
		}

		// 创建并执行登录任务
		loginTask := task.NewLoginTask(ce.Browser, url, username, password)
		if err := loginTask.Execute(); err != nil {
			ce.Logger.Error("登录失败: %v", err)
			return fmt.Errorf("登录失败: %v", err)
		} else {
			ce.Logger.Info("成功登录到 %s", url)
			if ce.Interactive {
				fmt.Println("登录成功!")
			}
		}

	case "screenshot":
		path := cmd.Params["path"]
		if path == "" {
			path = "screenshot.png"
		}

		// 截图
		screenshotTask := task.NewScreenshotTask(ce.Browser, path)
		if err := screenshotTask.Execute(); err != nil {
			ce.Logger.Error("截图失败: %v", err)
			return fmt.Errorf("截图失败: %v", err)
		} else {
			ce.Logger.Info("截图已保存到 %s", path)
			if ce.Interactive {
				fmt.Println("截图已保存到:", path)
			}
		}

	case "go", "navigate":
		url := cmd.Target
		if url == "" {
			url = cmd.Params["url"]
		}

		if url == "" {
			return fmt.Errorf("缺少URL参数")
		}

		// 导航
		navTask := task.NewNavigationTask(ce.Browser, url)
		if err := navTask.Execute(); err != nil {
			ce.Logger.Error("导航失败: %v", err)
			return fmt.Errorf("导航失败: %v", err)
		} else {
			ce.Logger.Info("已导航到 %s", url)
			if ce.Interactive {
				fmt.Println("已导航到:", url)
			}
		}

	case "wait":
		durationStr := cmd.Params["duration"]
		if durationStr == "" {
			durationStr = cmd.Target
		}

		if durationStr == "" {
			return fmt.Errorf("缺少持续时间参数")
		}

		// 解析持续时间
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			// 尝试添加默认单位（秒）
			if num, err := strconv.Atoi(durationStr); err == nil {
				duration = time.Duration(num) * time.Second
			} else {
				ce.Logger.Error("无效的持续时间: %v", durationStr)
				return fmt.Errorf("无效的持续时间: %v", durationStr)
			}
		}

		// 等待
		waitTask := task.NewWaitTask(ce.Browser, duration)
		if ce.Interactive {
			fmt.Printf("等待 %v...\n", duration)
		}
		if err := waitTask.Execute(); err != nil {
			ce.Logger.Error("等待失败: %v", err)
			return fmt.Errorf("等待失败: %v", err)
		} else {
			ce.Logger.Info("等待完成: %v", duration)
			if ce.Interactive {
				fmt.Println("等待完成")
			}
		}

	case "search":
		url := cmd.Target
		if url == "" {
			url = cmd.Params["url"]
		}

		query := cmd.Params["query"]
		input := cmd.Params["input"]
		button := cmd.Params["button"]

		// 参数验证
		missing := ce.Validator.Validate(cmd)
		if len(missing) > 0 || url == "" || input == "" || button == "" {
			return fmt.Errorf("search命令缺少参数: url=%v, query=%v, input=%v, button=%v",
				url == "", contains(missing, "query"), input == "", button == "")
		}

		// 创建并执行搜索任务
		searchTask := task.NewSearchTask(ce.Browser, url, query, input, button)
		if err := searchTask.Execute(); err != nil {
			ce.Logger.Error("搜索失败: %v", err)
			return fmt.Errorf("搜索失败: %v", err)
		} else {
			ce.Logger.Info("在 %s 上搜索 %s 成功", url, query)
			if ce.Interactive {
				fmt.Println("搜索成功!")
			}
		}

	case "form":
		url := cmd.Target
		if url == "" {
			url = cmd.Params["url"]
		}

		formSelector := cmd.Params["form"]
		fieldsStr := cmd.Params["fields"]
		submitButton := cmd.Params["submit"]

		// 参数验证
		if url == "" {
			return fmt.Errorf("表单填写缺少URL参数")
		}

		// 解析字段
		fields := make(map[string]string)
		if fieldsStr != "" {
			fieldPairs := strings.Split(fieldsStr, ",")
			for _, pair := range fieldPairs {
				kv := strings.SplitN(pair, ":", 2)
				if len(kv) == 2 {
					fields[kv[0]] = kv[1]
				}
			}
		}

		if len(fields) == 0 {
			return fmt.Errorf("表单填写至少需要一个字段")
		}

		// 创建并执行表单填写任务
		formTask := task.NewFormFillTask(ce.Browser, url, formSelector, fields, submitButton)
		if err := formTask.Execute(); err != nil {
			ce.Logger.Error("表单填写失败: %v", err)
			return fmt.Errorf("表单填写失败: %v", err)
		} else {
			ce.Logger.Info("在 %s 上填写表单成功", url)
			if ce.Interactive {
				fmt.Println("表单填写成功!")
			}
		}

	case "script":
		scriptCode := cmd.Target
		if scriptCode == "" {
			scriptCode = cmd.Params["code"]
		}

		if scriptCode == "" {
			return fmt.Errorf("缺少脚本代码")
		}

		// 创建并执行脚本任务
		scriptTask := task.NewScriptTask(ce.Browser, scriptCode)
		if err := scriptTask.Execute(); err != nil {
			ce.Logger.Error("脚本执行失败: %v", err)
			return fmt.Errorf("脚本执行失败: %v", err)
		} else {
			ce.Logger.Info("脚本执行成功")
			if ce.Interactive {
				fmt.Println("脚本执行成功!")
			}
		}

	case "exit", "quit":
		if ce.Interactive {
			fmt.Println("退出程序")
		}
		return nil

	case "help":
		if ce.Interactive {
			printHelp()
		}

	default:
		return fmt.Errorf("未知命令: %s", cmd.Action)
	}

	return nil
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("可用命令:")
	fmt.Println("  go <url> - 导航到指定URL")
	fmt.Println("  navigate <url> - 导航到指定URL")
	fmt.Println("  login <url> username=<username> password=<password> - 登录到指定网站")
	fmt.Println("  screenshot [path=<filename>] - 截取当前页面")
	fmt.Println("  wait <duration> - 等待指定时间（如5s, 1m）")
	fmt.Println("  search <url> query=<query> input=<selector> button=<selector> - 在网站上执行搜索")
	fmt.Println("  form <url> [form=<selector>] field_<name>=<value>... [submit=<selector>] - 填写表单")
	fmt.Println("  script \"<javascript-code>\" - 执行JavaScript代码")
	fmt.Println("  exit/quit - 退出程序")
	fmt.Println("  help - 显示帮助信息")
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

func main() {
	// 命令行参数解析
	scriptFile := flag.String("script", "", "要执行的脚本文件路径")
	verbose := flag.Bool("verbose", false, "显示详细的执行信息")
	configFile := flag.String("config", "", "配置文件路径")
	webMode := flag.Bool("web", false, "启动Web界面模式")
	flag.Parse()

	// 1. 初始化配置
	if *configFile != "" {
		var err error
		cfg, err = config.LoadConfig(*configFile)
		if err != nil {
			fmt.Printf("加载配置文件失败: %v，将使用默认配置\n", err)
			cfg = config.DefaultConfig()
		}
	} else {
		cfg = config.DefaultConfig()
	}

	// 2. 初始化日志
	var err error
	logger, err = log.NewLogger(cfg.Log.Level, cfg.Log.OutputPath)
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}

	// 3. 如果是Web模式，启动Web服务器
	if *webMode {
		server, err := web.NewServer(cfg, logger)
		if err != nil {
			logger.Fatal("创建Web服务器失败: %v", err)
		}

		fmt.Printf("启动Web模式，访问 http://localhost:%d 开始对话\n", cfg.Web.Port)
		if err := server.Start(); err != nil {
			logger.Fatal("Web服务器启动失败: %v", err)
		}
		return
	}

	// 4. 启动浏览器（仅在非Web模式下）
	b, err := browser.NewBrowser(cfg.Browser)
	if err != nil {
		logger.Fatal("启动浏览器失败: %v", err)
	}
	defer b.Close()

	// 5. 创建指令解析器
	p := parser.NewSimpleParser()

	// 6. 创建命令验证器
	validator := parser.NewCommandValidator()

	// 7. 创建命令执行器
	executor := &CommandExecutor{
		Parser:      p,
		Browser:     b,
		Validator:   validator,
		Logger:      logger,
		Interactive: *scriptFile == "", // 脚本模式下不交互
	}

	// 8. 如果提供了脚本文件，执行脚本
	if *scriptFile != "" {
		script, err := utils.LoadScript(*scriptFile)
		if err != nil {
			logger.Fatal("加载脚本失败: %v", err)
		}

		err = utils.ExecuteScript(script, executor, *verbose)
		if err != nil {
			logger.Fatal("执行脚本失败: %v", err)
		}
		return
	}

	// 9. 交互循环
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("欢迎使用Go浏览器代理！输入 'help' 获取帮助。")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("读取输入失败: %v", err)
			continue
		}

		// 去除空白字符
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 退出命令
		if input == "exit" || input == "quit" {
			break
		}

		// 执行命令
		err = executor.ExecuteCommand(input)
		if err != nil {
			fmt.Println("错误:", err)
		}
	}

	fmt.Println("再见!")
}
