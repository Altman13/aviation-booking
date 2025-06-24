package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFlightsFromAmadeus(t *testing.T) {
	// Тестируем поиск рейсов с известными параметрами
	origin := "LON"
	destination := "NYC"
	departureDate := "2023-12-01"

	// Выполняем запрос
	flights, err := GetFlightsFromAmadeus(origin, destination, departureDate)

	// Проверяем, что ошибок нет
	assert.NoError(t, err)

	// Проверяем, что возвращен хотя бы один рейс
	assert.Greater(t, len(flights), 0)

	// Проверяем структуру данных
	assert.NotEmpty(t, flights[0].ID)
	assert.NotEmpty(t, flights[0].Price.Total)
	assert.NotEmpty(t, flights[0].Departure.AirportCode)
}
