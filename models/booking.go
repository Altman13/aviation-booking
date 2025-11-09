package models

import (
	"time"
)

type Booking struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id"`
	FlightID      uint      `json:"flight_id"`
	Seats         int       `json:"seats"`
	Class         string    `json:"class"`
	BookingDate   time.Time `json:"booking_date"`
	TotalAmount   float64   `json:"total_amount"`
	Status        string    `json:"status"`         // confirmed, pending, cancelled
	PaymentStatus string    `json:"payment_status"` // paid, pending, refunded
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Связи
	User    User     `json:"user" gorm:"foreignKey:UserID"`
	Flight  Flight   `json:"flight" gorm:"foreignKey:FlightID"`
	Tickets []Ticket `json:"tickets" gorm:"foreignKey:BookingID"`
}
