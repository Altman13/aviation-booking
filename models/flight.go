package models

import (
	"time"
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
