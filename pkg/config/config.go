package config

import (
	"log"
	"os"
	"probe/pkg/logger"
	"probe/pkg/utils"
)

// config/config.go
type Config struct {
	MachineID        int64  `env:"MACHINE_ID"`
	DBDriver         string `env:"DB_DRIVER"`
	DSN              string `env:"DSN"`
	Log              logger.LogConfig
	Addr             string `env:"ADDR"`
	Mode             string `env:"MODE"`
	DocsPrefix       string `env:"DOCS_PREFIX"`
	APIPrefix        string `env:"API_PREFIX"`
	AdminPrefix      string `env:"ADMIN_PREFIX"`
	AuthPrefix       string `env:"AUTH_PREFIX"`
	SessionSecret    string `env:"SESSION_SECRET"`
	SecretExpireDays string `env:"SESSION_EXPIRE_DAYS"`
}

var GlobalConfig *Config

func Load() error {
	// 1. 根据环境加载 .env 文件
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "test" // 默认使用开发环境
	}
	err := utils.LoadEnv(env)
	if err != nil {
		log.Printf("Failed to load .env file: %v", err)
	}

	// 2. 加载全局配置
	GlobalConfig = &Config{
		MachineID:        utils.GetIntEnv("MACHINE_ID"),
		DBDriver:         utils.GetEnv("DB_DRIVER"),
		DSN:              utils.GetEnv("DSN"),
		Addr:             utils.GetEnv("ADDR"),
		Mode:             utils.GetEnv("MODE"),
		DocsPrefix:       utils.GetEnv("DOCS_PREFIX"),
		APIPrefix:        utils.GetEnv("API_PREFIX"),
		AdminPrefix:      utils.GetEnv("ADMIN_PREFIX"),
		AuthPrefix:       utils.GetEnv("AUTH_PREFIX"),
		SecretExpireDays: utils.GetEnv("SESSION_EXPIRE_DAYS"),
		SessionSecret:    utils.GetEnv("SESSION_SECRET"),
		Log: logger.LogConfig{
			Level:      utils.GetEnv("LOG_LEVEL"),
			Filename:   utils.GetEnv("LOG_FILENAME"),
			MaxSize:    int(utils.GetIntEnv("LOG_MAX_SIZE")),
			MaxAge:     int(utils.GetIntEnv("LOG_MAX_AGE")),
			MaxBackups: int(utils.GetIntEnv("LOG_MAX_BACKUPS")),
			Daily:      utils.GetBoolEnv("LOG_DAILY"),
		},
	}
	return nil
}
