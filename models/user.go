package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	err := db.Where("id = ?", id).First(&user).Error
	return &user, err
}
