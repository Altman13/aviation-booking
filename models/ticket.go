package models

import "time"

type Ticket struct {
	ID            uint      `json:"id"`
	BookingID     uint      `json:"booking_id"`
	PassengerName string    `json:"passenger_name"`
	FlightID      uint      `json:"flight_id"`
	SeatNumber    string    `json:"seat_number"`
	FareCondition string    `json:"fare_condition"`
	Price         float64   `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
