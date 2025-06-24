package services

import (
	"aviation-booking/config" // импортируем конфиг
	"aviation-booking/models" // импортируем модель данных
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Структура для получения данных о рейсах
type FlightOffer struct {
	ID        string `json:"id"`
	Departure struct {
		AirportCode string `json:"iataCode"`
		DateTime    string `json:"at"`
	} `json:"departure"`
	Arrival struct {
		AirportCode string `json:"iataCode"`
		DateTime    string `json:"at"`
	} `json:"arrival"`
	Price struct {
		Currency string `json:"currency"`
		Total    string `json:"total"`
	} `json:"price"`
}

// Структура ответа от API Amadeus
type FlightResponse struct {
	Data []FlightOffer `json:"data"`
}

// Функция для получения рейсов из Amadeus API
func GetFlightsFromAmadeus(origin, destination, departureDate string) ([]FlightOffer, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("https://api.amadeus.com/v2/shopping/flight-offers?origin=%s&destination=%s&departureDate=%s", origin, destination, departureDate)

	// Отправляем запрос с авторизацией
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Добавляем заголовок с API ключом
	req.Header.Add("Authorization", "Bearer "+config.AmadeusAPIKey)

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверка на успешный ответ
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get data: %s", resp.Status)
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Десериализация ответа
	var flightResp FlightResponse
	err = json.Unmarshal(body, &flightResp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return flightResp.Data, nil
}

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
