package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var config Config

type Config struct {
	TelegramToken string
	DatabaseUrl   string
	Sample        string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file")
	}
	config = Config{
		TelegramToken: getEnv("TELEGRAM_TOKEN", ""),
		DatabaseUrl:   getEnv("DATABASE_URL", ""),
		Sample:        getEnv("SAMPLE", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func GetConfig() Config {
	return config
}
