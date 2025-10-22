package controllers

import (
	"aviation-booking/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BookFlightRequest - структура для запроса бронирования
type BookFlightRequest struct {
	FlightID string `json:"flight_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Flight booked successfully",
		"data":    booking,
	})
}

// GetUserBookings - получение всех бронирований пользователя
func GetUserBookings(c *gin.Context) {
	userID := c.Param("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	bookings, err := services.GetUserBookings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
	bookingID := c.Param("booking_id")

	if bookingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Booking ID is required",
		})
		return
	}

	err := services.CancelBooking(bookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Booking cancelled successfully",
	})
}
