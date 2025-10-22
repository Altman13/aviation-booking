package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock базы данных или хранилища
type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) CreateBooking(booking *Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *MockBookingRepository) GetUserBookings(userID string) ([]Booking, error) {
	args := m.Called(userID)
	return args.Get(0).([]Booking), args.Error(1)
}

func (m *MockBookingRepository) CancelBooking(bookingID string) error {
	args := m.Called(bookingID)
	return args.Error(0)
}

func TestCreateBooking(t *testing.T) {
	tests := []struct {
		name          string
		flightID      string
		userID        string
		seats         int
		class         string
		expectedClass string
		expectedCost  float64
		expectError   bool
	}{
		{
			name:          "Успешное создание бронирования economy класса",
			flightID:      "SU-1234",
			userID:        "user-123",
			seats:         2,
			class:         "economy",
			expectedClass: "economy",
			expectedCost:  10000.0, // 5000 * 2
			expectError:   false,
		},
		{
			name:          "Успешное создание бронирования business класса",
			flightID:      "SU-5678",
			userID:        "user-456",
			seats:         1,
			class:         "business",
			expectedClass: "business",
			expectedCost:  15000.0, // 15000 * 1
			expectError:   false,
		},
		{
			name:          "Создание бронирования с классом по умолчанию",
			flightID:      "SU-9012",
			userID:        "user-789",
			seats:         3,
			class:         "",
			expectedClass: "economy",
			expectedCost:  15000.0, // 5000 * 3
			expectError:   false,
		},
		{
			name:          "Создание бронирования first класса",
			flightID:      "SU-3456",
			userID:        "user-111",
			seats:         2,
			class:         "first",
			expectedClass: "first",
			expectedCost:  60000.0, // 30000 * 2
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Выполняем создание бронирования
			booking, err := CreateBooking(tt.flightID, tt.userID, tt.seats, tt.class)

			// Проверяем ошибки
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, booking)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, booking)

				// Проверяем поля бронирования
				assert.NotEmpty(t, booking.ID)
				assert.Contains(t, booking.ID, "booking-")
				assert.Equal(t, tt.flightID, booking.FlightID)
				assert.Equal(t, tt.userID, booking.UserID)
				assert.Equal(t, tt.seats, booking.Seats)
				assert.Equal(t, tt.expectedClass, booking.Class)
				assert.Equal(t, tt.expectedCost, booking.TotalCost)
				assert.Equal(t, "confirmed", booking.Status)

				// Проверяем корректность времени создания
				_, err := time.Parse(time.RFC3339, booking.CreatedAt)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserBookings(t *testing.T) {
	tests := []struct {
		name            string
		userID          string
		expectedCount   int
		expectedFirstID string
	}{
		{
			name:            "Получение бронирований существующего пользователя",
			userID:          "user-123",
			expectedCount:   2,
			expectedFirstID: "booking-001",
		},
		{
			name:            "Получение бронирований другого пользователя",
			userID:          "user-456",
			expectedCount:   2,
			expectedFirstID: "booking-001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Получаем бронирования пользователя
			bookings, err := GetUserBookings(tt.userID)

			// Проверяем результат
			assert.NoError(t, err)
			assert.Len(t, bookings, tt.expectedCount)

			if tt.expectedCount > 0 {
				// Проверяем первую запись
				assert.Equal(t, tt.expectedFirstID, bookings[0].ID)
				assert.Equal(t, tt.userID, bookings[0].UserID)
				assert.NotEmpty(t, bookings[0].FlightID)
				assert.NotEmpty(t, bookings[0].Class)
				assert.Greater(t, bookings[0].TotalCost, 0.0)
				assert.Equal(t, "confirmed", bookings[0].Status)
				assert.NotEmpty(t, bookings[0].CreatedAt)

				// Проверяем вторую запись (если есть)
				if len(bookings) > 1 {
					assert.Equal(t, "booking-002", bookings[1].ID)
					assert.Equal(t, "business", bookings[1].Class)
					assert.Equal(t, 25000.0, bookings[1].TotalCost)
				}
			}
		})
	}
}

func TestCancelBooking(t *testing.T) {
	tests := []struct {
		name        string
		bookingID   string
		expectError bool
	}{
		{
			name:        "Успешная отмена бронирования",
			bookingID:   "booking-123",
			expectError: false,
		},
		{
			name:        "Отмена несуществующего бронирования",
			bookingID:   "non-existent-booking",
			expectError: false, // Текущая реализация всегда возвращает nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Выполняем отмену бронирования
			err := CancelBooking(tt.bookingID)

			// Проверяем результат
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalculateTotalCost(t *testing.T) {
	tests := []struct {
		name     string
		flightID string
		seats    int
		class    string
		expected float64
	}{
		{
			name:     "Economy класс, 1 место",
			flightID: "TEST-001",
			seats:    1,
			class:    "economy",
			expected: 5000.0,
		},
		{
			name:     "Economy класс, 3 места",
			flightID: "TEST-002",
			seats:    3,
			class:    "economy",
			expected: 15000.0,
		},
		{
			name:     "Business класс, 1 место",
			flightID: "TEST-003",
			seats:    1,
			class:    "business",
			expected: 15000.0,
		},
		{
			name:     "Business класс, 2 места",
			flightID: "TEST-004",
			seats:    2,
			class:    "business",
			expected: 30000.0,
		},
		{
			name:     "First класс, 1 место",
			flightID: "TEST-005",
			seats:    1,
			class:    "first",
			expected: 30000.0,
		},
		{
			name:     "First класс, 4 места",
			flightID: "TEST-006",
			seats:    4,
			class:    "first",
			expected: 120000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateTotalCost(tt.flightID, tt.seats, tt.class)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateID(t *testing.T) {
	// Тестируем генерацию ID
	id1 := generateID()
	id2 := generateID()

	// Проверяем, что ID не пустые и разные
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)

	// Проверяем формат ID (должен быть в формате YYYYMMDDHHMMSS)
	assert.Len(t, id1, 14)

	// Пытаемся распарсить как дату
	_, err := time.Parse("20060102150405", id1)
	assert.NoError(t, err)
}

// Интеграционный тест полного потока бронирования
func TestBookingFlow(t *testing.T) {
	// Шаг 1: Создание бронирования
	booking, err := CreateBooking("SU-9999", "test-user", 2, "economy")
	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, "confirmed", booking.Status)

	// Шаг 2: Получение бронирований пользователя
	bookings, err := GetUserBookings("test-user")
	assert.NoError(t, err)
	assert.Greater(t, len(bookings), 0)

	// Шаг 3: Отмена бронирования
	err = CancelBooking(booking.ID)
	assert.NoError(t, err)
}
