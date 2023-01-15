package donationdto

import (
	"server/models"
	"time"
)

type DonationResponse struct {
	ID        int                     `json:"id"`
	Money     int                     `json:"money"`
	Status    string                  `json:"status"`
	UserID    int                     `json:"user_id"`
	User      models.UserProfile      `json:"user"`
	FundingID int                     `json:"funding_id"`
	Funding   models.FundingResponses `json:"funding"`
	CreatedAt time.Time               `json:"-"`
	UpdatedAt time.Time               `json:"-"`
}
