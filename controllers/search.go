package controllers

import (
	"aviation-booking/services" // Создадим сервис для логики поиска
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchFlights - обработчик для поиска рейсов
func SearchFlights(c *gin.Context) {
	// Параметры поиска
	departureCity := c.DefaultQuery("departure", "")
	arrivalCity := c.DefaultQuery("arrival", "")

	// Валидация входных данных
	if departureCity == "" || arrivalCity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both departure and arrival cities are required"})
		return
	}

	// Получаем рейсы с помощью сервисного слоя
	flights, err := services.SearchFlights(departureCity, arrivalCity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправляем результат
	c.JSON(http.StatusOK, flights)
}
