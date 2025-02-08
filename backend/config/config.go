package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
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

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка чтения конфигурации: %v", err)
	}

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("Ошибка разбора конфигурации: %v", err)
	}

	log.Println("Конфигурация загружена успешно")
}
