package config

import (
	"github.com/spf13/viper"
	"log"
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
		log.Printf("Ошибка чтения конфигурации: %v (используем только ENV)", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("Ошибка разбора конфигурации: %v", err)
	}

	if broker := os.Getenv("QUEUE_BROKER"); broker != "" {
		Cfg.Queue.Broker = broker
	}

	if topic := os.Getenv("PING_TOPIC"); topic != "" {
		Cfg.Queue.Topic = topic
	}

	log.Println("Конфигурация загружена успешно")
}
