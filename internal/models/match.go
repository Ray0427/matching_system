package models

import "time"

type Match struct {
	Person1ID string    `json:"person1Id"`
	Person2ID string    `json:"person2Id"`
	Timestamp time.Time `json:"timestamp"`
}