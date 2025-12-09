// services/flight_service.go
package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type FlightService struct {
	db *gorm.DB
}

func NewFlightService() *FlightService {
	return &FlightService{db: config.DB}
}

func (s *FlightService) WithDB(db *gorm.DB) *FlightService {
	s.db = db
	return s
}

// SearchFlights - поиск рейсов по параметрам
func (s *FlightService) SearchFlights(params FlightSearchParams) ([]models.Flight, error) {
	var flights []models.Flight

	query := s.db.Model(&models.Flight{})

	// Применяем фильтры, если они указаны
	if params.DepartureCity != "" {
		query = query.Where("departure_city = ?", params.DepartureCity)
	}

	if params.ArrivalCity != "" {
		query = query.Where("arrival_city = ?", params.ArrivalCity)
	}

	if params.DepartureTimeFrom != nil {
		query = query.Where("departure_time >= ?", params.DepartureTimeFrom)
	}

	if params.DepartureTimeTo != nil {
		query = query.Where("departure_time <= ?", params.DepartureTimeTo)
	}

	if params.MinPrice > 0 {
		query = query.Where("price >= ?", params.MinPrice)
	}

	if params.MaxPrice > 0 {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	if params.MinAvailableSeats > 0 {
		query = query.Where("available_seats >= ?", params.MinAvailableSeats)
	}

	// Сортировка
	if params.SortBy != "" {
		query = query.Order(params.SortBy)
	}

	// Пагинация
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	err := query.Find(&flights).Error
	if err != nil {
		return nil, err
	}

	return flights, nil
}

// GetAllFlights - получение всех рейсов
func (s *FlightService) GetAllFlights() ([]models.Flight, error) {
	var flights []models.Flight
	err := s.db.Find(&flights).Error
	return flights, err
}

// GetFlightByID - получение рейса по ID
func (s *FlightService) GetFlightByID(flightID uint) (*models.Flight, error) {
	var flight models.Flight
	if err := s.db.First(&flight, flightID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("рейс не найден")
		}
		return nil, err
	}
	return &flight, nil
}

// CreateFlight - создание нового рейса
func (s *FlightService) CreateFlight(flight *models.Flight) error {
	// Валидация
	if err := s.validateFlight(flight); err != nil {
		return err
	}

	return s.db.Create(flight).Error
}

// UpdateFlight - полное обновление рейса
func (s *FlightService) UpdateFlight(flightID uint, flightData *models.Flight) error {
	// Сначала проверяем существование рейса
	existingFlight, err := s.GetFlightByID(flightID)
	if err != nil {
		return err
	}

	// Обновляем поля
	flightData.ID = flightID
	flightData.CreatedAt = existingFlight.CreatedAt // Сохраняем оригинальное время создания

	return s.db.Save(flightData).Error
}

// UpdateFlightSeats - обновление только количества доступных мест
func (s *FlightService) UpdateFlightSeats(flightID uint, seats int) error {
	if seats < 0 {
		return errors.New("количество мест не может быть отрицательным")
	}

	return s.db.Model(&models.Flight{}).
		Where("id = ?", flightID).
		Update("available_seats", seats).Error
}

// DeleteFlight - удаление рейса
func (s *FlightService) DeleteFlight(flightID uint) error {
	result := s.db.Delete(&models.Flight{}, flightID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("рейс не найден")
	}
	return nil
}

// CheckFlightAvailability - проверка доступности мест на рейсе
func (s *FlightService) CheckFlightAvailability(flightID uint, requiredSeats int) (bool, error) {
	flight, err := s.GetFlightByID(flightID)
	if err != nil {
		return false, err
	}

	return flight.AvailableSeats >= requiredSeats, nil
}

// ReserveSeats - бронирование мест на рейсе
func (s *FlightService) ReserveSeats(flightID uint, seats int) error {
	// Проверяем доступность
	available, err := s.CheckFlightAvailability(flightID, seats)
	if err != nil {
		return err
	}
	if !available {
		return errors.New("недостаточно мест на рейсе")
	}

	// Резервируем места (уменьшаем количество доступных)
	return s.db.Model(&models.Flight{}).
		Where("id = ? AND available_seats >= ?", flightID, seats).
		Update("available_seats", gorm.Expr("available_seats - ?", seats)).Error
}

// ReleaseSeats - освобождение мест (отмена бронирования)
func (s *FlightService) ReleaseSeats(flightID uint, seats int) error {
	return s.db.Model(&models.Flight{}).
		Where("id = ?", flightID).
		Update("available_seats", gorm.Expr("available_seats + ?", seats)).Error
}

// GetUpcomingFlights - получение предстоящих рейсов
func (s *FlightService) GetUpcomingFlights(limit int) ([]models.Flight, error) {
	var flights []models.Flight

	err := s.db.Where("departure_time > ?", time.Now()).
		Order("departure_time ASC").
		Limit(limit).
		Find(&flights).Error

	return flights, err
}

// GetFlightsByAirline - получение рейсов по авиакомпании
func (s *FlightService) GetFlightsByAirline(airline string) ([]models.Flight, error) {
	var flights []models.Flight
	err := s.db.Where("airline = ?", airline).Find(&flights).Error
	return flights, err
}

// validateFlight - валидация данных рейса
func (s *FlightService) validateFlight(flight *models.Flight) error {
	if flight.FlightNumber == "" {
		return errors.New("номер рейса обязателен")
	}
	if flight.DepartureCity == "" {
		return errors.New("город вылета обязателен")
	}
	if flight.ArrivalCity == "" {
		return errors.New("город прибытия обязателен")
	}
	if flight.DepartureTime.IsZero() {
		return errors.New("время вылета обязательно")
	}
	if flight.ArrivalTime.IsZero() {
		return errors.New("время прибытия обязательно")
	}
	if flight.DepartureTime.After(flight.ArrivalTime) {
		return errors.New("время вылета не может быть позже времени прибытия")
	}
	if flight.Price <= 0 {
		return errors.New("цена должна быть положительной")
	}
	if flight.AvailableSeats < 0 {
		return errors.New("количество мест не может быть отрицательным")
	}
	return nil
}

// FlightSearchParams - параметры поиска рейсов
type FlightSearchParams struct {
	DepartureCity     string
	ArrivalCity       string
	DepartureTimeFrom *time.Time
	DepartureTimeTo   *time.Time
	MinPrice          float64
	MaxPrice          float64
	MinAvailableSeats int
	Airline           string
	SortBy            string // например: "departure_time ASC", "price DESC"
	Limit             int
	Offset            int
}
