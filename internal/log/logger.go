package log

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger 日志记录器接口
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

// LogLevel 日志级别
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// SimpleLogger 简单日志实现
type SimpleLogger struct {
	level    LogLevel
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger
	logFile  *os.File
}

// NewLogger 创建新的日志记录器
func NewLogger(level string, outputPath string) (Logger, error) {
	// 创建日志目录
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return nil, err
	}

	// 创建日志文件
	logFileName := filepath.Join(outputPath, time.Now().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// 创建日志记录器
	debugLog := log.New(logFile, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(logFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLog := log.New(logFile, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(logFile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLog := log.New(logFile, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)

	// 解析日志级别
	var logLevel LogLevel
	switch level {
	case "debug":
		logLevel = DebugLevel
	case "info":
		logLevel = InfoLevel
	case "warn":
		logLevel = WarnLevel
	case "error":
		logLevel = ErrorLevel
	case "fatal":
		logLevel = FatalLevel
	default:
		logLevel = InfoLevel
	}

	return &SimpleLogger{
		level:    logLevel,
		debugLog: debugLog,
		infoLog:  infoLog,
		warnLog:  warnLog,
		errorLog: errorLog,
		fatalLog: fatalLog,
		logFile:  logFile,
	}, nil
}

// Debug 调试日志
func (l *SimpleLogger) Debug(format string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.debugLog.Printf(format, args...)
	}
}

// Info 信息日志
func (l *SimpleLogger) Info(format string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.infoLog.Printf(format, args...)
	}
}

// Warn 警告日志
func (l *SimpleLogger) Warn(format string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.warnLog.Printf(format, args...)
	}
}

// Error 错误日志
func (l *SimpleLogger) Error(format string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.errorLog.Printf(format, args...)
	}
}

// Fatal 致命错误日志
func (l *SimpleLogger) Fatal(format string, args ...interface{}) {
	if l.level <= FatalLevel {
		l.fatalLog.Printf(format, args...)
		os.Exit(1)
	}
}
