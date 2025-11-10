package testdata

import (
	"aviation-booking/models"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TestDataGenerator struct {
	random *rand.Rand
}

func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateUsers - генерация реалистичных пользователей
func (g *TestDataGenerator) GenerateUsers(count int) ([]models.User, error) {
	users := make([]models.User, count)

	// Мужские имена
	maleFirstNames := []string{
		"Александр", "Алексей", "Андрей", "Антон", "Артем", "Борис", "Вадим", "Валентин",
		"Валерий", "Виктор", "Виталий", "Владимир", "Вячеслав", "Геннадий", "Георгий",
		"Денис", "Дмитрий", "Евгений", "Иван", "Игорь", "Кирилл", "Константин", "Максим",
		"Михаил", "Никита", "Николай", "Олег", "Павел", "Петр", "Роман", "Сергей", "Станислав",
		"Степан", "Юрий", "Ярослав",
	}

	// Женские имена
	femaleFirstNames := []string{
		"Александра", "Алена", "Алина", "Алла", "Анастасия", "Анна", "Валентина", "Валерия",
		"Вера", "Виктория", "Галина", "Дарья", "Евгения", "Екатерина", "Елена", "Ирина",
		"Ксения", "Лариса", "Любовь", "Людмила", "Марина", "Мария", "Надежда", "Наталья",
		"Оксана", "Ольга", "Светлана", "Татьяна", "Юлия", "Яна",
	}

	// Мужские фамилии
	maleLastNames := []string{
		"Иванов", "Петров", "Сидоров", "Смирнов", "Кузнецов", "Попов", "Лебедев", "Козлов",
		"Новиков", "Морозов", "Волков", "Соловьев", "Васильев", "Зайцев", "Павлов",
		"Семенов", "Голубев", "Виноградов", "Богданов", "Воробьев", "Федоров", "Михайлов",
		"Беляев", "Тарасов", "Белов", "Комаров", "Орлов", "Киселев", "Макаров", "Андреев",
		"Ковалев", "Ильин", "Гусев", "Титов", "Кузьмин", "Кудрявцев", "Баранов", "Куликов",
		"Алексеев", "Степанов", "Яковлев", "Сорокин", "Сергеев", "Романов", "Захаров",
		"Борисов", "Королев", "Герасимов", "Пономарев", "Григорьев",
	}

	// Женские фамилии
	femaleLastNames := []string{
		"Иванова", "Петрова", "Сидорова", "Смирнова", "Кузнецова", "Попова", "Лебедева",
		"Козлова", "Новикова", "Морозова", "Волкова", "Соловьева", "Васильева", "Зайцева",
		"Павлова", "Семенова", "Голубева", "Виноградова", "Богданова", "Воробьева",
		"Федорова", "Михайлова", "Беляева", "Тарасова", "Белова", "Комарова", "Орлова",
		"Киселева", "Макарова", "Андреева", "Ковалева", "Ильина", "Гусева", "Титова",
		"Кузьмина", "Кудрявцева", "Баранова", "Куликова", "Алексеева", "Степанова",
		"Яковлева", "Сорокина", "Сергеева", "Романова", "Захарова", "Борисова", "Королева",
		"Герасимова", "Пономарева", "Григорьева",
	}

	domains := []string{"gmail.com", "yandex.ru", "mail.ru", "rambler.ru", "outlook.com", "yahoo.com"}

	// Коды операторов для России
	operators := []string{
		"901", "902", "903", "904", "905", "906", "908", "909",
		"910", "911", "912", "913", "914", "915", "916", "917", "918", "919",
		"920", "921", "922", "923", "924", "925", "926", "927", "928", "929",
		"930", "931", "932", "933", "934", "936", "937", "938", "939",
		"950", "951", "952", "953", "954", "955", "956", "958", "959",
		"960", "961", "962", "963", "964", "965", "966", "967", "968", "969",
		"980", "981", "982", "983", "984", "985", "986", "987", "988", "989",
	}

	for i := 0; i < count; i++ {
		// Случайно выбираем пол (50/50)
		isMale := g.random.Intn(2) == 0

		var firstName, lastName string
		if isMale {
			firstName = maleFirstNames[g.random.Intn(len(maleFirstNames))]
			lastName = maleLastNames[g.random.Intn(len(maleLastNames))]
		} else {
			firstName = femaleFirstNames[g.random.Intn(len(femaleFirstNames))]
			lastName = femaleLastNames[g.random.Intn(len(femaleLastNames))]
		}

		domain := domains[g.random.Intn(len(domains))]

		// Создаем email на основе имени и фамилии
		email := fmt.Sprintf("%s.%s%d@%s",
			transliterate(firstName),
			transliterate(lastName),
			g.random.Intn(1000),
			domain)

		// Генерируем телефонный номер в российском формате
		operator := operators[g.random.Intn(len(operators))]
		number := fmt.Sprintf("%07d", g.random.Intn(10000000))
		phoneNumber := fmt.Sprintf("+7%s%s", operator, number)

		// Генерируем дату рождения (от 18 до 70 лет)
		years := 18 + g.random.Intn(52) // 18-70 лет
		days := g.random.Intn(365)
		dateOfBirth := time.Now().AddDate(-years, 0, -days)

		// Нормализуем имена (первая буква заглавная, остальные строчные)
		caser := cases.Title(language.Russian)
		normalizedFirstName := caser.String(firstName)
		normalizedLastName := caser.String(lastName)

		users[i] = models.User{
			FirstName:   normalizedFirstName,
			LastName:    normalizedLastName,
			Email:       email,
			PhoneNumber: phoneNumber,
			DateOfBirth: dateOfBirth,
			CreatedAt:   time.Now().Add(-time.Duration(g.random.Intn(365)) * 24 * time.Hour),
			UpdatedAt:   time.Now(),
		}
	}

	return users, nil
}

// transliterate - транслитерация кириллицы в латиницу для email
func transliterate(text string) string {
	translitMap := map[rune]string{
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e", 'ж': "zh",
		'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n", 'о': "o",
		'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
		'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu",
		'я': "ya",
		'А': "a", 'Б': "b", 'В': "v", 'Г': "g", 'Д': "d", 'Е': "e", 'Ё': "e", 'Ж': "zh",
		'З': "z", 'И': "i", 'Й': "y", 'К': "k", 'Л': "l", 'М': "m", 'Н': "n", 'О': "o",
		'П': "p", 'Р': "r", 'С': "s", 'Т': "t", 'У': "u", 'Ф': "f", 'Х': "kh", 'Ц': "ts",
		'Ч': "ch", 'Ш': "sh", 'Щ': "shch", 'Ъ': "", 'Ы': "y", 'Ь': "", 'Э': "e", 'Ю': "yu",
		'Я': "ya",
	}

	var result string
	for _, char := range text {
		if engChar, exists := translitMap[char]; exists {
			result += engChar
		} else {
			result += string(char)
		}
	}
	return result
}

// GenerateAirports - генерация аэропортов
func (g *TestDataGenerator) GenerateAirports() []models.Airport {
	return []models.Airport{
		{Code: "SVO", Name: "Шереметьево", City: "Москва", Country: "Россия"},
		{Code: "DME", Name: "Домодедово", City: "Москва", Country: "Россия"},
		{Code: "VKO", Name: "Внуково", City: "Москва", Country: "Россия"},
		{Code: "LED", Name: "Пулково", City: "Санкт-Петербург", Country: "Россия"},
		{Code: "KRR", Name: "Пашковский", City: "Краснодар", Country: "Россия"},
		{Code: "OVB", Name: "Толмачёво", City: "Новосибирск", Country: "Россия"},
		{Code: "KHV", Name: "Новый", City: "Хабаровск", Country: "Россия"},
		{Code: "CDG", Name: "Шарль-де-Голль", City: "Париж", Country: "Франция"},
		{Code: "JFK", Name: "Джон Кеннеди", City: "Нью-Йорк", Country: "США"},
		{Code: "DXB", Name: "Дубай", City: "Дубай", Country: "ОАЭ"},
	}
}

// GenerateAirlines - генерация авиакомпаний
func (g *TestDataGenerator) GenerateAirlines() []models.Airline {
	return []models.Airline{
		{Code: "SU", Name: "Аэрофлот", Country: "Россия"},
		{Code: "S7", Name: "S7 Airlines", Country: "Россия"},
		{Code: "U6", Name: "Уральские авиалинии", Country: "Россия"},
		{Code: "FV", Name: "Россия", Country: "Россия"},
		{Code: "AF", Name: "Air France", Country: "Франция"},
		{Code: "LH", Name: "Lufthansa", Country: "Германия"},
		{Code: "TK", Name: "Turkish Airlines", Country: "Турция"},
		{Code: "EY", Name: "Etihad Airways", Country: "ОАЭ"},
	}
}

// GenerateFlights - генерация рейсов
func (g *TestDataGenerator) GenerateFlights(count int, airlines []models.Airline, airports []models.Airport) []models.Flight {
	flights := make([]models.Flight, count)

	for i := 0; i < count; i++ {
		airline := airlines[g.random.Intn(len(airlines))]
		departureAirport := airports[g.random.Intn(len(airports))]
		arrivalAirport := airports[g.random.Intn(len(airports))]

		// Убедимся что аэропорты разные (сравниваем по коду)
		for departureAirport.Code == arrivalAirport.Code {
			arrivalAirport = airports[g.random.Intn(len(airports))]
		}

		// Генерируем время вылета в ближайшие 30 дней
		departureTime := time.Now().Add(time.Duration(g.random.Intn(30*24)) * time.Hour)
		// Продолжительность полета 1-8 часов
		duration := time.Duration(1+g.random.Intn(8)) * time.Hour
		arrivalTime := departureTime.Add(duration)

		// Цена от 2000 до 30000 рублей
		price := 2000.0 + float64(g.random.Intn(28000))

		// Доступные места от 50 до 300
		availableSeats := 50 + g.random.Intn(250)

		flights[i] = models.Flight{
			FlightNumber:   fmt.Sprintf("%s%d", airline.Code, 1000+g.random.Intn(9000)),
			Airline:        airline.Name, // Используем название авиакомпании
			DepartureCity:  departureAirport.City,
			ArrivalCity:    arrivalAirport.City,
			DepartureTime:  departureTime,
			ArrivalTime:    arrivalTime,
			Price:          price,
			AvailableSeats: availableSeats,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
	}

	return flights
}

// GenerateBookings - генерация бронирований с использованием реальных ID
func (g *TestDataGenerator) GenerateBookings(count int, userIDs []uint, flightIDs []uint) []models.Booking {
	bookings := make([]models.Booking, count)

	classes := []string{"economy", "business", "first"}
	statuses := []string{"confirmed", "pending", "cancelled"}
	paymentStatuses := []string{"paid", "pending", "failed"}

	for i := 0; i < count; i++ {
		// Выбираем случайные ID из существующих
		userID := userIDs[g.random.Intn(len(userIDs))]
		flightID := flightIDs[g.random.Intn(len(flightIDs))]

		seats := 1 + g.random.Intn(4) // 1-4 места
		class := classes[g.random.Intn(len(classes))]
		status := statuses[g.random.Intn(len(statuses))]
		paymentStatus := paymentStatuses[g.random.Intn(len(paymentStatuses))]

		// Общая сумма (можно сделать более сложную логику)
		totalAmount := float64(seats) * (1000.0 + float64(g.random.Intn(20000)))

		// Дата бронирования - случайная дата в последние 30 дней
		bookingDate := time.Now().Add(-time.Duration(g.random.Intn(30*24)) * time.Hour)

		bookings[i] = models.Booking{
			UserID:        userID,
			FlightID:      flightID,
			Seats:         seats,
			Class:         class,
			BookingDate:   bookingDate,
			TotalAmount:   totalAmount,
			Status:        status,
			PaymentStatus: paymentStatus,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}

	return bookings
}
