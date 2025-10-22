package controllers

import (
	"aviation-booking/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookFlightRequest - структура для запроса бронирования
type BookFlightRequest struct {
	FlightID uint   `json:"flight_id" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
	Seats    int    `json:"seats" binding:"required,min=1"`
	Class    string `json:"class,omitempty"`
}

// BookFlight - обработчик бронирования рейса
func BookFlight(c *gin.Context) {
	var req BookFlightRequest

	// Привязываем JSON из тела запроса к структуре
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Создаем бронирование через сервис
	booking, err := services.CreateBooking(req.FlightID, req.UserID, req.Seats, req.Class)
	if err != nil {
		// Определяем тип ошибки для возврата соответствующего статуса
		statusCode := http.StatusInternalServerError
		if err.Error() == "flight not found" || err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "not enough available seats" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Flight booked successfully",
		"data":    booking,
	})
}

// GetUserBookings - получение всех бронирований пользователя
func GetUserBookings(c *gin.Context) {
	userIDStr := c.Param("user_id")

	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Конвертируем string в uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	bookings, err := services.GetUserBookings(uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user_id":  userID,
			"bookings": bookings,
		},
	})
}

// CancelBooking - отмена бронирования
func CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")

	if bookingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Booking ID is required",
		})
		return
	}

	// Конвертируем string в uint
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid booking ID format",
		})
		return
	}

	err = services.CancelBooking(uint(bookingID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "booking not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "booking already cancelled" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Booking cancelled successfully",
	})
}

// GetBooking - получение информации о конкретном бронировании
func GetBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")

	if bookingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Booking ID is required",
		})
		return
	}

	// Конвертируем string в uint
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid booking ID format",
		})
		return
	}

	booking, err := services.GetBookingByID(uint(bookingID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "booking not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   booking,
	})
}

// UpdateBookingRequest - структура для запроса обновления бронирования
type UpdateBookingRequest struct {
	Seats int    `json:"seats" binding:"required,min=1"`
	Class string `json:"class" binding:"required"`
}

// UpdateBooking - обновление бронирования
func UpdateBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")

	if bookingIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Booking ID is required",
		})
		return
	}

	// Конвертируем string в uint
	bookingID, err := strconv.ParseUint(bookingIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid booking ID format",
		})
		return
	}

	var req UpdateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	booking, err := services.UpdateBooking(uint(bookingID), req.Seats, req.Class)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "booking not found" || err.Error() == "flight not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "not enough available seats" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Booking updated successfully",
		"data":    booking,
	})
}
