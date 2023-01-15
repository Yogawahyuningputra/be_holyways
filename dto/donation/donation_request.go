package donationdto

import "time"

type DonationRequest struct {
	Money     int       `json:"money" form:"money" validate:"required"`
	Status    string    `json:"status" form:"status"`
	UserID    int       `json:"user_id"`
	FundingID int       `json:"funding_id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
