package services

import (
	"aviation-booking/config" // Импортируем конфигурацию с DB
	"aviation-booking/models"
	"errors"
)

// BookFlight - сервис для создания бронирования
func BookFlight(flight *models.Flight) error {
	// Проверка доступных мест
	if flight.AvailableSeats <= 0 {
		return errors.New("No available seats")
	}

	// Создаем транзакцию
	tx := config.DB.Begin()

	// Убедимся, что транзакция правильно началась
	if tx.Error != nil {
		return tx.Error
	}

	// Создание бронирования
	if err := flight.CreateFlight(tx); err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки
		return err
	}

	// Обновление количества доступных мест
	flight.AvailableSeats--

	// Обновляем рейс в базе данных в рамках транзакции
	if err := tx.Save(flight).Error; err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки
		return err
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
