package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DB_USER       = getEnv("DB_USER", "root")
	DB_PASSWORD   = getEnv("DB_PASSWORD", "root")
	DB_NAME       = getEnv("DB_NAME", "aviation_booking")
	DB_HOST       = getEnv("DB_HOST", "127.0.0.1")
	DB_PORT       = getEnv("DB_PORT", "3306")
	AmadeusAPIKey = getEnv("AMADEUS_API_KEY", "")
)

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func Init() {
	// Загружаем переменные окружения из .env файла
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Логируем значения переменных окружения для отладки
	log.Printf("DB_USER: %s", DB_USER)
	log.Printf("DB_PASSWORD: %s", DB_PASSWORD)
	log.Printf("DB_HOST: %s", DB_HOST)
	log.Printf("DB_PORT: %s", DB_PORT)
	log.Printf("DB_NAME: %s", DB_NAME)

	log.Println("Config initialized")
}
