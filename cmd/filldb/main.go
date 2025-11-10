package main

import (
	"fmt"
	"log"

	"aviation-booking/config"
	"aviation-booking/models"
	"aviation-booking/testdata"

	"gorm.io/gorm"
)

func main() {
	// Инициализируем конфигурацию
	config.Init()
	log.Println("Configuration loaded")

	// Подключаемся к ТЕСТОВОЙ базе данных
	err := config.InitTestDB()
	if err != nil {
		log.Fatal("Failed to connect to TEST database:", err)
	}

	db := config.GetTestDB()
	defer config.CloseTestDB()

	log.Println("Starting test data generation in TEST database...")

	// Очищаем существующие данные
	if err := clearExistingData(db); err != nil {
		log.Fatal("Failed to clear existing data:", err)
	}

	generator := testdata.NewTestDataGenerator()

	// Генерируем тестовые данные
	fmt.Println("Generating test data...")

	users, err := generator.GenerateUsers(20)
	if err != nil {
		log.Fatal("Failed to generate users:", err)
	}
	// Сохраняем пользователей и получаем их с ID
	savedUsers, err := insertAndGetUsers(db, users)
	if err != nil {
		log.Fatal("Failed to insert users:", err)
	}
	fmt.Printf("✅ Inserted %d users\n", len(savedUsers))

	airports := generator.GenerateAirports()
	// Сохраняем аэропорты и получаем их с ID
	savedAirports, err := insertAndGetAirports(db, airports)
	if err != nil {
		log.Fatal("Failed to insert airports:", err)
	}
	fmt.Printf("✅ Inserted %d airports\n", len(savedAirports))

	airlines := generator.GenerateAirlines()
	// Сохраняем авиакомпании и получаем их с ID
	savedAirlines, err := insertAndGetAirlines(db, airlines)
	if err != nil {
		log.Fatal("Failed to insert airlines:", err)
	}
	fmt.Printf("✅ Inserted %d airlines\n", len(savedAirlines))

	flights := generator.GenerateFlights(200, savedAirlines, savedAirports)
	// Сохраняем рейсы и получаем их с ID
	savedFlights, err := insertAndGetFlights(db, flights)
	if err != nil {
		log.Fatal("Failed to insert flights:", err)
	}
	fmt.Printf("✅ Inserted %d flights\n", len(savedFlights))

	// Генерируем и сохраняем бронирования, получая их ID
	bookings := generateBookings(30, savedUsers, savedFlights)
	savedBookings, err := insertAndGetBookings(db, bookings)
	if err != nil {
		log.Fatal("Failed to insert bookings:", err)
	}
	fmt.Printf("✅ Inserted %d bookings\n", len(savedBookings))

	// Теперь генерируем билеты с правильными booking.ID
	tickets := generateTickets(savedBookings, savedFlights, savedUsers)
	if err := insertTickets(db, tickets); err != nil {
		log.Fatal("Failed to insert tickets:", err)
	}
	fmt.Printf("✅ Inserted %d tickets\n", len(tickets))

	fmt.Println("🎉 Test data generated successfully in TEST database!")
	fmt.Println("📊 Database: aviation_booking_test")
	fmt.Println("🚀 You can now run your application with: go run cmd/server/main.go")
}

// Вставка пользователей и возврат с ID
func insertAndGetUsers(db *gorm.DB, users []models.User) ([]models.User, error) {
	var savedUsers []models.User

	for _, user := range users {
		result := db.Create(&user)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to insert user %s: %v", user.Email, result.Error)
		}
		// Добавляем пользователя с заполненным ID
		savedUsers = append(savedUsers, user)
	}

	return savedUsers, nil
}

// Вставка аэропортов и возврат с ID
func insertAndGetAirports(db *gorm.DB, airports []models.Airport) ([]models.Airport, error) {
	var savedAirports []models.Airport

	for _, airport := range airports {
		result := db.Create(&airport)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to insert airport %s: %v", airport.Code, result.Error)
		}
		// Добавляем аэропорт с заполненным ID
		savedAirports = append(savedAirports, airport)
	}

	return savedAirports, nil
}

// Вставка авиакомпаний и возврат с ID
func insertAndGetAirlines(db *gorm.DB, airlines []models.Airline) ([]models.Airline, error) {
	var savedAirlines []models.Airline

	for _, airline := range airlines {
		result := db.Create(&airline)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to insert airline %s: %v", airline.Code, result.Error)
		}
		// Добавляем авиакомпанию с заполненным ID
		savedAirlines = append(savedAirlines, airline)
	}

	return savedAirlines, nil
}

// Вставка рейсов и возврат с ID
func insertAndGetFlights(db *gorm.DB, flights []models.Flight) ([]models.Flight, error) {
	var savedFlights []models.Flight

	for _, flight := range flights {
		result := db.Create(&flight)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to insert flight %s: %v", flight.FlightNumber, result.Error)
		}
		// Добавляем рейс с заполненным ID
		savedFlights = append(savedFlights, flight)
	}

	return savedFlights, nil
}

// Генерация бронирований
func generateBookings(count int, users []models.User, flights []models.Flight) []models.Booking {
	var bookings []models.Booking

	bookingStatuses := []string{"confirmed", "pending", "cancelled"}
	paymentStatuses := []string{"paid", "pending", "failed"}

	for i := 0; i < count; i++ {
		user := users[i%len(users)]
		flight := flights[i%len(flights)]

		// Генерируем случайное количество мест (1-4)
		seats := 1 + (i % 4)
		totalAmount := flight.Price * float64(seats)

		booking := models.Booking{
			UserID:        user.ID,
			FlightID:      flight.ID,
			Seats:         seats,
			BookingDate:   user.CreatedAt,
			TotalAmount:   totalAmount,
			Status:        bookingStatuses[i%len(bookingStatuses)],
			PaymentStatus: paymentStatuses[i%len(paymentStatuses)],
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
		}
		bookings = append(bookings, booking)
	}

	return bookings
}

// Вставка бронирований и возврат сохраненных с ID
func insertAndGetBookings(db *gorm.DB, bookings []models.Booking) ([]models.Booking, error) {
	var savedBookings []models.Booking

	for _, booking := range bookings {
		result := db.Create(&booking)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to insert booking: %v", result.Error)
		}
		// Добавляем сохраненное бронирование с ID
		savedBookings = append(savedBookings, booking)
	}

	return savedBookings, nil
}

// Генерация билетов с правильными booking.ID
func generateTickets(bookings []models.Booking, flights []models.Flight, users []models.User) []models.Ticket {
	var tickets []models.Ticket

	fareConditions := []string{"economy", "business", "first"}

	for _, booking := range bookings {
		// Создаем количество билетов равное количеству мест в бронировании
		for j := 0; j < booking.Seats; j++ {
			// Используем рейс из бронирования
			flightID := booking.FlightID
			var flight models.Flight
			for _, f := range flights {
				if f.ID == flightID {
					flight = f
					break
				}
			}

			// Используем пользователя из бронирования
			userID := booking.UserID
			var user models.User
			for _, u := range users {
				if u.ID == userID {
					user = u
					break
				}
			}

			ticket := models.Ticket{
				BookingID:     booking.ID,
				PassengerName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
				FlightID:      flight.ID,
				SeatNumber:    fmt.Sprintf("%d%c", 10+(j*2), 'A'+j),
				FareCondition: fareConditions[j%len(fareConditions)],
				Price:         flight.Price,
				CreatedAt:     booking.CreatedAt,
				UpdatedAt:     booking.UpdatedAt,
			}
			tickets = append(tickets, ticket)
		}
	}

	return tickets
}

// Очистка существующих данных
func clearExistingData(db *gorm.DB) error {
	// Отключаем проверку внешних ключей
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %v", err)
	}

	tables := []string{"tickets", "bookings", "flights", "airlines", "airports", "users"}

	for _, table := range tables {
		result := db.Exec("DELETE FROM " + table)
		if result.Error != nil {
			// Если таблицы не существует, пропускаем ошибку
			if isTableNotExistsError(result.Error) {
				fmt.Printf("Table %s doesn't exist, skipping\n", table)
				continue
			}
			return fmt.Errorf("failed to clear table %s: %v", table, result.Error)
		}
	}

	// Включаем проверку внешних ключей обратно
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return fmt.Errorf("failed to enable foreign key checks: %v", err)
	}

	fmt.Println("✅ Cleared existing data")
	return nil
}

// Проверка ошибки "таблица не существует"
func isTableNotExistsError(err error) bool {
	if err == nil {
		return false
	}
	errorStr := err.Error()
	return contains(errorStr, "doesn't exist") ||
		contains(errorStr, "Unknown table") ||
		contains(errorStr, "table not found")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Вставка билетов
func insertTickets(db *gorm.DB, tickets []models.Ticket) error {
	for _, ticket := range tickets {
		result := db.Create(&ticket)
		if result.Error != nil {
			return fmt.Errorf("failed to insert ticket: %v", result.Error)
		}
	}
	return nil
}
