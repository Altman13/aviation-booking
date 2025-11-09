package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_USER       string
	DB_PASSWORD   string
	DB_NAME       string
	DB_HOST       string
	DB_PORT       string
	AmadeusAPIKey string
)

func Init() {
	// Загружаем переменные окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using environment variables or defaults")
	}

	// Инициализируем переменные
	DB_USER = getEnv("DB_USER", "root")
	DB_PASSWORD = getEnv("DB_PASSWORD", "")
	DB_NAME = getEnv("DB_NAME", "aviation_booking")
	DB_HOST = getEnv("DB_HOST", "127.0.0.1")
	DB_PORT = getEnv("DB_PORT", "3306")
	AmadeusAPIKey = getEnv("AMADEUS_API_KEY", "")

	log.Printf("DB_USER: %s", DB_USER)
	log.Printf("DB_HOST: %s", DB_HOST)
	log.Printf("DB_PORT: %s", DB_PORT)
	log.Printf("DB_NAME: %s", DB_NAME)

	log.Println("Config initialized")
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetEnv для использования в других пакетах
func GetEnv(key string, defaultValue string) string {
	return getEnv(key, defaultValue)
}
