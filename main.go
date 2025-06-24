package main

import (
	"aviation-booking/config" // Путь к конфигурации базы данных
	"aviation-booking/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	config.InitDB()
	defer config.CloseDB() // Закрытие подключения при завершении работы

	// Создание маршрутов
	r := gin.Default()
	r.POST("/book", controllers.BookFlight)

	// Запуск сервера
	r.Run(":8080")
}
