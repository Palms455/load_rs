package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
)

// Структура с конфигурацией бд, воркеров
var MainConfig *Config
func init() {
	MainConfig = NewConfig()
	configPath := filepath.Join("configs", "rs.toml")
	_, err := toml.DecodeFile(configPath, MainConfig)
	if err != nil {
		log.Println("Не найден .toml файл:", err)
	}
}


type Config struct {
	//URL DB
	DbAddr string `toml:"db"`
	//Logger Level
	LoggerLevel string `toml:"logger_level"`
}

func NewConfig() *Config {
	return &Config{
		DbAddr:    "postgres://username:password@localhost:5432/database_name",
		LoggerLevel: "debug",
	}
}

