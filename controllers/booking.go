package controllers

import (
	"aviation-booking/models"
	"aviation-booking/services" // Импортируем сервис для бронирования
	"net/http"

	"github.com/gin-gonic/gin"
)

// BookFlight - обработчик для бронирования рейса
func BookFlight(c *gin.Context) {
	// Получаем данные от клиента
	var flight models.Flight
	if err := c.BindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid flight data", "details": err.Error()})
		return
	}

	// Проверка, что количество мест не меньше 1 (на всякий случай)
	if flight.AvailableSeats <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No available seats"})
		return
	}

	// Вызываем сервис для бронирования
	if err := services.BookFlight(&flight); err != nil {
		// В случае ошибки в сервисе бронирования, возвращаем статус 500
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book flight", "details": err.Error()})
		return
	}

	// Ответ пользователю с успешным статусом
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking successful",
		"flight":  flight,
	})
}
