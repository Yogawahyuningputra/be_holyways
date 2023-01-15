package models

import "time"

type Donation struct {
	ID         int              `json:"id"`
	Money      int              `json:"money" gorm:"type: varchar(255)"`
	Status     string           `json:"status" gorm:"type: varchar(255)"`
	Attachment string           `json:"image" gorm:"type: varchar(255)"`
	UserID     int              `json:"user_id"`
	User       UserProfile      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FundingID  int              `json:"funding_id"`
	Funding    FundingResponses `json:"funding" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
