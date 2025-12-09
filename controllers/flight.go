// controllers/flight_controller.go
package controllers

import (
	"aviation-booking/models"
	"aviation-booking/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FlightController struct {
	flightService *services.FlightService
}

// NewFlightController создает новый контроллер рейсов
func NewFlightController() *FlightController {
	return &FlightController{
		flightService: services.NewFlightService(),
	}
}

// FlightSearchRequest - структура для запроса поиска рейсов
type FlightSearchRequest struct {
	DepartureCity     string  `form:"departure_city"`
	ArrivalCity       string  `form:"arrival_city"`
	DepartureTimeFrom string  `form:"departure_time_from"`
	DepartureTimeTo   string  `form:"departure_time_to"`
	MinPrice          float64 `form:"min_price"`
	MaxPrice          float64 `form:"max_price"`
	MinAvailableSeats int     `form:"min_available_seats"`
	Airline           string  `form:"airline"`
	SortBy            string  `form:"sort_by"`
	Limit             int     `form:"limit,default=50"`
	Page              int     `form:"page,default=1"`
}

// SearchFlights - обработчик для поиска рейсов
func (fc *FlightController) SearchFlights(c *gin.Context) {
	var req FlightSearchRequest

	// Привязываем параметры запроса к структуре
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Конвертируем параметры в структуру сервиса
	params := services.FlightSearchParams{
		DepartureCity:     req.DepartureCity,
		ArrivalCity:       req.ArrivalCity,
		MinPrice:          req.MinPrice,
		MaxPrice:          req.MaxPrice,
		MinAvailableSeats: req.MinAvailableSeats,
		Airline:           req.Airline,
		SortBy:            req.SortBy,
		Limit:             req.Limit,
		Offset:            (req.Page - 1) * req.Limit,
	}

	// Парсим время
	if req.DepartureTimeFrom != "" {
		if t, err := parseTime(req.DepartureTimeFrom); err == nil {
			params.DepartureTimeFrom = &t
		}
	}

	if req.DepartureTimeTo != "" {
		if t, err := parseTime(req.DepartureTimeTo); err == nil {
			params.DepartureTimeTo = &t
		}
	}

	// Получаем рейсы с помощью сервисного слоя
	flights, err := fc.flightService.SearchFlights(params)
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
			"total":   len(flights),
			"page":    req.Page,
			"limit":   req.Limit,
		},
	})
}

// GetFlightByID - получение деталей рейса по ID
func (fc *FlightController) GetFlightByID(c *gin.Context) {
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
	flight, err := fc.flightService.GetFlightByID(uint(flightID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   flight,
	})
}

// GetAllFlights - получение всех рейсов
func (fc *FlightController) GetAllFlights(c *gin.Context) {
	// Получаем все рейсы
	flights, err := fc.flightService.GetAllFlights()
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
	FlightNumber         string  `json:"flight_number" binding:"required"`
	Airline              string  `json:"airline" binding:"required"`
	DepartureCity        string  `json:"departure_city" binding:"required"`
	DepartureAirportCode string  `json:"departure_airport_code"`
	ArrivalCity          string  `json:"arrival_city" binding:"required"`
	ArrivalAirportCode   string  `json:"arrival_airport_code"`
	DepartureTime        string  `json:"departure_time" binding:"required"`
	ArrivalTime          string  `json:"arrival_time" binding:"required"`
	Price                float64 `json:"price" binding:"required,min=0"`
	AvailableSeats       int     `json:"available_seats" binding:"required,min=0"`
}

// CreateFlight - создание нового рейса
func (fc *FlightController) CreateFlight(c *gin.Context) {
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

	// Создаем рейс
	flight := &models.Flight{
		FlightNumber:         req.FlightNumber,
		Airline:              req.Airline,
		DepartureCity:        req.DepartureCity,
		DepartureAirportCode: req.DepartureAirportCode,
		ArrivalCity:          req.ArrivalCity,
		ArrivalAirportCode:   req.ArrivalAirportCode,
		DepartureTime:        departureTime,
		ArrivalTime:          arrivalTime,
		Price:                req.Price,
		AvailableSeats:       req.AvailableSeats,
	}

	err = fc.flightService.CreateFlight(flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

// UpdateFlightRequest - структура для полного обновления рейса
type UpdateFlightRequest struct {
	FlightNumber         string  `json:"flight_number"`
	Airline              string  `json:"airline"`
	DepartureCity        string  `json:"departure_city"`
	DepartureAirportCode string  `json:"departure_airport_code"`
	ArrivalCity          string  `json:"arrival_city"`
	ArrivalAirportCode   string  `json:"arrival_airport_code"`
	DepartureTime        string  `json:"departure_time"`
	ArrivalTime          string  `json:"arrival_time"`
	Price                float64 `json:"price" min="0"`
	AvailableSeats       int     `json:"available_seats" min="0"`
}

// UpdateFlight - полное обновление рейса
func (fc *FlightController) UpdateFlight(c *gin.Context) {
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

	var req UpdateFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Парсим время
	var departureTime, arrivalTime time.Time
	if req.DepartureTime != "" {
		departureTime, err = parseTime(req.DepartureTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid departure time format",
			})
			return
		}
	}

	if req.ArrivalTime != "" {
		arrivalTime, err = parseTime(req.ArrivalTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid arrival time format",
			})
			return
		}
	}

	// Создаем объект рейса для обновления
	flightData := &models.Flight{
		FlightNumber:         req.FlightNumber,
		Airline:              req.Airline,
		DepartureCity:        req.DepartureCity,
		DepartureAirportCode: req.DepartureAirportCode,
		ArrivalCity:          req.ArrivalCity,
		ArrivalAirportCode:   req.ArrivalAirportCode,
		Price:                req.Price,
		AvailableSeats:       req.AvailableSeats,
	}

	// Устанавливаем время только если оно передано
	if !departureTime.IsZero() {
		flightData.DepartureTime = departureTime
	}
	if !arrivalTime.IsZero() {
		flightData.ArrivalTime = arrivalTime
	}

	// Обновляем рейс
	err = fc.flightService.UpdateFlight(uint(flightID), flightData)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "рейс не найден" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Flight updated successfully",
	})
}

// UpdateFlightSeatsRequest - структура для запроса обновления мест
type UpdateFlightSeatsRequest struct {
	AvailableSeats int `json:"available_seats" binding:"required,min=0"`
}

// UpdateFlightSeats - обновление только количества доступных мест
func (fc *FlightController) UpdateFlightSeats(c *gin.Context) {
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
	err = fc.flightService.UpdateFlightSeats(uint(flightID), req.AvailableSeats)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "рейс не найден" {
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

// DeleteFlight - удаление рейса
func (fc *FlightController) DeleteFlight(c *gin.Context) {
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

	// Удаляем рейс
	err = fc.flightService.DeleteFlight(uint(flightID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "рейс не найден" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Flight deleted successfully",
	})
}

// CheckAvailabilityRequest - структура для проверки доступности мест
type CheckAvailabilityRequest struct {
	RequiredSeats int `json:"required_seats" binding:"required,min=1"`
}

// CheckFlightAvailability - проверка доступности мест на рейсе
func (fc *FlightController) CheckFlightAvailability(c *gin.Context) {
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

	var req CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Проверяем доступность мест
	available, err := fc.flightService.CheckFlightAvailability(uint(flightID), req.RequiredSeats)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"available": available,
			"message":   map[bool]string{true: "Seats are available", false: "Not enough seats"}[available],
		},
	})
}

// ReserveSeatsRequest - структура для бронирования мест
type ReserveSeatsRequest struct {
	Seats int `json:"seats" binding:"required,min=1"`
}

// ReserveSeats - бронирование мест на рейсе
func (fc *FlightController) ReserveSeats(c *gin.Context) {
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

	var req ReserveSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Резервируем места
	err = fc.flightService.ReserveSeats(uint(flightID), req.Seats)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "рейс не найден" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Seats reserved successfully",
	})
}

// ReleaseSeatsRequest - структура для освобождения мест
type ReleaseSeatsRequest struct {
	Seats int `json:"seats" binding:"required,min=1"`
}

// ReleaseSeats - освобождение мест (отмена бронирования)
func (fc *FlightController) ReleaseSeats(c *gin.Context) {
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

	var req ReleaseSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Освобождаем места
	err = fc.flightService.ReleaseSeats(uint(flightID), req.Seats)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "рейс не найден" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Seats released successfully",
	})
}

// GetUpcomingFlights - получение предстоящих рейсов
func (fc *FlightController) GetUpcomingFlights(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Получаем предстоящие рейсы
	flights, err := fc.flightService.GetUpcomingFlights(limit)
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
			"limit":   limit,
		},
	})
}

// GetFlightsByAirline - получение рейсов по авиакомпании
func (fc *FlightController) GetFlightsByAirline(c *gin.Context) {
	airline := c.Param("airline")

	if airline == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Airline name is required",
		})
		return
	}

	// Получаем рейсы авиакомпании
	flights, err := fc.flightService.GetFlightsByAirline(airline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"airline": airline,
			"flights": flights,
			"total":   len(flights),
		},
	})
}

// parseTime - вспомогательная функция для парсинга времени
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
