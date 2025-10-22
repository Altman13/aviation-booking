package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchFlights(t *testing.T) {
	// Инициализируем тестовую базу данных
	config.InitTestDB()
	defer config.CloseTestDB()

	// Создаем тестовые данные
	flight1 := &models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  "2023-12-01T10:00:00Z",
		ArrivalTime:    "2023-12-01T14:00:00Z",
		Price:          150.00,
		AvailableSeats: 10,
	}
	flight2 := &models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "London",
		DepartureTime:  "2023-12-01T11:00:00Z",
		ArrivalTime:    "2023-12-01T15:00:00Z",
		Price:          120.00,
		AvailableSeats: 5,
	}

	// Добавляем данные в базу данных
	config.TestDB.Create(flight1)
	config.TestDB.Create(flight2)

	// Тестируем поиск рейсов
	tests := []struct {
		departureCity string
		arrivalCity   string
		expectedCount int
	}{
		{"Moscow", "Paris", 1},
		{"Moscow", "London", 1},
		{"Moscow", "New York", 0}, // Нет таких рейсов
	}

	for _, tt := range tests {
		t.Run(tt.departureCity+"-"+tt.arrivalCity, func(t *testing.T) {
			flights, err := SearchFlights(tt.departureCity, tt.arrivalCity)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCount, len(flights))
		})
	}
}
