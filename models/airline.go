package models

type Airline struct {
	Code    string `json:"airline_code" faker:"len=2"`
	Name    string `json:"airline_name" faker:"sentence"`
	Country string `json:"country_code" faker:"len=2"`
}
