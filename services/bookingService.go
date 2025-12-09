// services/booking_service.go
package services

import (
	"aviation-booking/config"
	"aviation-booking/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookingService struct {
	db            *gorm.DB
	flightService *FlightService
	userService   *UserService
}

// NewBookingService создает новый экземпляр сервиса бронирований
func NewBookingService() *BookingService {
	flightService := NewFlightService()
	userService := NewUserService()

	return &BookingService{
		db:            config.DB,
		flightService: flightService,
		userService:   userService,
	}
}

// WithDB позволяет передать кастомное соединение с БД
func (s *BookingService) WithDB(db *gorm.DB) *BookingService {
	s.db = db
	s.flightService = s.flightService.WithDB(db)
	s.userService = s.userService.WithDB(db)
	return s
}

// BookFlight создает новое бронирование
func (s *BookingService) BookFlight(flightID, userID uint, seats int, class string) (*models.Booking, error) {
	// Проверяем существование рейса
	flight, err := s.flightService.GetFlightByID(flightID)
	if err != nil {
		return nil, errors.New("рейс не найден")
	}

	// Проверяем существование пользователя
	if _, err := s.userService.GetUserByID(userID); err != nil {
		return nil, errors.New("пользователь не найден")
	}

	// Проверяем доступность мест
	available, err := s.flightService.CheckFlightAvailability(flightID, seats)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("недостаточно мест на рейсе")
	}

	// Валидация класса
	bookingClass := models.BookingClass(class)
	if !isValidClass(bookingClass) {
		return nil, errors.New("неверный класс бронирования")
	}

	// Рассчитываем цену
	totalPrice := s.calculatePrice(flight.Price, seats, bookingClass)

	// Генерируем номер бронирования
	bookingNumber := s.generateBookingNumber()

	// Резервируем места
	if err := s.flightService.ReserveSeats(flightID, seats); err != nil {
		return nil, err
	}

	// Создаем бронирование
	booking := &models.Booking{
		BookingNumber: bookingNumber,
		FlightID:      flightID,
		UserID:        userID,
		Seats:         seats,
		Class:         bookingClass,
		TotalPrice:    totalPrice,
		Status:        models.BookingStatusConfirmed,
		PaymentStatus: "pending",
		BookingDate:   time.Now(),
	}

	// Сохраняем в БД
	if err := s.db.Create(booking).Error; err != nil {
		// Если не удалось сохранить, освобождаем места
		s.flightService.ReleaseSeats(flightID, seats)
		return nil, err
	}

	return booking, nil
}

// GetBookingByID получает бронирование по ID
func (s *BookingService) GetBookingByID(bookingID uint) (*models.Booking, error) {
	var booking models.Booking
	if err := s.db.Preload("Flight").Preload("User").
		First(&booking, bookingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("бронирование не найдено")
		}
		return nil, err
	}
	return &booking, nil
}

// GetBookingByNumber получает бронирование по номеру
func (s *BookingService) GetBookingByNumber(bookingNumber string) (*models.Booking, error) {
	var booking models.Booking
	if err := s.db.Preload("Flight").Preload("User").
		Where("booking_number = ?", bookingNumber).
		First(&booking).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("бронирование не найдено")
		}
		return nil, err
	}
	return &booking, nil
}

// GetUserBookings получает все бронирования пользователя
func (s *BookingService) GetUserBookings(userID uint) ([]models.Booking, error) {
	// Проверяем существование пользователя
	if _, err := s.userService.GetUserByID(userID); err != nil {
		return nil, err
	}

	var bookings []models.Booking
	err := s.db.Preload("Flight").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error

	return bookings, err
}

// UpdateBooking обновляет бронирование
func (s *BookingService) UpdateBooking(bookingID uint, seats int, class string) (*models.Booking, error) {
	// Получаем текущее бронирование
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return nil, err
	}

	// Проверяем статус
	if booking.Status == models.BookingStatusCancelled {
		return nil, errors.New("нельзя изменить отмененное бронирование")
	}

	// Если изменяется количество мест
	if seats != booking.Seats {
		seatDifference := seats - booking.Seats

		if seatDifference > 0 {
			// Нужно добавить места
			available, err := s.flightService.CheckFlightAvailability(booking.FlightID, seatDifference)
			if err != nil {
				return nil, err
			}
			if !available {
				return nil, errors.New("недостаточно мест на рейсе")
			}
			if err := s.flightService.ReserveSeats(booking.FlightID, seatDifference); err != nil {
				return nil, err
			}
		} else {
			// Нужно освободить места (seatDifference отрицательный)
			if err := s.flightService.ReleaseSeats(booking.FlightID, -seatDifference); err != nil {
				return nil, err
			}
		}

		booking.Seats = seats
	}

	// Если изменяется класс
	if class != "" {
		bookingClass := models.BookingClass(class)
		if !isValidClass(bookingClass) {
			return nil, errors.New("неверный класс бронирования")
		}
		booking.Class = bookingClass
	}

	// Пересчитываем цену
	flight, err := s.flightService.GetFlightByID(booking.FlightID)
	if err != nil {
		return nil, err
	}
	booking.TotalPrice = s.calculatePrice(flight.Price, booking.Seats, booking.Class)

	// Сохраняем изменения
	if err := s.db.Save(booking).Error; err != nil {
		return nil, err
	}

	return booking, nil
}

// CancelBooking отменяет бронирование
func (s *BookingService) CancelBooking(bookingID uint) error {
	// Получаем бронирование
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return err
	}

	// Проверяем статус
	if booking.Status == models.BookingStatusCancelled {
		return errors.New("бронирование уже отменено")
	}

	// Возвращаем места
	if err := s.flightService.ReleaseSeats(booking.FlightID, booking.Seats); err != nil {
		return err
	}

	// Рассчитываем штраф за отмену
	cancellationFee := s.calculateCancellationFee(booking.TotalPrice, booking.BookingDate)

	// Обновляем статус
	booking.Status = models.BookingStatusCancelled
	booking.CancellationFee = cancellationFee
	booking.PaymentStatus = "refunded" // или "cancelled"

	return s.db.Save(booking).Error
}

// ConfirmPayment подтверждает оплату
func (s *BookingService) ConfirmPayment(bookingID uint) error {
	booking, err := s.GetBookingByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status != models.BookingStatusConfirmed {
		return errors.New("только подтвержденные бронирования могут быть оплачены")
	}

	booking.PaymentStatus = "paid"
	return s.db.Save(booking).Error
}

// SearchBookings ищет бронирования по параметрам
func (s *BookingService) SearchBookings(params BookingSearchParams) ([]models.Booking, error) {
	var bookings []models.Booking

	query := s.db.Preload("Flight").Preload("User")

	// Применяем фильтры
	if params.UserID != 0 {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.FlightID != 0 {
		query = query.Where("flight_id = ?", params.FlightID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.PaymentStatus != "" {
		query = query.Where("payment_status = ?", params.PaymentStatus)
	}
	if params.Class != "" {
		query = query.Where("class = ?", params.Class)
	}
	if !params.DateFrom.IsZero() {
		query = query.Where("booking_date >= ?", params.DateFrom)
	}
	if !params.DateTo.IsZero() {
		query = query.Where("booking_date <= ?", params.DateTo)
	}

	// Пагинация
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// Сортировка
	if params.SortBy != "" {
		query = query.Order(params.SortBy)
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Find(&bookings).Error
	return bookings, err
}

// Вспомогательные методы

func (s *BookingService) calculatePrice(basePrice float64, seats int, class models.BookingClass) float64 {
	// Множители для разных классов
	multipliers := map[models.BookingClass]float64{
		models.BookingClassEconomy:  1.0,
		models.BookingClassBusiness: 2.5,
		models.BookingClassFirst:    4.0,
	}

	multiplier := multipliers[class]
	if multiplier == 0 {
		multiplier = 1.0
	}

	return basePrice * float64(seats) * multiplier
}

func (s *BookingService) calculateCancellationFee(totalPrice float64, bookingDate time.Time) float64 {
	now := time.Now()
	hoursUntilFlight := bookingDate.Sub(now).Hours()

	// Разные штрафы в зависимости от времени до вылета
	if hoursUntilFlight > 48 {
		return totalPrice * 0.1 // 10%
	} else if hoursUntilFlight > 24 {
		return totalPrice * 0.3 // 30%
	} else if hoursUntilFlight > 6 {
		return totalPrice * 0.5 // 50%
	} else {
		return totalPrice * 0.8 // 80%
	}
}

func (s *BookingService) generateBookingNumber() string {
	// Генерация уникального номера бронирования
	timestamp := time.Now().Unix()
	random := time.Now().Nanosecond() % 10000
	return fmt.Sprintf("BK-%d-%04d", timestamp, random)
}

func isValidClass(class models.BookingClass) bool {
	switch class {
	case models.BookingClassEconomy, models.BookingClassBusiness, models.BookingClassFirst:
		return true
	default:
		return false
	}
}

// BookingSearchParams - параметры поиска бронирований
type BookingSearchParams struct {
	UserID        uint
	FlightID      uint
	Status        string
	PaymentStatus string
	Class         string
	DateFrom      time.Time
	DateTo        time.Time
	SortBy        string
	Limit         int
	Offset        int
}
