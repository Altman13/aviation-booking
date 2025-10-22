package services

import (
	"strconv"
)

// Flight представляет информацию о рейсе
type Flight struct {
	ID               string  `json:"id"`
	FlightNumber     string  `json:"flight_number"`
	Airline          string  `json:"airline"`
	DepartureAirport string  `json:"departure_airport"`
	ArrivalAirport   string  `json:"arrival_airport"`
	DepartureTime    string  `json:"departure_time"`
	ArrivalTime      string  `json:"arrival_time"`
	Price            float64 `json:"price"`
	Currency         string  `json:"currency"`
	Duration         string  `json:"duration"`
	AvailableSeats   int     `json:"available_seats"`
	Transfers        int     `json:"transfers,omitempty"`
	Aircraft         string  `json:"aircraft,omitempty"`
}

// SearchFlights - поиск рейсов по параметрам
func SearchFlights(origin, destination, date, adults, class string) ([]Flight, error) {
	// Валидация параметров
	if origin == "" || destination == "" || date == "" {
		return nil, nil
	}

	// Конвертируем adults в int
	adultsCount, err := strconv.Atoi(adults)
	if err != nil {
		adultsCount = 1
	}

	// Временная заглушка с тестовыми данными
	flights := []Flight{
		{
			ID:               "SU-1234",
			FlightNumber:     "SU 1234",
			Airline:          "Аэрофлот",
			DepartureAirport: origin,
			ArrivalAirport:   destination,
			DepartureTime:    date + "T10:00:00",
			ArrivalTime:      date + "T11:30:00",
			Price:            calculatePrice(5000, adultsCount, class),
			Currency:         "RUB",
			Duration:         "1h 30m",
			AvailableSeats:   150,
			Transfers:        0,
			Aircraft:         "Airbus A320",
		},
		{
			ID:               "S7-5678",
			FlightNumber:     "S7 5678",
			Airline:          "S7 Airlines",
			DepartureAirport: origin,
			ArrivalAirport:   destination,
			DepartureTime:    date + "T14:20:00",
			ArrivalTime:      date + "T15:50:00",
			Price:            calculatePrice(4500, adultsCount, class),
			Currency:         "RUB",
			Duration:         "1h 30m",
			AvailableSeats:   120,
			Transfers:        0,
			Aircraft:         "Boeing 737",
		},
		{
			ID:               "U6-9012",
			FlightNumber:     "U6 9012",
			Airline:          "Уральские авиалинии",
			DepartureAirport: origin,
			ArrivalAirport:   destination,
			DepartureTime:    date + "T18:45:00",
			ArrivalTime:      date + "T21:15:00",
			Price:            calculatePrice(3800, adultsCount, class),
			Currency:         "RUB",
			Duration:         "2h 30m",
			AvailableSeats:   80,
			Transfers:        1,
			Aircraft:         "Airbus A319",
		},
	}

	return flights, nil
}

// GetFlightByID - получение рейса по ID
func GetFlightByID(flightID string) (*Flight, error) {
	// Временная заглушка - в реальном приложении здесь будет запрос к БД
	flight := &Flight{
		ID:               flightID,
		FlightNumber:     "SU 1234",
		Airline:          "Аэрофлот",
		DepartureAirport: "SVO",
		ArrivalAirport:   "LED",
		DepartureTime:    "2024-01-15T10:00:00",
		ArrivalTime:      "2024-01-15T11:30:00",
		Price:            5000.00,
		Currency:         "RUB",
		Duration:         "1h 30m",
		AvailableSeats:   150,
		Transfers:        0,
		Aircraft:         "Airbus A320",
	}

	return flight, nil
}

// Вспомогательная функция для расчета цены
func calculatePrice(basePrice int, passengers int, class string) float64 {
	multiplier := 1.0

	switch class {
	case "business":
		multiplier = 2.5
	case "first":
		multiplier = 5.0
	}

	return float64(basePrice*passengers) * multiplier
}
