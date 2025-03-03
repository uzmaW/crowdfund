package models

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Goal        float64   `json:"goal"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	UserID      uint      `json:"user_id"` // Creator of the project
}

type CreateProject struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Goal        float64   `json:"goal"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}
