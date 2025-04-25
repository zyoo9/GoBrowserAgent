# Go浏览器代理（GoBrowserAgent）

Go浏览器代理是一个基于Go语言开发的命令行工具，用于自动化浏览器操作。它使用[Rod](https://github.com/go-rod/rod)库提供浏览器自动化能力，可以通过简单的命令语法执行各种浏览器操作。

## 功能特点

- 命令行交互界面
- 浏览器自动化操作
- 灵活的任务模型
- 简单的指令解析
- 支持批处理脚本
- 可配置的运行环境
- 内置Web界面，支持与LLM对话

## 支持的命令

- `go <url>` 或 `navigate <url>` - 导航到指定URL
- `login <url> username=<username> password=<password>` - 登录到指定网站
- `screenshot [path=<filename>]` - 截取当前页面
- `wait <duration>` - 等待指定时间（如5s, 1m）
- `search <url> query=<query> input=<selector> button=<selector>` - 在网站上执行搜索
- `form <url> [form=<selector>] field_<n>=<value>... [submit=<selector>]` - 填写表单
- `script "<javascript-code>"` - 执行JavaScript代码
- `exit` 或 `quit` - 退出程序
- `help` - 显示帮助信息

## 安装

确保已安装Go 1.16或更高版本。

```bash
# 克隆仓库
git clone https://github.com/yourusername/GoBrowserAgent.git
cd GoBrowserAgent

# 安装依赖
go mod download

# 构建项目
go build
```

## 命令行参数

GoBrowserAgent支持以下命令行参数：

- `-script <path>`: 要执行的脚本文件路径
- `-verbose`: 显示详细的执行信息
- `-config <path>`: 配置文件路径
- `-web`: 启动Web界面模式，提供LLM对话功能

示例：

```bash
# 交互模式
./GoBrowserAgent

# 脚本模式
./GoBrowserAgent -script examples/baidu_search.txt -verbose

# 使用配置文件
./GoBrowserAgent -config my_config.json

# 启动Web界面模式
./GoBrowserAgent -web
```

## Web界面

GoBrowserAgent支持Web界面模式，提供与LLM对话的功能。

启动Web模式：

```bash
./GoBrowserAgent -web
```

然后在浏览器中访问：http://localhost:8080

在Web界面中，您可以：
- 与大型语言模型进行对话
- 获取智能问答和建议
- 通过输入问题并点击发送或按Enter键来与模型交互

Web界面配置可以在config.json中的web部分进行设置：

```json
"web": {
  "port": 8080  // Web服务器端口
}
```

## 脚本文件

脚本文件是一系列要按顺序执行的命令，每行一个命令。支持注释（使用`#`或`//`开头的行）和空行。

脚本示例：

```
# 百度搜索示例
go https://www.baidu.com
wait 2s
search https://www.baidu.com query=Golang开发 input=#kw button=#su
wait 3s
screenshot path=baidu_search_result.png
```

您可以参考`examples/`目录下的示例脚本。

## 配置文件

配置文件使用JSON格式，允许您自定义浏览器、日志和LLM设置。

配置示例：

```json
{
  "browser": {
    "headless": false,
    "user_data_dir": "./user_data",
    "default_width": 1280,
    "default_height": 800,
    "timeout": 30
  },
  "llm": {
    "provider": "openai",
    "api_key": "your-api-key-here",
    "api_url": "https://api.openai.com/v1",
    "model_name": "gpt-3.5-turbo",
    "max_tokens": 2048,
    "system_prompt": "你是一个专注于Go编程语言的技术助手，擅长解释代码并提供编程建议。"
  },
  "log": {
    "level": "info",
    "output_path": "./logs"
  },
  "web": {
    "port": 8080
  }
}
```

复制`config.json.example`文件并根据需要进行修改。

### 配置不同的LLM提供商

GoBrowserAgent支持各种LLM提供商，只需在配置文件中相应调整即可：

1. **OpenAI (GPT)**
```json
"llm": {
  "provider": "openai",
  "api_key": "sk-your-openai-key",
  "api_url": "https://api.openai.com/v1",
  "model_name": "gpt-3.5-turbo",
  "max_tokens": 2048
}
```

2. **通义千问 (Qwen)**
```json
"llm": {
  "provider": "qwen",
  "api_key": "your-qwen-key",
  "api_url": "https://dashscope.aliyuncs.com/api/v1",
  "model_name": "qwen-max",
  "max_tokens": 2048
}
```

3. **百度文心一言 (ERNIE Bot)**
```json
"llm": {
  "provider": "ernie",
  "api_key": "your-ernie-key",
  "api_url": "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat",
  "model_name": "ernie-bot-4",
  "max_tokens": 2048
}
```

4. **讯飞星火 (Spark)**
```json
"llm": {
  "provider": "spark",
  "api_key": "your-spark-key",
  "api_url": "https://spark-api.xf-yun.com/v3.1",
  "model_name": "spark-3.5",
  "max_tokens": 2048
}
```

您可以根据自己使用的模型来调整配置，系统会自动适配不同模型的响应格式。

### 配置系统提示（System Prompt）

系统提示是向LLM提供的初始指令，用于设置模型的行为和回答风格。在配置文件中，您可以通过`system_prompt`字段来自定义它：

```json
"llm": {
  "provider": "openai",
  "api_key": "your-api-key",
  "api_url": "https://api.openai.com/v1",
  "model_name": "gpt-3.5-turbo",
  "max_tokens": 2048,
  "system_prompt": "你是一个专注于Go编程语言的技术助手，擅长解释代码并提供编程建议。"
}
```

系统提示的一些使用建议：

1. **指定角色和专业领域**：例如"你是一个Go语言专家"或"你是一个网络安全顾问"
2. **设定回应风格**：例如"以简洁明了的方式回答问题"或"提供详细且有教育意义的回答"
3. **限制回应范围**：例如"仅回答与编程相关的问题"
4. **针对特定任务定制**：例如"帮助用户理解浏览器自动化的原理和实践"

合理设置系统提示可以让模型更好地满足您的需求，提供更加精准的回答。

## 使用示例

### 导航到网页

```
> go https://www.baidu.com
已导航到: https://www.baidu.com
```

### 在百度搜索

```
> search https://www.baidu.com query=golang input=#kw button=#su
搜索成功!
```

### 截图

```
> screenshot path=baidu.png
截图已保存到: baidu.png
```

### 使用脚本修改页面

```
> script "document.querySelector('.title').style.color = 'red';"
脚本执行成功!
```

### 填写表单

```
> form https://example.com field_username=user field_password=pass submit=#login-button
表单填写成功!
```

## 项目结构

```
GoBrowserAgent/
├── internal/          # 内部包
│   ├── config/        # 配置管理
│   ├── log/           # 日志工具
│   └── web/           # Web服务器实现
│       └── static/    # 静态网页文件
├── pkg/               # 核心包
│   ├── browser/       # 浏览器控制接口和实现
│   ├── parser/        # 命令解析器
│   ├── task/          # 任务接口和实现
│   ├── llm/           # 语言模型接口
│   └── utils/         # 通用工具函数
├── examples/          # 示例脚本
├── main.go            # 主程序入口
├── go.mod             # Go模块定义
└── go.sum             # 依赖校验和
```

## 下一步计划

- 支持更多的浏览器操作
- 添加更多的选择器类型（XPath、CSS等）
- 实现脚本变量和条件逻辑
- 实现更高级的命令解析（例如，使用大型语言模型）
- 添加更多的输出格式选项
- 支持并发任务执行
- 增强Web界面功能

## 贡献

欢迎贡献代码、报告问题或提出建议！

## 许可证

MIT 