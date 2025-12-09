package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BookingServiceTestSuite struct {
	suite.Suite
	db             *gorm.DB
	bookingService *BookingService
	flightService  *FlightService
	userService    *UserService
}

func (suite *BookingServiceTestSuite) SetupTest() {
	// Используем тестовую базу данных
	config.InitTestDB()
	suite.db = config.TestDB

	// Создаем сервисы с тестовой БД
	suite.flightService = NewFlightService().WithDB(suite.db)
	suite.userService = NewUserService().WithDB(suite.db)
	suite.bookingService = NewBookingService().WithDB(suite.db)

	// Устанавливаем зависимости
	suite.bookingService.flightService = suite.flightService
	suite.bookingService.userService = suite.userService

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

func (suite *BookingServiceTestSuite) TestBookFlight_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Тестируем создание бронирования
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")

	suite.NoError(err)
	suite.NotNil(booking)
	suite.Equal(flight.ID, booking.FlightID)
	suite.Equal(user.ID, booking.UserID)
	suite.Equal(2, booking.Seats)
	suite.Equal(models.BookingClassEconomy, booking.Class)
	suite.Equal(models.BookingStatusConfirmed, booking.Status)
	suite.NotEmpty(booking.BookingNumber)

	// Проверяем расчет цены
	expectedPrice := 5000.0 * 2 * 1.0 // 2 места * economy multiplier (1.0)
	suite.Equal(expectedPrice, booking.TotalPrice)

	// Проверяем, что количество мест уменьшилось
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(98, updatedFlight.AvailableSeats)
}

func (suite *BookingServiceTestSuite) TestBookFlight_NotEnoughSeats() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 5,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Пытаемся забронировать больше мест, чем доступно
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 10, "economy")

	suite.Error(err)
	suite.Nil(booking)
	suite.Equal("недостаточно мест на рейсе", err.Error())

	// Проверяем, что количество мест не изменилось
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(5, updatedFlight.AvailableSeats)
}

func (suite *BookingServiceTestSuite) TestBookFlight_FlightNotFound() {
	// Создаем только пользователя
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	// Пытаемся забронировать несуществующий рейс
	booking, err := suite.bookingService.BookFlight(999, user.ID, 2, "economy")

	suite.Error(err)
	suite.Nil(booking)
	suite.Equal("рейс не найден", err.Error())
}

func (suite *BookingServiceTestSuite) TestBookFlight_UserNotFound() {
	// Создаем только рейс
	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Пытаемся забронировать от несуществующего пользователя
	booking, err := suite.bookingService.BookFlight(flight.ID, 999, 2, "economy")

	suite.Error(err)
	suite.Nil(booking)
	suite.Equal("пользователь не найден", err.Error())
}

func (suite *BookingServiceTestSuite) TestBookFlight_InvalidClass() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Пытаемся забронировать с неверным классом
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "invalid_class")

	suite.Error(err)
	suite.Nil(booking)
	suite.Equal("неверный класс бронирования", err.Error())
}

func (suite *BookingServiceTestSuite) TestGetUserBookings_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем несколько бронирований
	booking1, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	booking2, err := suite.bookingService.BookFlight(flight.ID, user.ID, 1, "business")
	suite.NoError(err)

	// Получаем бронирования пользователя
	bookings, err := suite.bookingService.GetUserBookings(user.ID)

	suite.NoError(err)
	suite.Len(bookings, 2)

	// Проверяем, что бронирования содержат правильные данные
	foundBooking1 := false
	foundBooking2 := false

	for _, booking := range bookings {
		if booking.ID == booking1.ID {
			foundBooking1 = true
			suite.Equal(models.BookingClassEconomy, booking.Class)
			suite.Equal(2, booking.Seats)
		}
		if booking.ID == booking2.ID {
			foundBooking2 = true
			suite.Equal(models.BookingClassBusiness, booking.Class)
			suite.Equal(1, booking.Seats)
		}
	}

	suite.True(foundBooking1, "Booking 1 should be found")
	suite.True(foundBooking2, "Booking 2 should be found")
}

func (suite *BookingServiceTestSuite) TestGetUserBookings_UserNotFound() {
	// Пытаемся получить бронирования несуществующего пользователя
	bookings, err := suite.bookingService.GetUserBookings(999)

	suite.Error(err)
	suite.Nil(bookings)
	suite.Equal("пользователь не найден", err.Error())
}

func (suite *BookingServiceTestSuite) TestCancelBooking_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирование
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	// Отменяем бронирование
	err = suite.bookingService.CancelBooking(booking.ID)

	suite.NoError(err)

	// Проверяем, что статус изменился
	var cancelledBooking models.Booking
	suite.db.First(&cancelledBooking, booking.ID)
	suite.Equal(models.BookingStatusCancelled, cancelledBooking.Status)
	suite.NotZero(cancelledBooking.CancellationFee)

	// Проверяем, что места вернулись
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(100, updatedFlight.AvailableSeats)
}

func (suite *BookingServiceTestSuite) TestCancelBooking_AlreadyCancelled() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирование и отменяем его
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	err = suite.bookingService.CancelBooking(booking.ID)
	suite.NoError(err)

	// Пытаемся отменить уже отмененное бронирование
	err = suite.bookingService.CancelBooking(booking.ID)

	suite.Error(err)
	suite.Equal("бронирование уже отменено", err.Error())
}

func (suite *BookingServiceTestSuite) TestCancelBooking_NotFound() {
	// Пытаемся отменить несуществующее бронирование
	err := suite.bookingService.CancelBooking(999)

	suite.Error(err)
	suite.Equal("бронирование не найдено", err.Error())
}

func (suite *BookingServiceTestSuite) TestGetBookingByID_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирование
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	// Получаем бронирование по ID
	result, err := suite.bookingService.GetBookingByID(booking.ID)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(booking.ID, result.ID)
	suite.Equal(flight.ID, result.FlightID)
	suite.Equal(user.ID, result.UserID)
	suite.Equal(2, result.Seats)
	suite.Equal(models.BookingClassEconomy, result.Class)
}

func (suite *BookingServiceTestSuite) TestGetBookingByID_NotFound() {
	// Пытаемся получить несуществующее бронирование
	result, err := suite.bookingService.GetBookingByID(999)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal("бронирование не найдено", err.Error())
}

func (suite *BookingServiceTestSuite) TestGetBookingByNumber_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирование
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	// Получаем бронирование по номеру
	result, err := suite.bookingService.GetBookingByNumber(booking.BookingNumber)

	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(booking.ID, result.ID)
	suite.Equal(booking.BookingNumber, result.BookingNumber)
}

func (suite *BookingServiceTestSuite) TestUpdateBooking_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем начальное бронирование
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	// Обновляем бронирование
	updatedBooking, err := suite.bookingService.UpdateBooking(booking.ID, 3, "business")

	suite.NoError(err)
	suite.NotNil(updatedBooking)
	suite.Equal(3, updatedBooking.Seats)
	suite.Equal(models.BookingClassBusiness, updatedBooking.Class)

	// Пересчитываем ожидаемую цену
	expectedPrice := 5000.0 * 3 * 2.5 // 3 места * business multiplier (2.5)
	suite.Equal(expectedPrice, updatedBooking.TotalPrice)

	// Проверяем, что количество мест обновилось
	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(99, updatedFlight.AvailableSeats) // 100 - 1 (увеличили с 2 до 3 мест)
}

func (suite *BookingServiceTestSuite) TestUpdateBooking_NotEnoughSeats() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 2,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем начальное бронирование
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 1, "economy")
	suite.NoError(err)

	// Пытаемся увеличить количество мест больше чем доступно
	updatedBooking, err := suite.bookingService.UpdateBooking(booking.ID, 5, "economy")

	suite.Error(err)
	suite.Nil(updatedBooking)
	suite.Equal("недостаточно мест на рейсе", err.Error())
}

func (suite *BookingServiceTestSuite) TestUpdateBooking_CancelledBooking() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирование и отменяем его
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	err = suite.bookingService.CancelBooking(booking.ID)
	suite.NoError(err)

	// Пытаемся обновить отмененное бронирование
	updatedBooking, err := suite.bookingService.UpdateBooking(booking.ID, 3, "business")

	suite.Error(err)
	suite.Nil(updatedBooking)
	suite.Equal("нельзя изменить отмененное бронирование", err.Error())
}

func (suite *BookingServiceTestSuite) TestConfirmPayment_Success() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирование
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	// Подтверждаем оплату
	err = suite.bookingService.ConfirmPayment(booking.ID)

	suite.NoError(err)

	// Проверяем, что статус оплаты изменился
	var confirmedBooking models.Booking
	suite.db.First(&confirmedBooking, booking.ID)
	suite.Equal("paid", confirmedBooking.PaymentStatus)
}

func (suite *BookingServiceTestSuite) TestConfirmPayment_CancelledBooking() {
	// Создаем тестовые данные
	user := models.User{
		FirstName:    "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирование и отменяем его
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 2, "economy")
	suite.NoError(err)

	err = suite.bookingService.CancelBooking(booking.ID)
	suite.NoError(err)

	// Пытаемся подтвердить оплату отмененного бронирования
	err = suite.bookingService.ConfirmPayment(booking.ID)

	suite.Error(err)
	suite.Equal("только подтвержденные бронирования могут быть оплачены", err.Error())
}

func (suite *BookingServiceTestSuite) TestSearchBookings_Success() {
	// Создаем тестовые данные
	user1 := models.User{
		FirstName:    "User 1",
		Email:        "user1@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user1)

	user2 := models.User{
		FirstName:    "User 2",
		Email:        "user2@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user2)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "St. Petersburg",
		Price:          5000.0,
		AvailableSeats: 100,
		FlightNumber:   "SU100",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Создаем бронирования для разных пользователей
	booking1, err := suite.bookingService.BookFlight(flight.ID, user1.ID, 2, "economy")
	suite.NoError(err)

	booking2, err := suite.bookingService.BookFlight(flight.ID, user2.ID, 1, "business")
	suite.NoError(err)

	// Ищем бронирования пользователя 1
	params := BookingSearchParams{
		UserID: user1.ID,
	}
	bookings, err := suite.bookingService.SearchBookings(params)

	suite.NoError(err)
	suite.Len(bookings, 1)
	suite.Equal(booking1.ID, bookings[0].ID)

	// Ищем бронирования по классу
	params = BookingSearchParams{
		Class: "business",
	}
	bookings, err = suite.bookingService.SearchBookings(params)

	suite.NoError(err)
	suite.Len(bookings, 1)
	suite.Equal(booking2.ID, bookings[0].ID)
}

func (suite *BookingServiceTestSuite) TestCalculatePrice() {
	tests := []struct {
		name      string
		basePrice float64
		seats     int
		class     models.BookingClass
		expected  float64
	}{
		{
			name:      "Economy класс, 1 место",
			basePrice: 5000.0,
			seats:     1,
			class:     models.BookingClassEconomy,
			expected:  5000.0,
		},
		{
			name:      "Economy класс, 3 места",
			basePrice: 5000.0,
			seats:     3,
			class:     models.BookingClassEconomy,
			expected:  15000.0,
		},
		{
			name:      "Business класс, 1 место",
			basePrice: 5000.0,
			seats:     1,
			class:     models.BookingClassBusiness,
			expected:  12500.0, // 5000 * 2.5
		},
		{
			name:      "Business класс, 2 места",
			basePrice: 5000.0,
			seats:     2,
			class:     models.BookingClassBusiness,
			expected:  25000.0, // 5000 * 2 * 2.5
		},
		{
			name:      "First класс, 1 место",
			basePrice: 5000.0,
			seats:     1,
			class:     models.BookingClassFirst,
			expected:  20000.0, // 5000 * 4.0
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result := suite.bookingService.calculatePrice(tt.basePrice, tt.seats, tt.class)
			suite.Equal(tt.expected, result)
		})
	}
}

// Интеграционный тест полного потока бронирования
func (suite *BookingServiceTestSuite) TestBookingFlow() {
	// Шаг 1: Создаем пользователя и рейс
	user := models.User{
		FirstName:    "Integration Test User",
		Email:        "integration@example.com",
		PasswordHash: "hashed_password",
	}
	suite.db.Create(&user)

	flight := models.Flight{
		DepartureCity:  "Moscow",
		ArrivalCity:    "Sochi",
		Price:          8000.0,
		AvailableSeats: 50,
		FlightNumber:   "SU200",
		DepartureTime:  time.Now().Add(24 * time.Hour),
		ArrivalTime:    time.Now().Add(26 * time.Hour),
	}
	suite.db.Create(&flight)

	// Шаг 2: Создаем бронирование
	booking, err := suite.bookingService.BookFlight(flight.ID, user.ID, 3, "business")
	suite.NoError(err)
	suite.NotNil(booking)
	suite.Equal(models.BookingStatusConfirmed, booking.Status)

	// Шаг 3: Получаем бронирования пользователя
	bookings, err := suite.bookingService.GetUserBookings(user.ID)
	suite.NoError(err)
	suite.Len(bookings, 1)
	suite.Equal(booking.ID, bookings[0].ID)

	// Шаг 4: Подтверждаем оплату
	err = suite.bookingService.ConfirmPayment(booking.ID)
	suite.NoError(err)

	// Шаг 5: Отменяем бронирование
	err = suite.bookingService.CancelBooking(booking.ID)
	suite.NoError(err)

	// Проверяем, что статус изменился и места вернулись
	var cancelledBooking models.Booking
	suite.db.First(&cancelledBooking, booking.ID)
	suite.Equal(models.BookingStatusCancelled, cancelledBooking.Status)
	suite.NotZero(cancelledBooking.CancellationFee)

	var updatedFlight models.Flight
	suite.db.First(&updatedFlight, flight.ID)
	suite.Equal(50, updatedFlight.AvailableSeats)
}
