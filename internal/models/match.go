package models

import (
	"time"
)

type Match struct {
	Person1   Person    `json:"person1"`
	Person2   Person    `json:"person2"`
	Timestamp time.Time `json:"timestamp"`
}
