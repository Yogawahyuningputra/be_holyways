package models

import "time"

type Donation struct {
	ID         int    `json:"id"`
	Money      int    `json:"money" gorm:"type:varchar(255)"`
	Status     string `json:"status" gorm:"type:varchar(255)"`
	Attachment string `json:"image" gorm:"type:varchar(255)"`
	// User       UserProfile      `json:"user" gorm:"type:varchar(255)"`
	UserID int `json:"user_id" gorm:"foreignKey:user_id"`
	// Funding    FundingResponses `json:"funding" gorm:"type:varchar(255)"`
	FundingID int       `json:"funding_id" gorm:"foreignKey:funding_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAT time.Time `json:"-"`
}
