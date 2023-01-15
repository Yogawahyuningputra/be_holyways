package models

import "time"

type Funding struct {
	ID          int         `json:"id"`
	Title       string      `json:"title" gorm:"type: varchar(255)"`
	Goals       int         `json:"goals" gorm:"type: varchar(255)"`
	Description string      `json:"description" gorm:"type: varchar(255)"`
	Image       string      `json:"image" gorm:"type: varchar(255)"`
	UserID      int         `json:"user_id"`
	User        UserProfile `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type FundingResponses struct {
	ID          int    `json:"id"`
	Title       string `json:"title" gorm:"type: varchar(255)"`
	Goals       int    `json:"goals" gorm:"type: varchar(255)"`
	Description string `json:"description" gorm:"type: varchar(255)"`
	Image       string `json:"image" gorm:"type: varchar(255)"`
}

func (FundingResponses) TableName() string {
	return "fundings"
}
