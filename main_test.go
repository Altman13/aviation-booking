package main

import (
	"fmt"
	"testing"

	"aviation-booking/models"
	"aviation-booking/testdata"
)

func TestFlightSearch(t *testing.T) {
	generator := testdata.NewTestDataGenerator()

	// Генерируем тестовые данные
	users, _ := generator.GenerateUsers(10)
	airports := generator.GenerateAirports()
	airlines := generator.GenerateAirlines()
	flights := generator.GenerateFlights(50, airlines, airports)

	fmt.Printf("Generated %d users\n", len(users))
	fmt.Printf("Generated %d airports\n", len(airports))
	fmt.Printf("Generated %d flights\n", len(flights))

	// Пример использования в тестах
	t.Run("Search flights from Moscow", func(t *testing.T) {
		moscowFlights := filterFlightsFromMoscow(flights)
		if len(moscowFlights) == 0 {
			t.Error("No flights from Moscow found")
		}
	})
}

func filterFlightsFromMoscow(flights []models.Flight) []models.Flight {
	var result []models.Flight
	for _, flight := range flights {
		if flight.DepartureAirportCode == "SVO" || flight.DepartureAirportCode == "DME" {
			result = append(result, flight)
		}
	}
	return result
}
