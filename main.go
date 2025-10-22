package main

import (
	"aviation-booking/config"
	"aviation-booking/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация базы данных
	config.InitDB()
	defer config.CloseDB()

	// Создание маршрутов
	r := gin.Default()

	// Health check endpoints
	r.GET("/health", controllers.HealthCheck)
	r.GET("/health/simple", controllers.SimpleHealth)
	r.GET("/ready", controllers.ReadyCheck)
	r.GET("/live", controllers.LiveCheck)
	r.GET("/stats/db", controllers.DatabaseStats)

	// Flight search endpoints - используем функции из flight.go
	r.GET("/flights/search", controllers.SearchFlights) // Из flight.go
	r.GET("/flights/:id", controllers.GetFlightDetails) // Из flight.go

	// Booking endpoints
	r.POST("/book", controllers.BookFlight)
	r.GET("/bookings/:user_id", controllers.GetUserBookings)
	r.DELETE("/bookings/:booking_id", controllers.CancelBooking)

	// User endpoints
	r.POST("/auth/register", controllers.RegisterUser)
	r.POST("/auth/login", controllers.LoginUser)

	// Запуск сервера
	r.Run(":8080")
}
