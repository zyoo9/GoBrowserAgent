package task

import (
	"GoBrowserAgent/pkg/browser"
	"fmt"
	"time"
)

// Task 表示一个待执行的任务
type Task interface {
	// Execute 执行任务
	Execute() error

	// GetName 获取任务名称
	GetName() string

	// GetDescription 获取任务描述
	GetDescription() string
}

// BaseTask 基础任务结构
type BaseTask struct {
	Name        string
	Description string
	Browser     browser.Browser
}

// GetName 获取任务名称
func (t *BaseTask) GetName() string {
	return t.Name
}

// GetDescription 获取任务描述
func (t *BaseTask) GetDescription() string {
	return t.Description
}

// LoginTask 登录任务
type LoginTask struct {
	*BaseTask
	URL      string
	Username string
	Password string
}

// NewLoginTask 创建登录任务
func NewLoginTask(b browser.Browser, url, username, password string) *LoginTask {
	return &LoginTask{
		BaseTask: &BaseTask{
			Name:        "登录任务",
			Description: "自动登录到指定系统",
			Browser:     b,
		},
		URL:      url,
		Username: username,
		Password: password,
	}
}

// Execute 执行登录任务
func (t *LoginTask) Execute() error {
	// 1. 导航到登录页面
	if err := t.Browser.Navigate(t.URL); err != nil {
		return err
	}

	// 2. 等待页面加载
	if err := t.Browser.WaitLoad(); err != nil {
		return err
	}

	// 3. 查找用户名输入框并输入
	userInput, err := t.Browser.Find("input[name='username']")
	if err != nil {
		return err
	}

	if err := userInput.Clear(); err != nil {
		return err
	}

	if err := userInput.Input(t.Username); err != nil {
		return err
	}

	// 4. 查找密码输入框并输入
	passInput, err := t.Browser.Find("input[name='password']")
	if err != nil {
		return err
	}

	if err := passInput.Clear(); err != nil {
		return err
	}

	if err := passInput.Input(t.Password); err != nil {
		return err
	}

	// 5. 查找登录按钮并点击
	loginBtn, err := t.Browser.Find("button[type='submit']")
	if err != nil {
		return err
	}

	if err := loginBtn.Click(); err != nil {
		return err
	}

	// 6. 等待页面加载
	return t.Browser.WaitLoad()
}

// SearchTask 搜索任务
type SearchTask struct {
	*BaseTask
	URL          string
	SearchQuery  string
	SearchInput  string
	SearchButton string
}

// NewSearchTask 创建搜索任务
func NewSearchTask(b browser.Browser, url, query, inputSelector, buttonSelector string) *SearchTask {
	return &SearchTask{
		BaseTask: &BaseTask{
			Name:        "搜索任务",
			Description: "在网站中执行搜索",
			Browser:     b,
		},
		URL:          url,
		SearchQuery:  query,
		SearchInput:  inputSelector,
		SearchButton: buttonSelector,
	}
}

// Execute 执行搜索任务
func (t *SearchTask) Execute() error {
	// 1. 导航到页面
	if err := t.Browser.Navigate(t.URL); err != nil {
		return err
	}

	// 2. 等待页面加载
	if err := t.Browser.WaitLoad(); err != nil {
		return err
	}

	// 3. 查找搜索输入框并输入
	searchInput, err := t.Browser.Find(t.SearchInput)
	if err != nil {
		return err
	}

	if err := searchInput.Clear(); err != nil {
		return err
	}

	if err := searchInput.Input(t.SearchQuery); err != nil {
		return err
	}

	// 4. 点击搜索按钮
	searchBtn, err := t.Browser.Find(t.SearchButton)
	if err != nil {
		return err
	}

	if err := searchBtn.Click(); err != nil {
		return err
	}

	// 5. 等待页面加载
	return t.Browser.WaitLoad()
}

// FormFillTask 表单填写任务
type FormFillTask struct {
	*BaseTask
	URL          string
	FormSelector string
	Fields       map[string]string
	SubmitButton string
}

// NewFormFillTask 创建表单填写任务
func NewFormFillTask(b browser.Browser, url, formSelector string, fields map[string]string, submitButton string) *FormFillTask {
	return &FormFillTask{
		BaseTask: &BaseTask{
			Name:        "表单填写任务",
			Description: "自动填写网页表单",
			Browser:     b,
		},
		URL:          url,
		FormSelector: formSelector,
		Fields:       fields,
		SubmitButton: submitButton,
	}
}

// Execute 执行表单填写任务
func (t *FormFillTask) Execute() error {
	// 1. 导航到页面
	if err := t.Browser.Navigate(t.URL); err != nil {
		return err
	}

	// 2. 等待页面加载
	if err := t.Browser.WaitLoad(); err != nil {
		return err
	}

	// 3. 填写表单字段
	for selector, value := range t.Fields {
		// 查找字段
		field, err := t.Browser.Find(selector)
		if err != nil {
			return fmt.Errorf("找不到字段 %s: %v", selector, err)
		}

		// 清空字段
		if err := field.Clear(); err != nil {
			return fmt.Errorf("清空字段 %s 失败: %v", selector, err)
		}

		// 输入值
		if err := field.Input(value); err != nil {
			return fmt.Errorf("在字段 %s 中输入值失败: %v", selector, err)
		}
	}

	// 4. 提交表单
	if t.SubmitButton != "" {
		submitBtn, err := t.Browser.Find(t.SubmitButton)
		if err != nil {
			return fmt.Errorf("找不到提交按钮 %s: %v", t.SubmitButton, err)
		}

		if err := submitBtn.Click(); err != nil {
			return fmt.Errorf("点击提交按钮失败: %v", err)
		}
	}

	// 5. 等待页面加载
	return t.Browser.WaitLoad()
}

// NavigationTask 导航任务
type NavigationTask struct {
	*BaseTask
	URL string
}

// NewNavigationTask 创建导航任务
func NewNavigationTask(b browser.Browser, url string) *NavigationTask {
	return &NavigationTask{
		BaseTask: &BaseTask{
			Name:        "导航任务",
			Description: "导航到指定URL",
			Browser:     b,
		},
		URL: url,
	}
}

// Execute 执行导航任务
func (t *NavigationTask) Execute() error {
	// 导航到页面
	if err := t.Browser.Navigate(t.URL); err != nil {
		return err
	}

	// 等待页面加载
	return t.Browser.WaitLoad()
}

// ScreenshotTask 截图任务
type ScreenshotTask struct {
	*BaseTask
	Path string
}

// NewScreenshotTask 创建截图任务
func NewScreenshotTask(b browser.Browser, path string) *ScreenshotTask {
	return &ScreenshotTask{
		BaseTask: &BaseTask{
			Name:        "截图任务",
			Description: "截取当前页面",
			Browser:     b,
		},
		Path: path,
	}
}

// Execute 执行截图任务
func (t *ScreenshotTask) Execute() error {
	// 截取当前页面
	return t.Browser.Screenshot(t.Path)
}

// ScriptTask 脚本执行任务
type ScriptTask struct {
	*BaseTask
	Script string
}

// NewScriptTask 创建脚本执行任务
func NewScriptTask(b browser.Browser, script string) *ScriptTask {
	return &ScriptTask{
		BaseTask: &BaseTask{
			Name:        "脚本执行任务",
			Description: "执行JavaScript脚本",
			Browser:     b,
		},
		Script: script,
	}
}

// Execute 执行脚本任务
func (t *ScriptTask) Execute() error {
	// 执行JavaScript脚本
	return t.Browser.ExecuteScript(t.Script)
}

// WaitTask 等待任务
type WaitTask struct {
	*BaseTask
	Duration time.Duration
}

// NewWaitTask 创建等待任务
func NewWaitTask(b browser.Browser, duration time.Duration) *WaitTask {
	return &WaitTask{
		BaseTask: &BaseTask{
			Name:        "等待任务",
			Description: "等待指定时间",
			Browser:     b,
		},
		Duration: duration,
	}
}

// Execute 执行等待任务
func (t *WaitTask) Execute() error {
	// 等待指定时间
	time.Sleep(t.Duration)
	return nil
}

// TaskManager 任务管理器
type TaskManager struct {
	tasks []Task
}

// NewTaskManager 创建任务管理器
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks: make([]Task, 0),
	}
}

// AddTask 添加任务
func (tm *TaskManager) AddTask(task Task) {
	tm.tasks = append(tm.tasks, task)
}

// ExecuteTasks 执行所有任务
func (tm *TaskManager) ExecuteTasks() []error {
	var errors []error

	for _, task := range tm.tasks {
		if err := task.Execute(); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// ClearTasks 清空任务列表
func (tm *TaskManager) ClearTasks() {
	tm.tasks = make([]Task, 0)
}
