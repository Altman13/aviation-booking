package services

import (
	"aviation-booking/config" // Импортируем конфигурацию с базой данных
	"aviation-booking/models" // Импортируем модели данных
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookFlight(t *testing.T) {
	// Инициализируем тестовую базу данных
	config.InitTestDB()
	defer config.CloseTestDB()

	// Создаем тестовые данные
	flight := &models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  "2023-12-01T10:00:00Z",
		ArrivalTime:    "2023-12-01T14:00:00Z",
		Price:          150.00,
		AvailableSeats: 10,
	}

	// Добавляем рейс в базу данных
	if err := config.TestDB.Create(flight).Error; err != nil {
		t.Fatalf("Failed to create flight: %v", err)
	}

	// Тестируем успешное бронирование
	err := BookFlight(flight)
	assert.NoError(t, err)
	assert.Equal(t, 9, flight.AvailableSeats) // Проверяем, что количество мест уменьшилось на 1

	// Тестируем бронирование без доступных мест
	flight.AvailableSeats = 0
	err = BookFlight(flight)
	assert.EqualError(t, err, "No available seats") // Проверяем, что ошибка "No available seats"
}
