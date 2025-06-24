package services

import (
	"aviation-booking/config" // Импортируем конфигурацию с DB
	"aviation-booking/models"
	"errors"
)

// SearchFlights - сервис для поиска рейсов
func SearchFlights(departureCity, arrivalCity string) ([]models.Flight, error) {
	// Если один из городов пуст, возвращаем ошибку
	if departureCity == "" || arrivalCity == "" {
		return nil, errors.New("both departure and arrival cities are required")
	}

	// Ищем рейсы в базе данных
	var flights []models.Flight
	err := config.DB.Where("departure_city = ? AND arrival_city = ?", departureCity, arrivalCity).Find(&flights).Error
	if err != nil {
		return nil, err
	}

	return flights, nil
}
