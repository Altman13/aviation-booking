package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"errors"
	"time"
)

// CreateBooking - создание нового бронирования
func CreateBooking(flightID, userID uint, seats int, class string) (*models.Booking, error) {
	db := config.DB

	// Проверяем существование рейса
	var flight models.Flight
	if err := db.First(&flight, flightID).Error; err != nil {
		return nil, errors.New("flight not found")
	}

	// Проверяем существование пользователя
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Проверяем доступность мест
	if flight.AvailableSeats < seats {
		return nil, errors.New("not enough available seats")
	}

	// Если класс не указан, используем economy по умолчанию
	if class == "" {
		class = "economy"
	}

	// Рассчитываем общую стоимость
	totalAmount := calculateTotalCost(flight.Price, seats, class)

	// Создаем бронирование с правильными полями
	booking := &models.Booking{
		UserID:        userID,
		FlightID:      flightID,
		Seats:         seats,
		Class:         class,
		BookingDate:   time.Now(),
		TotalAmount:   totalAmount,
		Status:        "confirmed",
		PaymentStatus: "paid",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Начинаем транзакцию вручную
	tx := db.Begin()

	// Сохраняем бронирование в БД
	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Обновляем количество доступных мест
	if err := tx.Model(&flight).Update("available_seats", flight.AvailableSeats-seats).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Коммитим транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return booking, nil
}

// GetUserBookings - получение бронирований пользователя
func GetUserBookings(userID uint) ([]models.Booking, error) {
	db := config.DB

	// Проверяем существование пользователя
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var bookings []models.Booking
	if err := db.Preload("Flight").Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, err
	}

	return bookings, nil
}

// CancelBooking - отмена бронирования
func CancelBooking(bookingID uint) error {
	db := config.DB

	// Находим бронирование
	var booking models.Booking
	if err := db.First(&booking, bookingID).Error; err != nil {
		return errors.New("booking not found")
	}

	// Проверяем, не отменено ли уже бронирование
	if booking.Status == "cancelled" {
		return errors.New("booking already cancelled")
	}

	// Находим рейс для возврата мест
	var flight models.Flight
	if err := db.First(&flight, booking.FlightID).Error; err != nil {
		return errors.New("flight not found")
	}

	// Начинаем транзакцию вручную
	tx := db.Begin()

	// Отменяем бронирование
	if err := tx.Model(&booking).Update("status", "cancelled").Error; err != nil {
		tx.Rollback()
		return err
	}

	// Возвращаем места
	if err := tx.Model(&flight).Update("available_seats", flight.AvailableSeats+booking.Seats).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Коммитим транзакцию
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// GetBookingByID - получение бронирования по ID
func GetBookingByID(bookingID uint) (*models.Booking, error) {
	db := config.DB

	var booking models.Booking
	if err := db.Preload("Flight").Preload("User").First(&booking, bookingID).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

// UpdateBooking - обновление бронирования
func UpdateBooking(bookingID uint, seats int, class string) (*models.Booking, error) {
	db := config.DB

	// Находим бронирование
	var booking models.Booking
	if err := db.First(&booking, bookingID).Error; err != nil {
		return nil, errors.New("booking not found")
	}

	// Находим рейс
	var flight models.Flight
	if err := db.First(&flight, booking.FlightID).Error; err != nil {
		return nil, errors.New("flight not found")
	}

	// Рассчитываем разницу в местах
	seatDifference := seats - booking.Seats

	// Проверяем доступность мест при увеличении количества
	if seatDifference > 0 && flight.AvailableSeats < seatDifference {
		return nil, errors.New("not enough available seats")
	}

	// Начинаем транзакцию вручную
	tx := db.Begin()

	// Обновляем количество доступных мест
	if err := tx.Model(&flight).Update("available_seats", flight.AvailableSeats-seatDifference).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Обновляем данные бронирования
	booking.Seats = seats
	booking.Class = class
	booking.TotalAmount = calculateTotalCost(flight.Price, seats, class)
	booking.UpdatedAt = time.Now()

	if err := tx.Save(&booking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Коммитим транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &booking, nil
}

// calculateTotalCost - расчет общей стоимости бронирования
func calculateTotalCost(basePrice float64, seats int, class string) float64 {
	// Множители для разных классов
	classMultipliers := map[string]float64{
		"economy":  1.0,
		"comfort":  1.5,
		"business": 2.5,
		"first":    4.0,
	}

	multiplier := classMultipliers["economy"] // по умолчанию
	if m, exists := classMultipliers[class]; exists {
		multiplier = m
	}

	return basePrice * float64(seats) * multiplier
}
