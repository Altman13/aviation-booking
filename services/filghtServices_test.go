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

	// Тестируем поиск рейсов с правильным количеством параметров
	tests := []struct {
		name          string
		origin        string
		destination   string
		date          string
		adults        string
		class         string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Успешный поиск рейсов Moscow-Paris",
			origin:        "Moscow",
			destination:   "Paris",
			date:          "2023-12-01",
			adults:        "1",
			class:         "economy",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Успешный поиск рейсов Moscow-London",
			origin:        "Moscow",
			destination:   "London",
			date:          "2023-12-01",
			adults:        "2",
			class:         "business",
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "Поиск несуществующих рейсов",
			origin:        "Moscow",
			destination:   "New York",
			date:          "2023-12-01",
			adults:        "1",
			class:         "economy",
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "Пустые обязательные параметры",
			origin:        "",
			destination:   "Paris",
			date:          "2023-12-01",
			adults:        "1",
			class:         "economy",
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flights, err := SearchFlights(tt.origin, tt.destination, tt.date, tt.adults, tt.class)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(flights))
			}
		})
	}
}

func TestSearchFlights_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		origin      string
		destination string
		date        string
		adults      string
		class       string
		expectError bool
	}{
		{
			name:        "Некорректное количество пассажиров",
			origin:      "Moscow",
			destination: "Paris",
			date:        "2023-12-01",
			adults:      "invalid",
			class:       "economy",
			expectError: false, // Функция должна обработать это
		},
		{
			name:        "Пустая дата",
			origin:      "Moscow",
			destination: "Paris",
			date:        "",
			adults:      "1",
			class:       "economy",
			expectError: false,
		},
		{
			name:        "Класс first",
			origin:      "Moscow",
			destination: "Paris",
			date:        "2023-12-01",
			adults:      "2",
			class:       "first",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flights, err := SearchFlights(tt.origin, tt.destination, tt.date, tt.adults, tt.class)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, flights)
			}
		})
	}
}

func TestGetFlightByID(t *testing.T) {
	tests := []struct {
		name        string
		flightID    string
		expectError bool
	}{
		{
			name:        "Успешное получение рейса по ID",
			flightID:    "SU-1234",
			expectError: false,
		},
		{
			name:        "Получение рейса с другим ID",
			flightID:    "S7-5678",
			expectError: false,
		},
		{
			name:        "Пустой ID",
			flightID:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flight, err := GetFlightByID(tt.flightID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, flight)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, flight)
				assert.Equal(t, tt.flightID, flight.ID)
				assert.NotEmpty(t, flight.FlightNumber)
				assert.NotEmpty(t, flight.Airline)
				assert.NotEmpty(t, flight.DepartureAirport)
				assert.NotEmpty(t, flight.ArrivalAirport)
				assert.Greater(t, flight.Price, 0.0)
			}
		})
	}
}

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		name       string
		basePrice  int
		passengers int
		class      string
		expected   float64
	}{
		{
			name:       "Economy класс, 1 пассажир",
			basePrice:  5000,
			passengers: 1,
			class:      "economy",
			expected:   5000.0,
		},
		{
			name:       "Economy класс, 2 пассажира",
			basePrice:  5000,
			passengers: 2,
			class:      "economy",
			expected:   10000.0,
		},
		{
			name:       "Business класс, 1 пассажир",
			basePrice:  5000,
			passengers: 1,
			class:      "business",
			expected:   12500.0, // 5000 * 2.5
		},
		{
			name:       "Business класс, 2 пассажира",
			basePrice:  5000,
			passengers: 2,
			class:      "business",
			expected:   25000.0, // 5000 * 2 * 2.5
		},
		{
			name:       "First класс, 1 пассажир",
			basePrice:  5000,
			passengers: 1,
			class:      "first",
			expected:   25000.0, // 5000 * 5.0
		},
		{
			name:       "First класс, 2 пассажира",
			basePrice:  5000,
			passengers: 2,
			class:      "first",
			expected:   50000.0, // 5000 * 2 * 5.0
		},
		{
			name:       "Неизвестный класс",
			basePrice:  5000,
			passengers: 1,
			class:      "unknown",
			expected:   5000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePrice(tt.basePrice, tt.passengers, tt.class)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSearchFlights_Integration(t *testing.T) {
	// Интеграционный тест с реальными параметрами поиска
	flights, err := SearchFlights("SVO", "LED", "2024-01-15", "2", "business")

	assert.NoError(t, err)
	assert.NotNil(t, flights)

	if len(flights) > 0 {
		flight := flights[0]
		assert.NotEmpty(t, flight.ID)
		assert.NotEmpty(t, flight.FlightNumber)
		assert.NotEmpty(t, flight.Airline)
		assert.Equal(t, "SVO", flight.DepartureAirport)
		assert.Equal(t, "LED", flight.ArrivalAirport)
		assert.Contains(t, flight.DepartureTime, "2024-01-15")
		assert.Greater(t, flight.Price, 0.0)
		assert.Greater(t, flight.AvailableSeats, 0)
	}
}
