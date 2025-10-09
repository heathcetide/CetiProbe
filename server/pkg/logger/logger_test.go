package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func makeTmpLogFile(t *testing.T, name string) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, name)
}

func readFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	return string(b)
}

func waitForWrite(assert func() bool) bool {
	deadline := time.Now().Add(800 * time.Millisecond)
	for time.Now().Before(deadline) {
		if assert() {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return assert()
}

func TestLogConfig(t *testing.T) {
	cfg := &LogConfig{
		Level:      "info",
		Filename:   "test.log",
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 3,
		Daily:      true,
	}

	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "test.log", cfg.Filename)
	assert.Equal(t, 100, cfg.MaxSize)
	assert.Equal(t, 7, cfg.MaxAge)
	assert.Equal(t, 3, cfg.MaxBackups)
	assert.True(t, cfg.Daily)
}

func TestInit(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	cfg := &LogConfig{
		Level:      "debug",
		Filename:   logFile,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Daily:      false,
	}

	// 测试开发模式
	err := Init(cfg, "dev")
	require.NoError(t, err)
	assert.NotNil(t, lg)

	// 测试日志输出
	Info("test info message")
	Debug("test debug message")
	Warn("test warn message")
	Error("test error message")

	// 测试生产模式
	err = Init(cfg, "prod")
	require.NoError(t, err)

	Info("test production message")
}

func TestInitWithInvalidLevel(t *testing.T) {
	cfg := &LogConfig{
		Level:      "invalid",
		Filename:   "test.log",
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Daily:      false,
	}

	err := Init(cfg, "prod")
	assert.Error(t, err)
}

func TestGetLogWriter(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	// 测试普通日志写入器
	writer := getLogWriter(logFile, 10, 1, 1, false)
	assert.NotNil(t, writer)

	// 测试按日期分割的日志写入器
	writer = getLogWriter(logFile, 10, 1, 1, true)
	assert.NotNil(t, writer)
}

func TestGetDailyLogFilename(t *testing.T) {
	// 测试带扩展名的文件名
	filename := "app.log"
	dailyFilename := GetDailyLogFilename(filename)
	expectedDate := time.Now().Format("2006-01-02")
	assert.Equal(t, "app-"+expectedDate+".log", dailyFilename)

	// 测试不带扩展名的文件名
	filename = "app"
	dailyFilename = GetDailyLogFilename(filename)
	assert.Equal(t, "app-"+expectedDate, dailyFilename)

	// 测试复杂文件名
	filename = "logs/app.log"
	dailyFilename = GetDailyLogFilename(filename)
	assert.Equal(t, "logs/app-"+expectedDate+".log", dailyFilename)
}

func TestLogMethods(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	cfg := &LogConfig{
		Level:      "debug",
		Filename:   logFile,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Daily:      false,
	}

	err := Init(cfg, "prod")
	require.NoError(t, err)

	// 测试各种日志级别
	Info("info message", zap.String("key", "value"))
	Debug("debug message", zap.Int("number", 42))
	Warn("warn message", zap.Bool("flag", true))
	Error("error message", zap.Float64("float", 3.14))

	// 测试Fatal和Panic（在测试中需要特殊处理）
	// 注意：这些方法会终止程序，所以在测试中需要谨慎使用
}

func TestSync(t *testing.T) {
	// 创建临时目录用于测试
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	cfg := &LogConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Daily:      false,
	}

	err := Init(cfg, "prod")
	require.NoError(t, err)

	// 测试Sync方法
	Sync() // 应该不会panic
}

func TestLogFileCreation(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	cfg := &LogConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Daily:      false,
	}

	err := Init(cfg, "prod")
	require.NoError(t, err)

	// 写入一些日志
	Info("test message 1")
	Info("test message 2")

	// 检查日志文件是否创建
	_, err = os.Stat(logFile)
	assert.NoError(t, err, "日志文件应该被创建")
}

func TestDailyLogFileCreation(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "app.log")

	cfg := &LogConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Daily:      true,
	}

	err := Init(cfg, "prod")
	require.NoError(t, err)

	// 写入一些日志
	Info("test daily message")

	// 检查按日期分割的日志文件是否创建
	expectedDate := time.Now().Format("2006-01-02")
	dailyLogFile := filepath.Join(tempDir, "app-"+expectedDate+".log")
	_, err = os.Stat(dailyLogFile)
	assert.NoError(t, err, "按日期分割的日志文件应该被创建")
}

func TestLogLevels(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	// 测试不同日志级别
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		cfg := &LogConfig{
			Level:      level,
			Filename:   logFile,
			MaxSize:    10,
			MaxAge:     1,
			MaxBackups: 1,
			Daily:      false,
		}

		err := Init(cfg, "prod")
		require.NoError(t, err, "初始化日志级别 %s 应该成功", level)

		// 测试该级别的日志输出
		switch level {
		case "debug":
			Debug("debug message")
		case "info":
			Info("info message")
		case "warn":
			Warn("warn message")
		case "error":
			Error("error message")
		}
	}
}

func TestConcurrentLogging(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "concurrent.log")

	cfg := &LogConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    100,
		MaxAge:     1,
		MaxBackups: 1,
		Daily:      false,
	}

	err := Init(cfg, "prod")
	require.NoError(t, err)

	// 并发写入日志
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				Info("concurrent message", zap.Int("goroutine", id), zap.Int("message", j))
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 检查日志文件是否创建
	_, err = os.Stat(logFile)
	assert.NoError(t, err, "并发日志文件应该被创建")
}
