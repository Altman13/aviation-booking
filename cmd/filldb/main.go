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

	// Подключаемся к тестовой базе данных
	err := config.InitTestDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := config.GetTestDB()
	defer config.CloseTestDB()

	generator := testdata.NewTestDataGenerator()

	fmt.Println("Generating test data...")

	// Очищаем существующие данные
	if err := clearExistingData(db); err != nil {
		log.Fatal("Failed to clear existing data:", err)
	}

	// Генерируем и вставляем данные
	users, err := generator.GenerateUsers(100)
	if err != nil {
		log.Fatal("Failed to generate users:", err)
	}
	if err := insertUsers(db, users); err != nil {
		log.Fatal("Failed to insert users:", err)
	}
	fmt.Printf("Inserted %d users\n", len(users))

	airports := generator.GenerateAirports()
	if err := insertAirports(db, airports); err != nil {
		log.Fatal("Failed to insert airports:", err)
	}
	fmt.Printf("Inserted %d airports\n", len(airports))

	airlines := generator.GenerateAirlines()
	if err := insertAirlines(db, airlines); err != nil {
		log.Fatal("Failed to insert airlines:", err)
	}
	fmt.Printf("Inserted %d airlines\n", len(airlines))

	flights := generator.GenerateFlights(1000, airlines, airports)
	if err := insertFlights(db, flights); err != nil {
		log.Fatal("Failed to insert flights:", err)
	}
	fmt.Printf("Inserted %d flights\n", len(flights))

	// Генерируем бронирования и билеты
	bookings, tickets := generateBookingsAndTickets(50, users, flights)
	if err := insertBookings(db, bookings); err != nil {
		log.Fatal("Failed to insert bookings:", err)
	}
	fmt.Printf("Inserted %d bookings\n", len(bookings))

	if err := insertTickets(db, tickets); err != nil {
		log.Fatal("Failed to insert tickets:", err)
	}
	fmt.Printf("Inserted %d tickets\n", len(tickets))

	fmt.Println("Test data generated successfully!")
}

// Очистка существующих данных
func clearExistingData(db *gorm.DB) error {
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

	fmt.Println("Cleared existing data")
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

// Вставка пользователей
func insertUsers(db *gorm.DB, users []models.User) error {
	for _, user := range users {
		result := db.Create(&user)
		if result.Error != nil {
			return fmt.Errorf("failed to insert user %s: %v", user.Email, result.Error)
		}
	}
	return nil
}

// Вставка аэропортов
func insertAirports(db *gorm.DB, airports []models.Airport) error {
	for _, airport := range airports {
		result := db.Create(&airport)
		if result.Error != nil {
			return fmt.Errorf("failed to insert airport %s: %v", airport.Code, result.Error)
		}
	}
	return nil
}

// Вставка авиакомпаний
func insertAirlines(db *gorm.DB, airlines []models.Airline) error {
	for _, airline := range airlines {
		result := db.Create(&airline)
		if result.Error != nil {
			return fmt.Errorf("failed to insert airline %s: %v", airline.Code, result.Error)
		}
	}
	return nil
}

// Вставка рейсов
func insertFlights(db *gorm.DB, flights []models.Flight) error {
	for _, flight := range flights {
		result := db.Create(&flight)
		if result.Error != nil {
			return fmt.Errorf("failed to insert flight %s: %v", flight.FlightNumber, result.Error)
		}
	}
	return nil
}

// Генерация бронирований и билетов
// Генерация бронирований и билетов
func generateBookingsAndTickets(count int, users []models.User, flights []models.Flight) ([]models.Booking, []models.Ticket) {
	var bookings []models.Booking
	var tickets []models.Ticket

	bookingStatuses := []string{"confirmed", "pending", "cancelled"}
	fareConditions := []string{"economy", "business", "first"}

	for i := 0; i < count; i++ {
		// Создаем бронирование - ID будет автоинкрементным
		booking := models.Booking{
			UserID:        users[i%len(users)].ID,
			BookingDate:   users[i%len(users)].CreatedAt,
			TotalAmount:   float64(10000 + (i * 1000)),
			Status:        bookingStatuses[i%len(bookingStatuses)],
			PaymentStatus: "paid",
			CreatedAt:     users[i%len(users)].CreatedAt,
			UpdatedAt:     users[i%len(users)].UpdatedAt,
		}
		bookings = append(bookings, booking)

		// Создаем 1-3 билета для каждого бронирования
		ticketCount := 1 + (i % 3)
		for j := 0; j < ticketCount; j++ {
			flight := flights[(i+j)%len(flights)]
			ticket := models.Ticket{
				PassengerName: fmt.Sprintf("%s %s", users[i%len(users)].FirstName, users[i%len(users)].LastName),
				FlightID:      flight.ID,
				SeatNumber:    fmt.Sprintf("%d%c", 10+(j*2), 'A'+j),
				FareCondition: fareConditions[j%len(fareConditions)],
				Price:         flight.Price * float64(j+1) * 0.8,
				CreatedAt:     booking.CreatedAt,
				UpdatedAt:     booking.UpdatedAt,
			}
			tickets = append(tickets, ticket)
		}
	}

	return bookings, tickets
}

// Вставка бронирований
func insertBookings(db *gorm.DB, bookings []models.Booking) error {
	for _, booking := range bookings {
		result := db.Create(&booking)
		if result.Error != nil {
			return fmt.Errorf("failed to insert booking %s: %v", booking.BookingID, result.Error)
		}
	}
	return nil
}

// Вставка билетов
func insertTickets(db *gorm.DB, tickets []models.Ticket) error {
	for _, ticket := range tickets {
		result := db.Create(&ticket)
		if result.Error != nil {
			return fmt.Errorf("failed to insert ticket for booking %s: %v", ticket.BookingID, result.Error)
		}
	}
	return nil
}
