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
// GenerateUsers - генерация пользователей
func (g *SimpleGenerator) GenerateUsers(count int) []models.User {
	users := make([]models.User, count)

	firstNames := []string{"Иван", "Анна", "Петр", "Мария", "Алексей", "Елена", "Сергей", "Ольга"}
	lastNames := []string{"Иванов", "Петрова", "Сидоров", "Смирнова", "Кузнецов", "Попова"}

	for i := 0; i < count; i++ {
		users[i] = models.User{
			Email:        fmt.Sprintf("user%d@example.com", i),
			PasswordHash: "password123", // В реальном приложении должен быть bcrypt/scrypt хэш
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
			Airline:              airline.Name, // Используем название авиакомпании
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

	// Генерируем рейсы
	flights := generator.GenerateFlights(200, savedAirlines, savedAirports)
	// Сохраняем рейсы и получаем их с ID
	savedFlights, err := insertAndGetFlights(db, flights)
	if err != nil {
		log.Fatal("Failed to insert flights:", err)
	}
	fmt.Printf("✅ Inserted %d flights\n", len(savedFlights))

	// Генерируем бронирования
	bookings := generateBookings(30, savedUsers, savedFlights)
	savedBookings, err := insertAndGetBookings(db, bookings)
	if err != nil {
		log.Fatal("Failed to insert bookings:", err)
	}
	fmt.Printf("✅ Inserted %d bookings\n", len(savedBookings))

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

// Генерация бронирований (используем существующую структуру моделей)
// Генерация бронирований (используем существующую структуру моделей)
func generateBookings(count int, users []models.User, flights []models.Flight) []models.Booking {
	var bookings []models.Booking

	// Сначала проверяем, какие типы статусов доступны в нашей модели
	// Создаем переменные с правильными типами
	var statusConfirmed models.BookingStatus
	var statusPending models.BookingStatus
	var statusCancelled models.BookingStatus

	// Пытаемся использовать существующие константы или создаем их
	// В зависимости от того, как определена ваша модель
	if hasBookingStatusConstants() {
		// Если есть константы в пакете models
		statusConfirmed = models.BookingStatusConfirmed
		statusPending = models.BookingStatusPending
		statusCancelled = models.BookingStatusCancelled
	} else {
		// Если нет, создаем строковые значения и конвертируем их
		// Это предполагает, что BookingStatus это type alias для string
		statusConfirmed = models.BookingStatus("confirmed")
		statusPending = models.BookingStatus("pending")
		statusCancelled = models.BookingStatus("cancelled")
	}

	statuses := []models.BookingStatus{statusConfirmed, statusPending, statusCancelled}
	paymentStatuses := []string{"paid", "pending", "failed"}

	for i := 0; i < count; i++ {
		user := users[i%len(users)]
		flight := flights[i%len(flights)]

		// Генерируем случайное количество мест (1-4)
		seats := 1 + (i % 4)

		booking := models.Booking{
			UserID:        user.ID,
			FlightID:      flight.ID,
			Seats:         seats,
			BookingDate:   time.Now().Add(-time.Duration(i*24) * time.Hour),
			Status:        statuses[i%len(statuses)],
			PaymentStatus: paymentStatuses[i%len(paymentStatuses)],
			CreatedAt:     time.Now().Add(-time.Duration(i*24) * time.Hour),
			UpdatedAt:     time.Now().Add(-time.Duration(i*12) * time.Hour),
		}
		bookings = append(bookings, booking)
	}

	return bookings
}

// Вспомогательная функция для проверки наличия констант
func hasBookingStatusConstants() bool {
	// Простая проверка - пытаемся использовать значения
	// Если компиляция пройдет, значит константы существуют
	return false // временно возвращаем false
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

// Очистка существующих данных
func clearExistingData(db *gorm.DB) error {
	// Для PostgreSQL используем TRUNCATE CASCADE
	tables := []string{
		"bookings",
		"flights",
		"airlines",
		"airports",
		"users",
	}

	for _, table := range tables {
		// Сначала проверяем, существует ли таблица
		var tableExists bool
		db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = ?)", table).Scan(&tableExists)

		if tableExists {
			result := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
			if result.Error != nil {
				return fmt.Errorf("failed to truncate table %s: %v", table, result.Error)
			}
			fmt.Printf("✅ Cleared table: %s\n", table)
		} else {
			fmt.Printf("ℹ️ Table %s doesn't exist, skipping\n", table)
		}
	}

	fmt.Println("✅ All tables cleared successfully")
	return nil
}
