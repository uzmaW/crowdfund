package models

import "time"

type Donation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `json:"project_id"`
	UserID    uint      `json:"user_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

type CreateDonation struct {
	ProjectID uint    `json:"project_id"`
	UserID    uint    `json:"user_id"`
	Amount    float64 `json:"amount"`
}

type UpdateDonation struct {
	Amount float64 `json:"amount"`
}

type DeleteDonation struct {
	ID uint `json:"id"`
}

type ListDonations struct {
	ProjectID uint `json:"project_id"`
}

type GetDonation struct {
	ID uint `json:"id"`
}

type GetDonationsByProjectID struct {
	ProjectID uint `json:"project_id"`
}

type GetDonationsByUserID struct {
	UserID uint `json:"user_id"`
}
