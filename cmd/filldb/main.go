package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"aviation-booking/config"
	"aviation-booking/models"

	"gorm.io/gorm"
)

// SimpleGenerator - упрощенный генератор тестовых данных
type SimpleGenerator struct {
	random *rand.Rand
}

func NewSimpleGenerator() *SimpleGenerator {
	return &SimpleGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateUsers - генерация пользователей
func (g *SimpleGenerator) GenerateUsers(count int) []models.User {
	users := make([]models.User, count)

	firstNames := []string{"Иван", "Анна", "Петр", "Мария", "Алексей", "Елена", "Сергей", "Ольга"}
	lastNames := []string{"Иванов", "Петрова", "Сидоров", "Смирнова", "Кузнецов", "Попова"}

	for i := 0; i < count; i++ {
		users[i] = models.User{
			Email:        fmt.Sprintf("user%d@example.com", i),
			PasswordHash: "password123",
			FirstName:    firstNames[i%len(firstNames)],
			LastName:     lastNames[i%len(lastNames)],
			CreatedAt:    time.Now().Add(-time.Duration(g.random.Intn(365)) * 24 * time.Hour),
			UpdatedAt:    time.Now(),
		}
	}

	return users
}

// GenerateAirports - генерация аэропортов
func (g *SimpleGenerator) GenerateAirports() []models.Airport {
	return []models.Airport{
		{Code: "SVO", Name: "Шереметьево", City: "Москва", Country: "Россия"},
		{Code: "DME", Name: "Домодедово", City: "Москва", Country: "Россия"},
		{Code: "LED", Name: "Пулково", City: "Санкт-Петербург", Country: "Россия"},
		{Code: "AER", Name: "Сочи", City: "Сочи", Country: "Россия"},
		{Code: "KRR", Name: "Краснодар", City: "Краснодар", Country: "Россия"},
		{Code: "KHV", Name: "Хабаровск", City: "Хабаровск", Country: "Россия"},
	}
}

// GenerateAirlines - генерация авиакомпаний
func (g *SimpleGenerator) GenerateAirlines() []models.Airline {
	return []models.Airline{
		{Code: "SU", Name: "Аэрофлот", Country: "Россия"},
		{Code: "S7", Name: "S7 Airlines", Country: "Россия"},
		{Code: "U6", Name: "Уральские авиалинии", Country: "Россия"},
		{Code: "FV", Name: "Россия", Country: "Россия"},
	}
}

// GenerateFlights - генерация рейсов
func (g *SimpleGenerator) GenerateFlights(count int, airlines []models.Airline, airports []models.Airport) []models.Flight {
	flights := make([]models.Flight, count)

	for i := 0; i < count; i++ {
		airline := airlines[g.random.Intn(len(airlines))]
		depAirport := airports[g.random.Intn(len(airports))]
		arrAirport := airports[g.random.Intn(len(airports))]

		// Убедимся что аэропорты разные
		for depAirport.Code == arrAirport.Code {
			arrAirport = airports[g.random.Intn(len(airports))]
		}

		departureTime := time.Now().Add(time.Duration(g.random.Intn(30*24)) * time.Hour)
		arrivalTime := departureTime.Add(time.Duration(1+g.random.Intn(8)) * time.Hour)

		flights[i] = models.Flight{
			FlightNumber:         fmt.Sprintf("%s%d", airline.Code, 1000+g.random.Intn(9000)),
			Airline:              airline.Name,
			DepartureCity:        depAirport.City,
			DepartureAirportCode: depAirport.Code,
			ArrivalCity:          arrAirport.City,
			ArrivalAirportCode:   arrAirport.Code,
			DepartureTime:        departureTime,
			ArrivalTime:          arrivalTime,
			Price:                2000.0 + float64(g.random.Intn(28000)),
			AvailableSeats:       50 + g.random.Intn(250),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
	}

	return flights
}

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

	generator := NewSimpleGenerator()

	// Генерируем тестовые данные
	fmt.Println("Generating test data...")

	users := generator.GenerateUsers(20)
	savedUsers, err := insertAndGetUsers(db, users)
	if err != nil {
		log.Fatal("Failed to insert users:", err)
	}
	fmt.Printf("✅ Inserted %d users\n", len(savedUsers))

	airports := generator.GenerateAirports()
	savedAirports, err := insertAndGetAirports(db, airports)
	if err != nil {
		log.Fatal("Failed to insert airports:", err)
	}
	fmt.Printf("✅ Inserted %d airports\n", len(savedAirports))

	airlines := generator.GenerateAirlines()
	savedAirlines, err := insertAndGetAirlines(db, airlines)
	if err != nil {
		log.Fatal("Failed to insert airlines:", err)
	}
	fmt.Printf("✅ Inserted %d airlines\n", len(savedAirlines))

	flights := generator.GenerateFlights(200, savedAirlines, savedAirports)
	savedFlights, err := insertAndGetFlights(db, flights)
	if err != nil {
		log.Fatal("Failed to insert flights:", err)
	}
	fmt.Printf("✅ Inserted %d flights\n", len(savedFlights))

	bookings, err := generateAndInsertBookings(db, 30, savedUsers, savedFlights, generator)
	if err != nil {
		log.Fatal("Failed to insert bookings:", err)
	}
	fmt.Printf("✅ Inserted %d bookings\n", len(bookings))

	tickets, err := generateAndInsertTickets(db, bookings, savedFlights, generator)
	if err != nil {
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
		savedFlights = append(savedFlights, flight)
	}

	return savedFlights, nil
}

// Генерация и вставка бронирований
func generateAndInsertBookings(db *gorm.DB, count int, users []models.User, flights []models.Flight, generator *SimpleGenerator) ([]models.Booking, error) {
	var savedBookings []models.Booking

	// Генерируем уникальные номера бронирований
	usedBookingNumbers := make(map[string]bool)

	// Используем константы из models пакета
	bookingClasses := []models.BookingClass{
		models.BookingClassEconomy,
		models.BookingClassBusiness,
		models.BookingClassFirst,
	}

	bookingStatuses := []models.BookingStatus{
		models.BookingStatusConfirmed,
		models.BookingStatusPending,
		models.BookingStatusCancelled,
	}

	paymentStatuses := []string{"paid", "pending", "failed"}

	for i := 0; i < count; i++ {
		// Генерируем уникальный номер бронирования
		var bookingNumber string
		for {
			bookingNumber = fmt.Sprintf("BK-%06d", generator.random.Intn(999999))
			if !usedBookingNumbers[bookingNumber] {
				usedBookingNumbers[bookingNumber] = true
				break
			}
		}

		user := users[i%len(users)]
		flight := flights[i%len(flights)]

		// Генерируем случайное количество мест (1-4)
		seats := 1 + (i % 4)

		// Расчет общей цены
		totalPrice := flight.Price * float64(seats)

		// Выбираем случайный класс из доступных
		class := bookingClasses[generator.random.Intn(len(bookingClasses))]

		// Выбираем случайный статус из доступных
		status := bookingStatuses[generator.random.Intn(len(bookingStatuses))]

		// Выбираем случайный статус оплаты
		paymentStatus := paymentStatuses[generator.random.Intn(len(paymentStatuses))]

		booking := models.Booking{
			BookingNumber: bookingNumber,
			UserID:        user.ID,
			FlightID:      flight.ID,
			Seats:         seats,
			Class:         class,
			TotalPrice:    totalPrice,
			Status:        status,
			PaymentStatus: paymentStatus,
			BookingDate:   time.Now().Add(-time.Duration(i*24) * time.Hour),
			CreatedAt:     time.Now().Add(-time.Duration(i*24) * time.Hour),
			UpdatedAt:     time.Now().Add(-time.Duration(i*12) * time.Hour),
		}

		result := db.Create(&booking)
		if result.Error != nil {
			return nil, fmt.Errorf("failed to insert booking: %v", result.Error)
		}
		savedBookings = append(savedBookings, booking)
	}

	return savedBookings, nil
}

// Генерация и вставка билетов
func generateAndInsertTickets(db *gorm.DB, bookings []models.Booking, flights []models.Flight, generator *SimpleGenerator) ([]models.Ticket, error) {
	var savedTickets []models.Ticket

	firstNames := []string{"Иван", "Анна", "Петр", "Мария", "Алексей", "Елена", "Сергей", "Ольга"}
	lastNames := []string{"Иванов", "Петрова", "Сидоров", "Смирнова", "Кузнецов", "Попова"}
	fareConditions := []string{"Economy", "Business", "First"}

	for _, booking := range bookings {
		// Для каждого места в бронировании создаем отдельный билет
		for seatNum := 1; seatNum <= booking.Seats; seatNum++ {
			// Генерируем имя пассажира
			passengerName := fmt.Sprintf("%s %s",
				firstNames[generator.random.Intn(len(firstNames))],
				lastNames[generator.random.Intn(len(lastNames))],
			)

			// Генерируем номер места
			seatNumber := fmt.Sprintf("%d%c",
				1+generator.random.Intn(30),
				'A'+rune(generator.random.Intn(6)),
			)

			// Определяем условия тарифа на основе класса бронирования
			fareCondition := fareConditions[0] // По умолчанию Economy
			switch booking.Class {
			case models.BookingClassBusiness:
				fareCondition = fareConditions[1]
			case models.BookingClassFirst:
				fareCondition = fareConditions[2]
			}

			// Расчет цены билета (общая цена / количество мест)
			ticketPrice := booking.TotalPrice / float64(booking.Seats)

			ticket := models.Ticket{
				BookingID:     booking.ID,
				PassengerName: passengerName,
				FlightID:      booking.FlightID,
				SeatNumber:    seatNumber,
				FareCondition: fareCondition,
				Price:         ticketPrice,
				CreatedAt:     booking.CreatedAt,
				UpdatedAt:     booking.UpdatedAt,
			}

			result := db.Create(&ticket)
			if result.Error != nil {
				return nil, fmt.Errorf("failed to insert ticket: %v", result.Error)
			}
			savedTickets = append(savedTickets, ticket)
		}
	}

	return savedTickets, nil
}

// Очистка существующих данных
func clearExistingData(db *gorm.DB) error {
	// Для MariaDB/MySQL
	tables := []string{
		"tickets",
		"bookings",
		"flights",
		"airlines",
		"airports",
		"users",
	}

	for _, table := range tables {
		// Сначала проверяем, существует ли таблица (исправленный синтаксис для MariaDB)
		var tableExists bool
		// Используем 1 вместо * для MariaDB
		db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?)", table).Scan(&tableExists)

		if tableExists {
			// Отключаем проверку внешних ключей для MySQL/MariaDB
			db.Exec("SET FOREIGN_KEY_CHECKS = 0")

			// Используем DELETE или TRUNCATE для MySQL/MariaDB
			result := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
			if result.Error != nil {
				db.Exec("SET FOREIGN_KEY_CHECKS = 1")
				return fmt.Errorf("failed to clear table %s: %v", table, result.Error)
			}

			// Сбрасываем автоинкремент для таблиц с ID
			if table != "users" && table != "bookings" && table != "tickets" {
				db.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table))
			}

			db.Exec("SET FOREIGN_KEY_CHECKS = 1")
			fmt.Printf("✅ Cleared table: %s\n", table)
		} else {
			fmt.Printf("ℹ️ Table %s doesn't exist, skipping\n", table)
		}
	}

	fmt.Println("✅ All tables cleared successfully")
	return nil
}
