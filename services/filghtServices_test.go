// services/flight_service_test.go
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
	db            *gorm.DB
	flightService *FlightService
}

func (suite *FlightServiceTestSuite) SetupTest() {
	// Используем тестовую базу данных
	config.InitTestDB()
	suite.db = config.TestDB
	suite.flightService = NewFlightService().WithDB(suite.db)

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
			params := FlightSearchParams{
				DepartureCity: tt.departureCity,
				ArrivalCity:   tt.arrivalCity,
			}

			result, err := suite.flightService.SearchFlights(params)

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

func (suite *FlightServiceTestSuite) TestSearchFlights_ExtendedFilters() {
	// Создаем тестовые данные
	departureTime1, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime1, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")
	departureTime2, _ := time.Parse(time.RFC3339, "2023-12-02T11:00:00Z")
	arrivalTime2, _ := time.Parse(time.RFC3339, "2023-12-02T15:00:00Z")
	departureTime3, _ := time.Parse(time.RFC3339, "2023-12-03T09:00:00Z")
	arrivalTime3, _ := time.Parse(time.RFC3339, "2023-12-03T12:00:00Z")

	flights := []models.Flight{
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "Paris",
			DepartureTime:  departureTime1,
			ArrivalTime:    arrivalTime1,
			Price:          150.00,
			AvailableSeats: 10,
			Airline:        "Aeroflot",
		},
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "London",
			DepartureTime:  departureTime2,
			ArrivalTime:    arrivalTime2,
			Price:          200.00,
			AvailableSeats: 5,
			Airline:        "British Airways",
		},
		{
			DepartureCity:  "Paris",
			ArrivalCity:    "Berlin",
			DepartureTime:  departureTime3,
			ArrivalTime:    arrivalTime3,
			Price:          80.00,
			AvailableSeats: 8,
			Airline:        "Air France",
		},
	}

	// Сохраняем данные в базу
	for i := range flights {
		suite.db.Create(&flights[i])
	}

	// Тест с фильтром по времени
	fromTime, _ := time.Parse(time.RFC3339, "2023-12-01T00:00:00Z")
	toTime, _ := time.Parse(time.RFC3339, "2023-12-02T23:59:59Z")

	params := FlightSearchParams{
		DepartureTimeFrom: &fromTime,
		DepartureTimeTo:   &toTime,
	}

	result, err := suite.flightService.SearchFlights(params)
	suite.NoError(err)
	suite.Len(result, 2) // Должны найти 2 рейса за 1-2 декабря

	// Тест с фильтром по цене
	params = FlightSearchParams{
		MinPrice: 100,
		MaxPrice: 180,
	}

	result, err = suite.flightService.SearchFlights(params)
	suite.NoError(err)
	suite.Len(result, 1) // Должен найти 1 рейс с ценой 150
	suite.Equal(150.00, result[0].Price)

	// Тест с фильтром по авиакомпании
	params = FlightSearchParams{
		Airline: "Aeroflot",
	}

	result, err = suite.flightService.SearchFlights(params)
	suite.NoError(err)
	suite.Len(result, 1)
	suite.Equal("Aeroflot", result[0].Airline)
}

func (suite *FlightServiceTestSuite) TestSearchFlights_EmptyDatabase() {
	// Тестируем поиск в пустой базе данных
	params := FlightSearchParams{
		DepartureCity: "Moscow",
		ArrivalCity:   "Paris",
	}

	result, err := suite.flightService.SearchFlights(params)

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
	result, err := suite.flightService.GetAllFlights()

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
	result, err := suite.flightService.GetFlightByID(flight.ID)

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
	result, err := suite.flightService.GetFlightByID(999)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal("рейс не найден", err.Error())
}

func (suite *FlightServiceTestSuite) TestCreateFlight_Success() {
	// Создаем новый рейс
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flight := &models.Flight{
		FlightNumber:   "SU123",
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 10,
	}

	err := suite.flightService.CreateFlight(flight)

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

func (suite *FlightServiceTestSuite) TestCreateFlight_ValidationError() {
	// Тест с некорректными данными
	tests := []struct {
		name          string
		flight        *models.Flight
		expectedError string
	}{
		{
			name: "Пустой номер рейса",
			flight: &models.Flight{
				DepartureCity:  "Moscow",
				ArrivalCity:    "Paris",
				DepartureTime:  time.Now(),
				ArrivalTime:    time.Now().Add(2 * time.Hour),
				Price:          150.00,
				AvailableSeats: 10,
			},
			expectedError: "номер рейса обязателен",
		},
		{
			name: "Пустой город вылета",
			flight: &models.Flight{
				FlightNumber:   "SU123",
				ArrivalCity:    "Paris",
				DepartureTime:  time.Now(),
				ArrivalTime:    time.Now().Add(2 * time.Hour),
				Price:          150.00,
				AvailableSeats: 10,
			},
			expectedError: "город вылета обязателен",
		},
		{
			name: "Пустой город прибытия",
			flight: &models.Flight{
				FlightNumber:   "SU123",
				DepartureCity:  "Moscow",
				DepartureTime:  time.Now(),
				ArrivalTime:    time.Now().Add(2 * time.Hour),
				Price:          150.00,
				AvailableSeats: 10,
			},
			expectedError: "город прибытия обязателен",
		},
		{
			name: "Некорректное время",
			flight: &models.Flight{
				FlightNumber:   "SU123",
				DepartureCity:  "Moscow",
				ArrivalCity:    "Paris",
				DepartureTime:  time.Now().Add(2 * time.Hour),
				ArrivalTime:    time.Now(),
				Price:          150.00,
				AvailableSeats: 10,
			},
			expectedError: "время вылета не может быть позже времени прибытия",
		},
		{
			name: "Отрицательная цена",
			flight: &models.Flight{
				FlightNumber:   "SU123",
				DepartureCity:  "Moscow",
				ArrivalCity:    "Paris",
				DepartureTime:  time.Now(),
				ArrivalTime:    time.Now().Add(2 * time.Hour),
				Price:          -100.00,
				AvailableSeats: 10,
			},
			expectedError: "цена должна быть положительной",
		},
		{
			name: "Отрицательное количество мест",
			flight: &models.Flight{
				FlightNumber:   "SU123",
				DepartureCity:  "Moscow",
				ArrivalCity:    "Paris",
				DepartureTime:  time.Now(),
				ArrivalTime:    time.Now().Add(2 * time.Hour),
				Price:          150.00,
				AvailableSeats: -5,
			},
			expectedError: "количество мест не может быть отрицательным",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := suite.flightService.CreateFlight(tt.flight)
			suite.Error(err)
			suite.Contains(err.Error(), tt.expectedError)
		})
	}
}

func (suite *FlightServiceTestSuite) TestUpdateFlight_Success() {
	// Создаем тестовый рейс
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flight := models.Flight{
		FlightNumber:   "SU123",
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 10,
	}
	suite.db.Create(&flight)

	// Обновляем рейс
	updatedFlight := &models.Flight{
		FlightNumber:   "SU456",
		DepartureCity:  "Moscow",
		ArrivalCity:    "London",
		DepartureTime:  departureTime.Add(24 * time.Hour),
		ArrivalTime:    arrivalTime.Add(24 * time.Hour),
		Price:          200.00,
		AvailableSeats: 5,
	}

	err := suite.flightService.UpdateFlight(flight.ID, updatedFlight)
	suite.NoError(err)

	// Проверяем обновление
	var result models.Flight
	suite.db.First(&result, flight.ID)
	suite.Equal("SU456", result.FlightNumber)
	suite.Equal("London", result.ArrivalCity)
	suite.Equal(200.00, result.Price)
	suite.Equal(5, result.AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestUpdateFlight_NotFound() {
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	updatedFlight := &models.Flight{
		FlightNumber:   "SU123",
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 10,
	}

	err := suite.flightService.UpdateFlight(999, updatedFlight)
	suite.Error(err)
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
	err := suite.flightService.UpdateFlightSeats(flight.ID, 5)

	suite.NoError(err)

	// Проверяем обновление
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(5, updatedFlight.AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestUpdateFlightSeats_InvalidSeats() {
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

	// Пытаемся установить отрицательное количество мест
	err := suite.flightService.UpdateFlightSeats(flight.ID, -5)
	suite.Error(err)
	suite.Equal("количество мест не может быть отрицательным", err.Error())
}

func (suite *FlightServiceTestSuite) TestUpdateFlightSeats_NotFound() {
	// Пытаемся обновить несуществующий рейс
	err := suite.flightService.UpdateFlightSeats(999, 5)

	suite.Error(err)
}

func (suite *FlightServiceTestSuite) TestDeleteFlight_Success() {
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

	// Удаляем рейс
	err := suite.flightService.DeleteFlight(flight.ID)
	suite.NoError(err)

	// Проверяем, что рейс удален
	var result models.Flight
	err = suite.db.First(&result, flight.ID).Error
	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *FlightServiceTestSuite) TestDeleteFlight_NotFound() {
	// Пытаемся удалить несуществующий рейс
	err := suite.flightService.DeleteFlight(999)
	suite.Error(err)
	suite.Equal("рейс не найден", err.Error())
}

func (suite *FlightServiceTestSuite) TestCheckFlightAvailability_Success() {
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

	// Проверяем доступность мест
	available, err := suite.flightService.CheckFlightAvailability(flight.ID, 5)
	suite.NoError(err)
	suite.True(available)

	available, err = suite.flightService.CheckFlightAvailability(flight.ID, 15)
	suite.NoError(err)
	suite.False(available)
}

func (suite *FlightServiceTestSuite) TestReserveSeats_Success() {
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

	// Резервируем места
	err := suite.flightService.ReserveSeats(flight.ID, 5)
	suite.NoError(err)

	// Проверяем, что места зарезервированы
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(5, updatedFlight.AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestReserveSeats_NotEnoughSeats() {
	// Создаем тестовый рейс
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 3,
	}
	suite.db.Create(&flight)

	// Пытаемся зарезервировать больше мест, чем доступно
	err := suite.flightService.ReserveSeats(flight.ID, 5)
	suite.Error(err)
	suite.Equal("недостаточно мест на рейсе", err.Error())
}

func (suite *FlightServiceTestSuite) TestReleaseSeats_Success() {
	// Создаем тестовый рейс
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Paris",
		DepartureTime:  departureTime,
		ArrivalTime:    arrivalTime,
		Price:          150.00,
		AvailableSeats: 5,
	}
	suite.db.Create(&flight)

	// Освобождаем места
	err := suite.flightService.ReleaseSeats(flight.ID, 3)
	suite.NoError(err)

	// Проверяем, что места освобождены
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(8, updatedFlight.AvailableSeats)
}

func (suite *FlightServiceTestSuite) TestGetUpcomingFlights() {
	now := time.Now()

	// Создаем тестовые рейсы
	flights := []models.Flight{
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "Paris",
			DepartureTime:  now.Add(-2 * time.Hour), // Прошедший
			ArrivalTime:    now.Add(-1 * time.Hour),
			Price:          150.00,
			AvailableSeats: 10,
		},
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "London",
			DepartureTime:  now.Add(1 * time.Hour), // Будущий
			ArrivalTime:    now.Add(3 * time.Hour),
			Price:          120.00,
			AvailableSeats: 5,
		},
		{
			DepartureCity:  "Paris",
			ArrivalCity:    "Berlin",
			DepartureTime:  now.Add(2 * time.Hour), // Будущий
			ArrivalTime:    now.Add(4 * time.Hour),
			Price:          80.00,
			AvailableSeats: 8,
		},
	}

	// Сохраняем данные в базу
	for i := range flights {
		suite.db.Create(&flights[i])
	}

	// Получаем предстоящие рейсы
	result, err := suite.flightService.GetUpcomingFlights(10)
	suite.NoError(err)
	suite.Len(result, 2) // Должны найти 2 будущих рейса
}

func (suite *FlightServiceTestSuite) TestGetFlightsByAirline() {
	// Создаем тестовые данные
	departureTime, _ := time.Parse(time.RFC3339, "2023-12-01T10:00:00Z")
	arrivalTime, _ := time.Parse(time.RFC3339, "2023-12-01T14:00:00Z")

	flights := []models.Flight{
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "Paris",
			DepartureTime:  departureTime,
			ArrivalTime:    arrivalTime,
			Price:          150.00,
			AvailableSeats: 10,
			Airline:        "Aeroflot",
		},
		{
			DepartureCity:  "Moscow",
			ArrivalCity:    "London",
			DepartureTime:  departureTime,
			ArrivalTime:    arrivalTime,
			Price:          120.00,
			AvailableSeats: 5,
			Airline:        "Aeroflot",
		},
		{
			DepartureCity:  "Paris",
			ArrivalCity:    "Berlin",
			DepartureTime:  departureTime,
			ArrivalTime:    arrivalTime,
			Price:          80.00,
			AvailableSeats: 8,
			Airline:        "Air France",
		},
	}

	// Сохраняем данные в базу
	for i := range flights {
		suite.db.Create(&flights[i])
	}

	// Получаем рейсы авиакомпании Aeroflot
	result, err := suite.flightService.GetFlightsByAirline("Aeroflot")
	suite.NoError(err)
	suite.Len(result, 2)

	for _, flight := range result {
		suite.Equal("Aeroflot", flight.Airline)
	}
}
