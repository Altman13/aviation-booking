// models/booking.go
package models

import (
	"time"
)

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusPending   BookingStatus = "pending"
)

type BookingClass string

const (
	BookingClassEconomy  BookingClass = "economy"
	BookingClassBusiness BookingClass = "business"
	BookingClassFirst    BookingClass = "first"
)

type Booking struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	BookingNumber   string        `json:"booking_number" gorm:"uniqueIndex"`
	FlightID        uint          `json:"flight_id"`
	UserID          uint          `json:"user_id"`
	Seats           int           `json:"seats"`
	Class           BookingClass  `json:"class"`
	TotalPrice      float64       `json:"total_price"`
	Status          BookingStatus `json:"status"`
	PaymentStatus   string        `json:"payment_status" gorm:"default:pending"`
	BookingDate     time.Time     `json:"booking_date"`
	CancellationFee float64       `json:"cancellation_fee,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`

	// Связи (опционально)
	Flight Flight `json:"flight,omitempty" gorm:"foreignKey:FlightID"`
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
