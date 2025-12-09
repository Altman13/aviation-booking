// controllers/booking_controller.go
package controllers

import (
	"aviation-booking/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingController struct {
	bookingService *services.BookingService
}

// NewBookingController создает новый контроллер бронирований
func NewBookingController() *BookingController {
	return &BookingController{
		bookingService: services.NewBookingService(),
	}
}

// BookFlightRequest - структура для запроса бронирования
type BookFlightRequest struct {
	FlightID uint   `json:"flight_id" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
	Seats    int    `json:"seats" binding:"required,min=1"`
	Class    string `json:"class,omitempty"`
}

// BookFlight - обработчик бронирования рейса
func (bc *BookingController) BookFlight(c *gin.Context) {
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
	booking, err := bc.bookingService.BookFlight(req.FlightID, req.UserID, req.Seats, req.Class)
	if err != nil {
		// Определяем тип ошибки для возврата соответствующего статуса
		statusCode := http.StatusInternalServerError
		if err.Error() == "рейс не найден" || err.Error() == "пользователь не найден" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "недостаточно мест на рейсе" {
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
func (bc *BookingController) GetUserBookings(c *gin.Context) {
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

	bookings, err := bc.bookingService.GetUserBookings(uint(userID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "пользователь не найден" {
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

// GetBookingByID - получение информации о конкретном бронировании
func (bc *BookingController) GetBookingByID(c *gin.Context) {
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

	booking, err := bc.bookingService.GetBookingByID(uint(bookingID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "бронирование не найдено" {
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

// GetBookingByNumber - получение бронирования по номеру
func (bc *BookingController) GetBookingByNumber(c *gin.Context) {
	bookingNumber := c.Param("booking_number")

	if bookingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Booking number is required",
		})
		return
	}

	booking, err := bc.bookingService.GetBookingByNumber(bookingNumber)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "бронирование не найдено" {
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

// CancelBooking - отмена бронирования
func (bc *BookingController) CancelBooking(c *gin.Context) {
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

	err = bc.bookingService.CancelBooking(uint(bookingID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "бронирование не найдено" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "бронирование уже отменено" {
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

// UpdateBookingRequest - структура для запроса обновления бронирования
type UpdateBookingRequest struct {
	Seats int    `json:"seats" binding:"required,min=1"`
	Class string `json:"class" binding:"required"`
}

// UpdateBooking - обновление бронирования
func (bc *BookingController) UpdateBooking(c *gin.Context) {
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

	booking, err := bc.bookingService.UpdateBooking(uint(bookingID), req.Seats, req.Class)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "бронирование не найдено" || err.Error() == "рейс не найден" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "недостаточно мест на рейсе" {
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

// ConfirmPayment - подтверждение оплаты
func (bc *BookingController) ConfirmPayment(c *gin.Context) {
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

	err = bc.bookingService.ConfirmPayment(uint(bookingID))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "бронирование не найдено" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "только подтвержденные бронирования могут быть оплачены" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Payment confirmed successfully",
	})
}

// SearchBookings - поиск бронирований
func (bc *BookingController) SearchBookings(c *gin.Context) {
	// Парсим параметры запроса
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	flightID, _ := strconv.ParseUint(c.Query("flight_id"), 10, 32)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	// Парсим даты (если есть)
	var dateFrom, dateTo time.Time
	if df := c.Query("date_from"); df != "" {
		dateFrom, _ = time.Parse("2006-01-02", df)
	}
	if dt := c.Query("date_to"); dt != "" {
		dateTo, _ = time.Parse("2006-01-02", dt)
	}

	params := services.BookingSearchParams{
		UserID:        uint(userID),
		FlightID:      uint(flightID),
		Status:        c.Query("status"),
		PaymentStatus: c.Query("payment_status"),
		Class:         c.Query("class"),
		DateFrom:      dateFrom,
		DateTo:        dateTo,
		SortBy:        c.Query("sort_by"),
		Limit:         limit,
		Offset:        (page - 1) * limit,
	}

	// Ищем бронирования
	bookings, err := bc.bookingService.SearchBookings(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"bookings": bookings,
			"total":    len(bookings),
			"page":     page,
			"limit":    limit,
		},
	})
}
