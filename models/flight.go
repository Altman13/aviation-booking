package models

import (
	"gorm.io/gorm" // Новый импорт для GORM
)

type Flight struct {
	ID             uint    `json:"id" gorm:"primaryKey"` // Изменено на `primaryKey` вместо `primary_key`
	DepartureCity  string  `json:"departure_city"`
	ArrivalCity    string  `json:"arrival_city"`
	DepartureTime  string  `json:"departure_time"`
	ArrivalTime    string  `json:"arrival_time"`
	Price          float64 `json:"price"`
	AvailableSeats int     `json:"available_seats"`
}

// CreateFlight создает новый рейс в базе данных
func (f *Flight) CreateFlight(db *gorm.DB) error {
	return db.Create(f).Error
}

// GetFlights получает все рейсы из базы данных
func GetFlights(db *gorm.DB) ([]Flight, error) {
	var flights []Flight
	err := db.Find(&flights).Error
	return flights, err
}
