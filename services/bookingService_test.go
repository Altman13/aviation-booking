package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BookingServiceTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *BookingServiceTestSuite) SetupTest() {
	// Используем тестовую базу данных
	config.InitTestDB()
	suite.db = config.TestDB

	// Очищаем таблицы перед каждым тестом
	suite.db.Exec("DELETE FROM bookings")
	suite.db.Exec("DELETE FROM flights")
	suite.db.Exec("DELETE FROM users")
}

func (suite *BookingServiceTestSuite) TearDownTest() {
	config.CloseTestDB()
}

func TestBookingServiceSuite(t *testing.T) {
	suite.Run(t, new(BookingServiceTestSuite))
}

func (suite *BookingServiceTestSuite) TestCreateBooking_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	// Тестируем создание бронирования
	booking, err := CreateBooking(flight.ID, user.ID, 2, "economy")

	suite.NoError(err)
	suite.NotNil(booking)
	suite.Equal(flight.ID, booking.FlightID)
	suite.Equal(user.ID, booking.UserID)
	suite.Equal(2, booking.Seats)
	suite.Equal("economy", booking.Class)
	suite.Equal("confirmed", booking.Status)
	suite.Equal(10000.0, booking.TotalPrice) // 5000 * 2 * 1.0

	// Проверяем, что количество мест уменьшилось
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(98, updatedFlight.AvailableSeats)
}

func (suite *BookingServiceTestSuite) TestCreateBooking_NotEnoughSeats() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 5,
	}
	suite.db.Create(&flight)

	// Пытаемся забронировать больше мест, чем доступно
	booking, err := CreateBooking(flight.ID, user.ID, 10, "economy")

	suite.Error(err)
	suite.Nil(booking)
	suite.Equal("not enough available seats", err.Error())

	// Проверяем, что количество мест не изменилось
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(5, updatedFlight.AvailableSeats)
}

func (suite *BookingServiceTestSuite) TestCreateBooking_FlightNotFound() {
	// Создаем только пользователя
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	// Пытаемся забронировать несуществующий рейс
	booking, err := CreateBooking(999, user.ID, 2, "economy")

	suite.Error(err)
	suite.Nil(booking)
	suite.Equal("flight not found", err.Error())
}

func (suite *BookingServiceTestSuite) TestCreateBooking_UserNotFound() {
	// Создаем только рейс
	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	// Пытаемся забронировать от несуществующего пользователя
	booking, err := CreateBooking(flight.ID, 999, 2, "economy")

	suite.Error(err)
	suite.Nil(booking)
	suite.Equal("user not found", err.Error())
}

func (suite *BookingServiceTestSuite) TestCreateBooking_DefaultClass() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	// Создаем бронирование без указания класса
	booking, err := CreateBooking(flight.ID, user.ID, 2, "")

	suite.NoError(err)
	suite.NotNil(booking)
	suite.Equal("economy", booking.Class) // Должен использоваться класс по умолчанию
}

func (suite *BookingServiceTestSuite) TestGetUserBookings_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	// Создаем несколько бронирований
	booking1 := models.Booking{
		FlightID:   flight.ID,
		UserID:     user.ID,
		Seats:      2,
		Class:      "economy",
		TotalPrice: 10000.0,
		Status:     "confirmed",
	}
	suite.db.Create(&booking1)

	booking2 := models.Booking{
		FlightID:   flight.ID,
		UserID:     user.ID,
		Seats:      1,
		Class:      "business",
		TotalPrice: 12500.0,
		Status:     "confirmed",
	}
	suite.db.Create(&booking2)

	// Получаем бронирования пользователя
	bookings, err := GetUserBookings(user.ID)

	suite.NoError(err)
	suite.Len(bookings, 2)
	suite.Equal(booking1.ID, bookings[0].ID)
	suite.Equal(booking2.ID, bookings[1].ID)
	suite.Equal("economy", bookings[0].Class)
	suite.Equal("business", bookings[1].Class)
}

func (suite *BookingServiceTestSuite) TestGetUserBookings_UserNotFound() {
	// Пытаемся получить бронирования несуществующего пользователя
	bookings, err := GetUserBookings(999)

	suite.Error(err)
	suite.Nil(bookings)
	suite.Equal("user not found", err.Error())
}

func (suite *BookingServiceTestSuite) TestCancelBooking_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	booking := models.Booking{
		FlightID:   flight.ID,
		UserID:     user.ID,
		Seats:      2,
		Class:      "economy",
		TotalPrice: 10000.0,
		Status:     "confirmed",
	}
	suite.db.Create(&booking)

	// Отменяем бронирование
	err := CancelBooking(booking.ID)

	suite.NoError(err)

	// Проверяем, что статус изменился
	var cancelledBooking models.Booking
	suite.db.First(&cancelledBooking, booking.ID)
	suite.Equal("cancelled", cancelledBooking.Status)

	// Проверяем, что места вернулись
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(102, updatedFlight.AvailableSeats) // 100 + 2
}

func (suite *BookingServiceTestSuite) TestCancelBooking_AlreadyCancelled() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	booking := models.Booking{
		FlightID:   flight.ID,
		UserID:     user.ID,
		Seats:      2,
		Class:      "economy",
		TotalPrice: 10000.0,
		Status:     "cancelled", // Уже отменено
	}
	suite.db.Create(&booking)

	// Пытаемся отменить уже отмененное бронирование
	err := CancelBooking(booking.ID)

	suite.Error(err)
	suite.Equal("booking already cancelled", err.Error())
}

func (suite *BookingServiceTestSuite) TestCancelBooking_NotFound() {
	// Пытаемся отменить несуществующее бронирование
	err := CancelBooking(999)

	suite.Error(err)
	suite.Equal("booking not found", err.Error())
}

func (suite *BookingServiceTestSuite) TestGetBookingByID_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	booking := models.Booking{
		FlightID:   flight.ID,
		UserID:     user.ID,
		Seats:      2,
		Class:      "economy",
		TotalPrice: 10000.0,
		Status:     "confirmed",
	}
	suite.db.Create(&booking)

	// Получаем бронирование по ID
	result, err := GetBookingByID(booking.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(booking.ID, result.ID)
	suite.Equal(flight.ID, result.FlightID)
	suite.Equal(user.ID, result.UserID)
}

func (suite *BookingServiceTestSuite) TestGetBookingByID_NotFound() {
	// Пытаемся получить несуществующее бронирование
	result, err := GetBookingByID(999)

	suite.Error(err)
	suite.Nil(result)
}

func (suite *BookingServiceTestSuite) TestUpdateBooking_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName: "Test User",
		Email:     "test@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
	}
	suite.db.Create(&flight)

	booking := models.Booking{
		FlightID:   flight.ID,
		UserID:     user.ID,
		Seats:      2,
		Class:      "economy",
		TotalPrice: 10000.0,
		Status:     "confirmed",
	}
	suite.db.Create(&booking)

	// Обновляем бронирование
	updatedBooking, err := UpdateBooking(booking.ID, 3, "business")

	suite.NoError(err)
	suite.NotNil(updatedBooking)
	suite.Equal(3, updatedBooking.Seats)
	suite.Equal("business", updatedBooking.Class)
	suite.Equal(37500.0, updatedBooking.TotalPrice) // 5000 * 3 * 2.5

	// Проверяем, что количество мест обновилось
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(99, updatedFlight.AvailableSeats) // 100 - 1 (увеличили с 2 до 3 мест)
}

func (suite *BookingServiceTestSuite) TestCalculateTotalCost() {
	tests := []struct {
		name      string
		basePrice float64
		seats     int
		class     string
		expected  float64
	}{
		{
			name:      "Economy класс, 1 место",
			basePrice: 5000.0,
			seats:     1,
			class:     "economy",
			expected:  5000.0,
		},
		{
			name:      "Economy класс, 3 места",
			basePrice: 5000.0,
			seats:     3,
			class:     "economy",
			expected:  15000.0,
		},
		{
			name:      "Business класс, 1 место",
			basePrice: 5000.0,
			seats:     1,
			class:     "business",
			expected:  12500.0,
		},
		{
			name:      "Business класс, 2 места",
			basePrice: 5000.0,
			seats:     2,
			class:     "business",
			expected:  25000.0,
		},
		{
			name:      "First класс, 1 место",
			basePrice: 5000.0,
			seats:     1,
			class:     "first",
			expected:  20000.0,
		},
		{
			name:      "Comfort класс, 2 места",
			basePrice: 5000.0,
			seats:     2,
			class:     "comfort",
			expected:  15000.0,
		},
		{
			name:      "Неизвестный класс (должен использовать economy)",
			basePrice: 5000.0,
			seats:     2,
			class:     "unknown",
			expected:  10000.0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := calculateTotalCost(tt.basePrice, tt.seats, tt.class)
			suite.Equal(tt.expected, result)
		})
	}
}

// Интеграционный тест полного потока бронирования
func (suite *BookingServiceTestSuite) TestBookingFlow() {
	// Шаг 1: Создаем пользователя и рейс
	user := models.User{
		FirstName: "Integration Test User",
		Email:     "integration@example.com",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Sochi",
		Price:          8000.0,
		AvailableSeats: 50,
	}
	suite.db.Create(&flight)

	// Шаг 2: Создаем бронирование
	booking, err := CreateBooking(flight.ID, user.ID, 3, "business")
	suite.NoError(err)
	suite.NotNil(booking)
	suite.Equal("confirmed", booking.Status)

	// Шаг 3: Получаем бронирования пользователя
	bookings, err := GetUserBookings(user.ID)
	suite.NoError(err)
	suite.Len(bookings, 1)
	suite.Equal(booking.ID, bookings[0].ID)

	// Шаг 4: Отменяем бронирование
	err = CancelBooking(booking.ID)
	suite.NoError(err)

	// Проверяем, что статус изменился и места вернулись
	var cancelledBooking models.Booking
	suite.db.First(&cancelledBooking, booking.ID)
	suite.Equal("cancelled", cancelledBooking.Status)

	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(50, updatedFlight.AvailableSeats)
}
