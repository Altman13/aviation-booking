package controllers

import (
	"aviation-booking/models"
	"aviation-booking/services" // Импортируем сервис
	"net/http"

	"github.com/gin-gonic/gin"
)

// BookFlight - обработчик для бронирования рейса
func BookFlight(c *gin.Context) {
	// Получаем данные от клиента
	var flight models.Flight
	if err := c.BindJSON(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вызываем сервис для бронирования
	if err := services.BookFlight(&flight); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ответ пользователю
	c.JSON(http.StatusOK, gin.H{
		"message": "Booking successful",
		"flight":  flight,
	})
}
