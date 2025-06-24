package config

import (
	"aviation-booking/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB // Основная база данных

// InitDB инициализирует подключение к основной базе данных MySQL
func InitDB() {
	var err error

	// Строка подключения к MySQL (формат: user:password@tcp(host:port)/dbname)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Автогенерация таблиц
	if err := DB.AutoMigrate(&models.Flight{}); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
}

// CloseDB закрывает соединение с основной базой данных
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get raw SQL database object: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Could not close database connection: %v", err)
	}
}

var TestDB *gorm.DB // Переменная для тестовой базы данных
func InitTestDB() {
	var err error

	// Строка подключения к тестовой базе данных MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s_test?charset=utf8&parseTime=True&loc=Local",
		DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	TestDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to test database: %v", err)
	}

	// Автогенерация таблиц для тестов
	if err := TestDB.AutoMigrate(&models.Flight{}); err != nil {
		log.Fatalf("Could not migrate test database: %v", err)
	}
}

// CloseTestDB закрывает соединение с тестовой базой данных
func CloseTestDB() {
	sqlDB, err := TestDB.DB()
	if err != nil {
		log.Fatalf("Failed to get raw SQL database object: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatalf("Failed to close connection: %v", err)
	}
}
