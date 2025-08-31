package dto

import "matching_system/internal/models"

// AddPersonRequest represents the request body for adding a new person
type AddPersonRequest struct {
	Name          string `json:"name" binding:"required"`
	Height        int    `json:"height" binding:"required,min=100,max=250"`
	Gender        string `json:"gender" binding:"required,oneof=male female"`
	NumberOfDates int    `json:"number_of_dates" binding:"required,min=0"`
}

type AddPersonResponse struct {
	Person  models.Person   `json:"person"`
	Matches []*models.Match `json:"matches"`
	Message string          `json:"message"`
}

type RemovePersonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type QueryPeopleResponse struct {
	People  []*models.Person `json:"people"`
	Total   int              `json:"total"`
	Message string           `json:"message"`
}
