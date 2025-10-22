package services

import "time"

// Booking представляет информацию о бронировании
type Booking struct {
	ID        string  `json:"id"`
	FlightID  string  `json:"flight_id"`
	UserID    string  `json:"user_id"`
	Seats     int     `json:"seats"`
	Class     string  `json:"class"`
	TotalCost float64 `json:"total_cost"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
}

// CreateBooking - создание нового бронирования
func CreateBooking(flightID, userID string, seats int, class string) (*Booking, error) {
	// Если класс не указан, используем economy по умолчанию
	if class == "" {
		class = "economy"
	}

	// Временная заглушка - в реальном приложении здесь будет логика работы с БД
	booking := &Booking{
		ID:        "booking-" + generateID(),
		FlightID:  flightID,
		UserID:    userID,
		Seats:     seats,
		Class:     class,
		TotalCost: calculateTotalCost(flightID, seats, class),
		Status:    "confirmed",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	return booking, nil
}

// GetUserBookings - получение бронирований пользователя
func GetUserBookings(userID string) ([]Booking, error) {
	// Временная заглушка - в реальном приложении здесь будет запрос к БД
	bookings := []Booking{
		{
			ID:        "booking-001",
			FlightID:  "SU-1234",
			UserID:    userID,
			Seats:     2,
			Class:     "economy",
			TotalCost: 10000.00,
			Status:    "confirmed",
			CreatedAt: "2024-01-15T12:00:00Z",
		},
		{
			ID:        "booking-002",
			FlightID:  "S7-5678",
			UserID:    userID,
			Seats:     1,
			Class:     "business",
			TotalCost: 25000.00,
			Status:    "confirmed",
			CreatedAt: "2024-01-16T14:30:00Z",
		},
	}

	return bookings, nil
}

// CancelBooking - отмена бронирования
func CancelBooking(bookingID string) error {
	// Временная заглушка - в реальном приложении здесь будет логика отмены в БД
	// Проверяем существование бронирования и т.д.

	// Симулируем успешную отмену
	return nil
}

// Вспомогательные функции
func generateID() string {
	return time.Now().Format("20060102150405")
}

func calculateTotalCost(flightID string, seats int, class string) float64 {
	// Временная логика расчета стоимости
	basePrice := 5000.0
	if class == "business" {
		basePrice = 15000.0
	} else if class == "first" {
		basePrice = 30000.0
	}
	return basePrice * float64(seats)
}
