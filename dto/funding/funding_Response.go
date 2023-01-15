package fundingdto

import (
	"server/models"
	"time"
)

type FundingResponse struct {
	ID          int                `json:"id"`
	Title       string             `json:"title"`
	Goals       int                `json:"goals"`
	Description string             `json:"description"`
	Image       string             `json:"image"`
	UserID      int                `json:"user_id"`
	User        models.UserProfile `json:"user"`
	CreatedAt   time.Time          `json:"-"`
	UpdatedAt   time.Time          `json:"-"`
}
