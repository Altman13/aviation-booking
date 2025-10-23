package testdata

import (
	"fmt"
	"math/rand"
	"time"

	"aviation-booking/models"

	"github.com/bxcodec/faker/v3"
)

type TestDataGenerator struct{}

func NewTestDataGenerator() *TestDataGenerator {
	rand.Seed(time.Now().UnixNano())
	return &TestDataGenerator{}
}

// GenerateUsers генерирует тестовых пользователей
func (g *TestDataGenerator) GenerateUsers(count int) ([]models.User, error) {
	var users []models.User

	for i := 0; i < count; i++ {
		var user models.User
		err := faker.FakeData(&user)
		if err != nil {
			return nil, err
		}

		user.ID = uint(i + 1)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		users = append(users, user)
	}

	return users, nil
}

// GenerateAirports генерирует тестовые аэропорты
func (g *TestDataGenerator) GenerateAirports() []models.Airport {
	return []models.Airport{
		{Code: "SVO", Name: "Шереметьево", City: "Москва", Country: "RU"},
		{Code: "DME", Name: "Домодедово", City: "Москва", Country: "RU"},
		{Code: "VKO", Name: "Внуково", City: "Москва", Country: "RU"},
		{Code: "LED", Name: "Пулково", City: "Санкт-Петербург", Country: "RU"},
		{Code: "IST", Name: "Стамбул", City: "Стамбул", Country: "TR"},
		{Code: "FRA", Name: "Франкфурт", City: "Франкфурт", Country: "DE"},
		{Code: "CDG", Name: "Шарль-де-Голль", City: "Париж", Country: "FR"},
		{Code: "LHR", Name: "Хитроу", City: "Лондон", Country: "GB"},
	}
}

// GenerateAirlines генерирует тестовые авиакомпании
func (g *TestDataGenerator) GenerateAirlines() []models.Airline {
	return []models.Airline{
		{Code: "SU", Name: "Аэрофлот", Country: "RU"},
		{Code: "S7", Name: "S7 Airlines", Country: "RU"},
		{Code: "U6", Name: "Уральские авиалинии", Country: "RU"},
		{Code: "TK", Name: "Turkish Airlines", Country: "TR"},
		{Code: "LH", Name: "Lufthansa", Country: "DE"},
		{Code: "AF", Name: "Air France", Country: "FR"},
		{Code: "BA", Name: "British Airways", Country: "GB"},
	}
}

// GenerateFlights генерирует тестовые рейсы
func (g *TestDataGenerator) GenerateFlights(count int, airlines []models.Airline, airports []models.Airport) []models.Flight {
	var flights []models.Flight

	flightNumbers := []string{"1234", "567", "789", "901", "345", "678", "234", "890"}

	for i := 0; i < count; i++ {
		departureAirport := airports[rand.Intn(len(airports))]
		arrivalAirport := airports[rand.Intn(len(airports))]

		// Гарантируем, что аэропорты разные
		for departureAirport.Code == arrivalAirport.Code {
			arrivalAirport = airports[rand.Intn(len(airports))]
		}

		airline := airlines[rand.Intn(len(airlines))]
		departureTime := time.Now().Add(time.Duration(rand.Intn(30)) * 24 * time.Hour)
		flightDuration := time.Duration(rand.Intn(8)+1) * time.Hour
		arrivalTime := departureTime.Add(flightDuration)

		flight := models.Flight{
			ID:                   uint(i + 1),
			FlightNumber:         fmt.Sprintf("%s %s", airline.Code, flightNumbers[rand.Intn(len(flightNumbers))]),
			Airline:              airline.Name,
			DepartureCity:        departureAirport.City,
			DepartureAirportCode: departureAirport.Code,
			ArrivalCity:          arrivalAirport.City,
			ArrivalAirportCode:   arrivalAirport.Code,
			DepartureTime:        departureTime,
			ArrivalTime:          arrivalTime,
			Price:                float64(rand.Intn(20000) + 5000),
			AvailableSeats:       rand.Intn(200) + 50,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		flights = append(flights, flight)
	}

	return flights
}
