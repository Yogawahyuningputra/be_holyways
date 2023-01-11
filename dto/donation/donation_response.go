package donationdto

import "server/models"

type DonationResponse struct {
	Money     int                     `json:"money"`
	Status    string                  `json:"status"`
	UserID    int                     `json:"user_id"`
	User      models.UserProfile      `json:"user"`
	FundingID int                     `json:"funding_id"`
	Funding   models.FundingResponses `json:"funding"`
}
