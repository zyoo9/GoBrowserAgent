package browser

import (
	"GoBrowserAgent/internal/config"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Browser 浏览器控制接口
type Browser interface {
	// Navigate 导航到指定URL
	Navigate(url string) error

	// Screenshot 截取当前页面
	Screenshot(path string) error

	// Close 关闭浏览器
	Close() error

	// Find 查找元素
	Find(selector string) (Element, error)

	// WaitLoad 等待页面加载完成
	WaitLoad() error

	// GetTitle 获取页面标题
	GetTitle() (string, error)

	// GetURL 获取当前URL
	GetURL() string

	// ExecuteScript 执行JavaScript
	ExecuteScript(script string) error
}

// Element 页面元素接口
type Element interface {
	// Click 点击元素
	Click() error

	// Input 输入文本
	Input(text string) error

	// Clear 清空输入
	Clear() error

	// GetText 获取元素文本
	GetText() (string, error)

	// IsVisible 检查元素是否可见
	IsVisible() (bool, error)
}

// RodBrowser rod浏览器实现
type RodBrowser struct {
	browser *rod.Browser
	page    *rod.Page
	config  config.BrowserConfig
	ctx     context.Context
	cancel  context.CancelFunc
}

// RodElement rod元素实现
type RodElement struct {
	element *rod.Element
}

// NewBrowser 创建新的浏览器实例
func NewBrowser(cfg config.BrowserConfig) (Browser, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// 配置启动参数
	l := launcher.New().
		Headless(cfg.Headless).
		UserDataDir(cfg.UserDataDir)

	// 启动浏览器
	url := l.MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()

	// 创建页面
	page := browser.MustPage()

	return &RodBrowser{
		browser: browser,
		page:    page,
		config:  cfg,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

// Navigate 导航到指定URL
func (b *RodBrowser) Navigate(url string) error {
	err := b.page.Navigate(url)
	return err
}

// Screenshot 截取当前页面
func (b *RodBrowser) Screenshot(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	data, err := b.page.Screenshot(false, nil)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}

// Close 关闭浏览器
func (b *RodBrowser) Close() error {
	defer b.cancel()
	return b.browser.Close()
}

// Find 查找元素
func (b *RodBrowser) Find(selector string) (Element, error) {
	el, err := b.page.Element(selector)
	if err != nil {
		return nil, err
	}

	return &RodElement{element: el}, nil
}

// WaitLoad 等待页面加载完成
func (b *RodBrowser) WaitLoad() error {
	return b.page.WaitLoad()
}

// GetTitle 获取页面标题
func (b *RodBrowser) GetTitle() (string, error) {
	var title string
	err := rod.Try(func() {
		title = b.page.MustInfo().Title
	})
	return title, err
}

// GetURL 获取当前URL
func (b *RodBrowser) GetURL() string {
	info, err := b.page.Info()
	if err != nil {
		return ""
	}
	return info.URL
}

// ExecuteScript 执行JavaScript
func (b *RodBrowser) ExecuteScript(script string) error {
	_, err := b.page.Eval(script)
	return err
}

// Click 点击元素
func (e *RodElement) Click() error {
	err := rod.Try(func() {
		e.element.MustClick()
	})
	return err
}

// Input 输入文本
func (e *RodElement) Input(text string) error {
	err := rod.Try(func() {
		e.element.MustInput(text)
	})
	return err
}

// Clear 清空输入
func (e *RodElement) Clear() error {
	err := rod.Try(func() {
		e.element.MustSelectAllText().MustInput("")
	})
	return err
}

// GetText 获取元素文本
func (e *RodElement) GetText() (string, error) {
	var text string
	err := rod.Try(func() {
		text = e.element.MustText()
	})
	return text, err
}

// IsVisible 检查元素是否可见
func (e *RodElement) IsVisible() (bool, error) {
	var visible bool
	err := rod.Try(func() {
		visible = e.element.MustVisible()
	})
	return visible, err
}
