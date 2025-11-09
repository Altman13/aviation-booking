package controllers

import (
	"aviation-booking/models"
	"aviation-booking/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SearchFlightsRequest - структура для запроса поиска рейсов
type SearchFlightsRequest struct {
	DepartureCity string `form:"departure_city"`
	ArrivalCity   string `form:"arrival_city"`
}

// SearchFlights - обработчик для поиска рейсов
func SearchFlights(c *gin.Context) {
	var req SearchFlightsRequest

	// Привязываем параметры запроса к структуре
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Получаем рейсы с помощью сервисного слоя
	flights, err := services.SearchFlights(req.DepartureCity, req.ArrivalCity)
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
				"departure_city": req.DepartureCity,
				"arrival_city":   req.ArrivalCity,
			},
			"total": len(flights),
		},
	})
}

// GetFlightDetails - получение деталей рейса по ID
func GetFlightDetails(c *gin.Context) {
	flightIDStr := c.Param("id")

	if flightIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Flight ID is required",
		})
		return
	}

	// Конвертируем string в uint
	flightID, err := strconv.ParseUint(flightIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid flight ID format",
		})
		return
	}

	// Получаем детали рейса
	flight, err := services.GetFlightByID(uint(flightID))
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

// GetAllFlights - получение всех рейсов
func GetAllFlights(c *gin.Context) {
	// Получаем все рейсы
	flights, err := services.GetAllFlights()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"flights": flights,
			"total":   len(flights),
		},
	})
}

// CreateFlightRequest - структура для запроса создания рейса
type CreateFlightRequest struct {
	DepartureCity  string  `json:"departure_city" binding:"required"`
	ArrivalCity    string  `json:"arrival_city" binding:"required"`
	DepartureTime  string  `json:"departure_time" binding:"required"`
	ArrivalTime    string  `json:"arrival_time" binding:"required"`
	Price          float64 `json:"price" binding:"required,min=0"`
	AvailableSeats int     `json:"available_seats" binding:"required,min=1"`
}

// CreateFlight - создание нового рейса
func CreateFlight(c *gin.Context) {
	var req CreateFlightRequest

	// Привязываем JSON из тела запроса к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Парсим время
	departureTime, err := parseTime(req.DepartureTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid departure time format",
		})
		return
	}

	arrivalTime, err := parseTime(req.ArrivalTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid arrival time format",
		})
		return
	}

	// Проверяем, что время прибытия после времени вылета
	if arrivalTime.Before(departureTime) || arrivalTime.Equal(departureTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Arrival time must be after departure time",
		})
		return
	}

	// Создаем рейс (используем models.Flight вместо services.Flight)
	flight := &models.Flight{
		DepartureCity:  req.DepartureCity,
		ArrivalCity:    req.ArrivalCity,
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          req.Price,
		AvailableSeats: req.AvailableSeats,
	}

	err = services.CreateFlight(flight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Flight created successfully",
		"data":    flight,
	})
}

// UpdateFlightSeatsRequest - структура для запроса обновления мест
type UpdateFlightSeatsRequest struct {
	AvailableSeats int `json:"available_seats" binding:"required,min=0"`
}

// UpdateFlightSeats - обновление количества доступных мест
func UpdateFlightSeats(c *gin.Context) {
	flightIDStr := c.Param("id")

	if flightIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Flight ID is required",
		})
		return
	}

	// Конвертируем string в uint
	flightID, err := strconv.ParseUint(flightIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid flight ID format",
		})
		return
	}

	var req UpdateFlightSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Обновляем количество мест
	err = services.UpdateFlightSeats(uint(flightID), req.AvailableSeats)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "record not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Flight seats updated successfully",
	})
}

// Вспомогательная функция для парсинга времени
func parseTime(timeStr string) (time.Time, error) {
	// Пробуем разные форматы времени
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, &time.ParseError{}
}