package config

import (
	"aviation-booking/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB     // Основная база данных
var TestDB *gorm.DB // Тестовая база данных

// InitDB инициализирует подключение к основной базе данных MySQL
func InitDB() error {
	var err error

	// Получаем значения из config/config.go
	dbUser := GetEnv("DB_USER", "root")
	dbPassword := GetEnv("DB_PASSWORD", "")
	dbHost := GetEnv("DB_HOST", "127.0.0.1")
	dbPort := GetEnv("DB_PORT", "3306")
	dbName := GetEnv("DB_NAME", "aviation_booking")

	// Строка подключения к основной базе данных
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Printf("Connecting to MAIN database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("could not connect to main database: %v", err)
	}

	// Автогенерация таблиц
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Airport{},
		&models.Airline{},
		&models.Flight{},
		&models.Booking{},
		&models.Ticket{},
	); err != nil {
		return fmt.Errorf("could not migrate main database: %v", err)
	}

	log.Println("✅ Main database connected and migrated successfully")
	return nil
}

// InitTestDB инициализирует тестовую базу данных
func InitTestDB() error {
	var err error

	// Получаем значения из config/config.go
	dbUser := GetEnv("DB_USER", "root")
	dbPassword := GetEnv("DB_PASSWORD", "root")
	dbHost := GetEnv("DB_HOST", "127.0.0.1")
	dbPort := GetEnv("DB_PORT", "3306")
	dbName := GetEnv("DB_NAME", "aviation_booking")
	testDBName := dbName + "_test"

	// Строка подключения к тестовой базе данных
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, testDBName)

	log.Printf("Connecting to TEST database: %s@%s:%s/%s", dbUser, dbHost, dbPort, testDBName)

	TestDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("could not connect to test database: %v", err)
	}

	// Автогенерация таблиц для тестов
	if err := TestDB.AutoMigrate(
		&models.User{},
		&models.Airport{},
		&models.Airline{},
		&models.Flight{},
		&models.Booking{},
		&models.Ticket{},
	); err != nil {
		return fmt.Errorf("could not migrate test database: %v", err)
	}

	log.Println("✅ Test database connected and migrated successfully")
	return nil
}

// GetDB возвращает экземпляр основной базы данных
func GetDB() *gorm.DB {
	return DB
}

// GetTestDB возвращает экземпляр тестовой базы данных
func GetTestDB() *gorm.DB {
	return TestDB
}

// CloseDB закрывает соединение с основной базой данных
func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Failed to get raw SQL database object: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Could not close main database connection: %v", err)
		} else {
			log.Println("Main database connection closed")
		}
	}
}

// CloseTestDB закрывает соединение с тестовой базой данных
func CloseTestDB() {
	if TestDB != nil {
		sqlDB, err := TestDB.DB()
		if err != nil {
			log.Printf("Failed to get raw SQL database object: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			log.Printf("Could not close test database connection: %v", err)
		} else {
			log.Println("Test database connection closed")
		}
	}
}
