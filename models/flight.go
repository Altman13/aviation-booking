package models

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	FlightNumber         string    `json:"flight_number"`
	Airline              string    `json:"airline"`
	DepartureCity        string    `json:"departure_city"`
	DepartureAirportCode string    `json:"departure_airport_code"`
	ArrivalCity          string    `json:"arrival_city"`
	ArrivalAirportCode   string    `json:"arrival_airport_code"`
	DepartureTime        time.Time `json:"departure_time"`
	ArrivalTime          time.Time `json:"arrival_time"`
	Price                float64   `json:"price"`
	AvailableSeats       int       `json:"available_seats"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Create создает новый рейс в базе данных
func (f *Flight) Create(db *gorm.DB) error {
	return db.Create(f).Error
}

// GetAll получает все рейсы из базы данных
func (f *Flight) GetAll(db *gorm.DB) ([]Flight, error) {
	var flights []Flight
	err := db.Find(&flights).Error
	return flights, err
}

// GetByID получает рейс по ID
func (f *Flight) GetByID(db *gorm.DB, id uint) error {
	return db.First(f, id).Error
}

// Update обновляет рейс в базе данных
func (f *Flight) Update(db *gorm.DB) error {
	return db.Save(f).Error
}

// Delete удаляет рейс из базы данных
func (f *Flight) Delete(db *gorm.DB) error {
	return db.Delete(f).Error
}
