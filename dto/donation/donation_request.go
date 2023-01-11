package donationdto

type DonationRequest struct {
	Money     int    `json:"money" form:"money" validate:"required"`
	Status    string `json:"status" form:"status" validate:"required"`
	UserID    int    `json:"user_id"`
	FundingID int    `json:"funding_id"`
}
