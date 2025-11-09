// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Проверяем текущую директорию
// 	dir, _ := os.Getwd()
// 	fmt.Printf("Current directory: %s\n", dir)

// 	// Пробуем загрузить .env
// 	err := godotenv.Load()
// 	if err != nil {
// 		fmt.Printf("Error loading .env: %v\n", err)

// 		// Проверяем существует ли файл
// 		if _, err := os.Stat(".env"); os.IsNotExist(err) {
// 			fmt.Println(".env file does not exist")
// 		} else {
// 			fmt.Println(".env file exists but couldn't be loaded")
// 		}
// 	} else {
// 		fmt.Println(".env loaded successfully")
// 		fmt.Printf("DB_USER: %s\n", os.Getenv("DB_USER"))
// 	}
// }
