package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type FlightServiceTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *FlightServiceTestSuite) SetupTest() {
	// Используем тестовую базу данных
	config.InitTestDB()
	suite.db = config.TestDB

	// Очищаем таблицы перед каждым тестом
	suite.db.Exec("DELETE FROM flights")
}

func (suite *FlightServiceTestSuite) TearDownTest() {
	config.CloseTestDB()
}

func TestFlightServiceSuite(t *testing.T) {
	suite.Run(t, new(FlightServiceTestSuite))
}

func (suite *FlightServiceTestSuite) TestSearchFlights_Success() {
	// Создаем тестовые данные с правильным форматом времени
	departureTime1, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime1, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")
	departureTime2, _ := time.Parse(time.RFC3339, "2023-12-01T11:00:00Z")
	arrivalTime2, _ := time.Parse(time.RFC3339, "2023-12-01T15:00:00Z")
	departureTime3, _ := time.Parse(time.RFC3339, "2023-12-02T09:00:00Z")
	arrivalTime3, _ := time.Parse(time.RFC3339, "2023-12-02T12:00:00Z")

	flights := []models.Flight{
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "Paris",
			DepartureTime:  departureTime1,
			ArrivalTime:    arrivalTime1,
			Price:          150.00,
			AvailableSeats: 10,
		},
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "London",
			DepartureTime:  departureTime2,
			ArrivalTime:    arrivalTime2,
			Price:          120.00,
			AvailableSeats: 5,
		},
		{
			DepartureCity:  "Paris",
			ArrivalCity:    "Berlin",
			DepartureTime:  departureTime3,
			ArrivalTime:    arrivalTime3,
			Price:          80.00,
			AvailableSeats: 8,
		},
	}

	// Сохраняем данные в базу
	for i := range flights {
		suite.db.Create(&flights[i])
	}

	tests := []struct {
		name          string
		departureCity string
		arrivalCity   string
		expectedCount int
	}{
		{
			name:          "Поиск рейсов Moscow-Paris",
			departureCity: "Moscow",
			arrivalCity:   "Paris",
			expectedCount: 1,
		},
		{
			name:          "Поиск рейсов Moscow-London",
			departureCity: "Moscow",
			arrivalCity:   "London",
			expectedCount: 1,
		},
		{
			name:          "Поиск рейсов Paris-Berlin",
			departureCity: "Paris",
			arrivalCity:   "Berlin",
			expectedCount: 1,
		},
		{
			name:          "Поиск рейсов Moscow-Berlin (нет рейсов)",
			departureCity: "Moscow",
			arrivalCity:   "Berlin",
			expectedCount: 0,
		},
		{
			name:          "Поиск всех рейсов из Moscow",
			departureCity: "Moscow",
			arrivalCity:   "",
			expectedCount: 2,
		},
		{
			name:          "Поиск всех рейсов в Paris",
			departureCity: "",
			arrivalCity:   "Paris",
			expectedCount: 1,
		},
		{
			name:          "Поиск всех рейсов (без фильтров)",
			departureCity: "",
			arrivalCity:   "",
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := SearchFlights(tt.departureCity, tt.arrivalCity)

			suite.NoError(err)
			suite.Len(result, tt.expectedCount)

			// Проверяем, что результаты соответствуют фильтрам
			for _, flight := range result {
				if tt.departureCity != "" {
					suite.Equal(tt.departureCity, flight.DepartureCity)
				}
				if tt.arrivalCity != "" {
					suite.Equal(tt.arrivalCity, flight.ArrivalCity)
				}
			}
		})
	}
}

func (suite *FlightServiceTestSuite) TestSearchFlights_EmptyDatabase() {
	// Тестируем поиск в пустой базе данных
	result, err := SearchFlights("Moscow", "Paris")

	suite.NoError(err)
	suite.Empty(result)
}

func (suite *FlightServiceTestSuite) TestGetAllFlights() {
	// Создаем тестовые данные
	departureTime1, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime1, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")
	departureTime2, _ := time.Parse(time.RFC3339, "2023-12-01T11:00:00Z")
	arrivalTime2, _ := time.Parse(time.RFC3339, "2023-12-01T15:00:00Z")

	flights := []models.Flight{
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "Paris",
			DepartureTime:  departureTime1,
			ArrivalTime:    arrivalTime1,
			Price:          150.00,
			AvailableSeats: 10,
		},
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "London",
			DepartureTime:  departureTime2,
			ArrivalTime:    arrivalTime2,
			Price:          120.00,
			AvailableSeats: 5,
		},
	}

	// Сохраняем данные в базу
	for i := range flights {
		suite.db.Create(&flights[i])
	}

	// Получаем все рейсы
	result, err := GetAllFlights()

	suite.NoError(err)
	suite.Len(result, 2)

	// Проверяем, что данные корректны
	suite.Equal("Moscow", result[0].DepartureCity)
	suite.Equal("Paris", result[0].ArrivalCity)
	suite.Equal(150.00, result[0].Price)
	suite.Equal(10, result[0].AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestGetFlightByID_Success() {
	// Создаем тестовый рейс
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 10,
	}
	suite.db.Create(&flight)

	// Получаем рейс по ID
	result, err := GetFlightByID(flight.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(flight.ID, result.ID)
	suite.Equal("Moscow", result.DepartureCity)
	suite.Equal("Paris", result.ArrivalCity)
	suite.Equal(150.00, result.Price)
	suite.Equal(10, result.AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestGetFlightByID_NotFound() {
	// Пытаемся получить несуществующий рейс
	result, err := GetFlightByID(999)

	suite.Error(err)
	suite.Nil(result)
}

func (suite *FlightServiceTestSuite) TestCreateFlight_Success() {
	// Создаем новый рейс
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flight := &models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 10,
	}

	err := CreateFlight(flight)

	suite.NoError(err)
	suite.NotZero(flight.ID) // ID должен быть установлен

	// Проверяем, что рейс сохранен в базе
	var savedFlight models.Flight
	suite.db.First(&savedFlight, flight.ID)
	suite.Equal("Moscow", savedFlight.DepartureCity)
	suite.Equal("Paris", savedFlight.ArrivalCity)
	suite.Equal(150.00, savedFlight.Price)
	suite.Equal(10, savedFlight.AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestUpdateFlightSeats_Success() {
	// Создаем тестовый рейс
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 10,
	}
	suite.db.Create(&flight)

	// Обновляем количество мест
	err := UpdateFlightSeats(flight.ID, 5)

	suite.NoError(err)

	// Проверяем обновление
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(5, updatedFlight.AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestUpdateFlightSeats_NotFound() {
	// Пытаемся обновить несуществующий рейс
	err := UpdateFlightSeats(999, 5)

	suite.Error(err)
}
