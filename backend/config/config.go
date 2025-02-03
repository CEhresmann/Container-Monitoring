package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		User     string
		Password string
		DBName   string
	}
	Queue struct {
		Type   string
		Broker string
	}
}

var Cfg Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found")
		} else {
			log.Println(err)
		}
	}
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Println(err)
	}
}
