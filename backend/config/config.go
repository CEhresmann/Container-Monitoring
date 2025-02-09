package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Queue struct {
		Type   string `yaml:"type"`
		Broker string `yaml:"broker"`
		Topic  string `yaml:"topic"`
	} `yaml:"queue"`
}

var Cfg Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		zap.L().Info("Ошибка чтения конфигурации: (используем только ENV)")
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		zap.L().Error("Ошибка разбора конфигурации")
		return
	}

	if broker := os.Getenv("QUEUE_BROKER"); broker != "" {
		Cfg.Queue.Broker = broker
	}

	if topic := os.Getenv("PING_TOPIC"); topic != "" {
		Cfg.Queue.Topic = topic
	}

	zap.L().Info("Конфигурация загружена успешно")
}
