package config

import (
	"log"
	"os"
)

var (
	DB_USER     = getEnv("DB_USER", "root")
	DB_PASSWORD = getEnv("DB_PASSWORD", "")
	DB_NAME     = getEnv("DB_NAME", "aviation_booking")
	DB_HOST     = getEnv("DB_HOST", "localhost")
	DB_PORT     = getEnv("DB_PORT", "3306")
)

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func Init() {
	// Здесь можно добавить логику для подключения к базе данных
	log.Println("Config initialized")
}
