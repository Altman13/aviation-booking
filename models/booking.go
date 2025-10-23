package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	FlightID      uint      `json:"flight_id" binding:"required"`
	UserID        uint      `json:"user_id" binding:"required"`
	Seats         int       `json:"seats" binding:"required,min=1"`
	Class         string    `json:"class" binding:"required"`
	Status        string    `json:"status" gorm:"default:confirmed"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	BookingID     string    `json:"booking_id" faker:"-"`
	BookingDate   time.Time `json:"booking_date"`
	TotalAmount   float64   `json:"total_amount"`
	PaymentStatus string    `json:"payment_status"`

	// Связи (опционально, для eager loading)
	Flight Flight `json:"flight" gorm:"foreignKey:FlightID"`
	User   User   `json:"user" gorm:"foreignKey:UserID"`
}

// Create создает новое бронирование в базе данных
func (b *Booking) Create(db *gorm.DB) error {
	return db.Create(b).Error
}

// GetByID получает бронирование по ID
func (b *Booking) GetByID(db *gorm.DB, id uint) error {
	return db.Preload("Flight").Preload("User").First(b, id).Error
}

// GetByUserID получает все бронирования пользователя
func (b *Booking) GetByUserID(db *gorm.DB, userID uint) ([]Booking, error) {
	var bookings []Booking
	err := db.Preload("Flight").Where("user_id = ?", userID).Find(&bookings).Error
	return bookings, err
}

// Update обновляет бронирование в базе данных
func (b *Booking) Update(db *gorm.DB) error {
	return db.Save(b).Error
}

// Cancel отменяет бронирование
func (b *Booking) Cancel(db *gorm.DB) error {
	b.Status = "cancelled"
	return db.Save(b).Error
}

// Delete удаляет бронирование из базы данных
func (b *Booking) Delete(db *gorm.DB) error {
	return db.Delete(b).Error
}
