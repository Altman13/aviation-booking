package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Create создает нового пользователя в базе данных
func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}

// GetByID получает пользователя по ID
func (u *User) GetByID(db *gorm.DB, id uint) error {
	return db.First(u, id).Error
}

// GetByEmail получает пользователя по email
func (u *User) GetByEmail(db *gorm.DB, email string) error {
	return db.Where("email = ?", email).First(u).Error
}

// Update обновляет пользователя в базе данных
func (u *User) Update(db *gorm.DB) error {
	return db.Save(u).Error
}

// Delete удаляет пользователя из базы данных
func (u *User) Delete(db *gorm.DB) error {
	return db.Delete(u).Error
}

// GetAll получает всех пользователей
func (u *User) GetAll(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}
