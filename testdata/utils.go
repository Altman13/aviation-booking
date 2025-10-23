package testdata

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomDate генерирует случайную дату в диапазоне
func RandomDate(from, to time.Time) time.Time {
	delta := to.Unix() - from.Unix()
	sec := rand.Int63n(delta) + from.Unix()
	return time.Unix(sec, 0)
}

// RandomFlightNumber генерирует случайный номер рейса
func RandomFlightNumber(airlineCode string) string {
	numbers := []string{
		fmt.Sprintf("%s %d", airlineCode, rand.Intn(2000)+1000),
		fmt.Sprintf("%s %d", airlineCode, rand.Intn(900)+100),
	}
	return numbers[rand.Intn(len(numbers))]
}
