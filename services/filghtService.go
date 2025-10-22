package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
)

// SearchFlights - поиск рейсов по параметрам
func SearchFlights(departureCity, arrivalCity string) ([]models.Flight, error) {
	db := config.DB

	var flights []models.Flight

	query := db.Model(&models.Flight{})

	if departureCity != "" {
		query = query.Where("departure_city = ?", departureCity)
	}

	if arrivalCity != "" {
		query = query.Where("arrival_city = ?", arrivalCity)
	}

	err := query.Find(&flights).Error
	if err != nil {
		return nil, err
	}

	return flights, nil
}

// GetAllFlights - получение всех рейсов
func GetAllFlights() ([]models.Flight, error) {
	db := config.DB

	var flights []models.Flight
	err := db.Find(&flights).Error
	return flights, err
}

// GetFlightByID - получение рейса по ID
func GetFlightByID(flightID uint) (*models.Flight, error) {
	db := config.DB

	var flight models.Flight
	if err := db.First(&flight, flightID).Error; err != nil {
		return nil, err
	}
	return &flight, nil
}

// CreateFlight - создание нового рейса
func CreateFlight(flight *models.Flight) error {
	db := config.DB
	return db.Create(flight).Error
}

// UpdateFlightSeats - обновление количества доступных мест
func UpdateFlightSeats(flightID uint, seats int) error {
	db := config.DB
	return db.Model(&models.Flight{}).Where("id = ?", flightID).Update("available_seats", seats).Error
}
