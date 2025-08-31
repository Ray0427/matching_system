package models

type Person struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Height int `json:"height"`
	Gender string `json:"gender"`
	NumberOfDates int `json:"number_of_dates"`
}