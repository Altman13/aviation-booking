package controllers

import (
	"aviation-booking/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SearchFlights - обработчик для поиска рейсов
func SearchFlights(c *gin.Context) {
	// Параметры поиска
	origin := c.Query("origin")
	destination := c.Query("destination")
	date := c.Query("date")
	adults := c.DefaultQuery("adults", "1")
	class := c.DefaultQuery("class", "economy")

	// Валидация входных данных
	if origin == "" || destination == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Parameters 'origin', 'destination', and 'date' are required",
		})
		return
	}

	// Получаем рейсы с помощью сервисного слоя
	flights, err := services.SearchFlights(origin, destination, date, adults, class)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Отправляем результат
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"flights": flights,
			"search_params": gin.H{
				"origin":      origin,
				"destination": destination,
				"date":        date,
				"adults":      adults,
				"class":       class,
			},
		},
	})
}

// GetFlightDetails - получение деталей рейса по ID
func GetFlightDetails(c *gin.Context) {
	flightID := c.Param("id")

	if flightID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Flight ID is required",
		})
		return
	}

	// Получаем детали рейса
	flight, err := services.GetFlightByID(flightID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Flight not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   flight,
	})
}
