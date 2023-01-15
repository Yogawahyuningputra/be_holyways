package fundingdto

import "time"

type FundingRequest struct {
	Title       string    `json:"title" form:"title" validate:"required"`
	Goals       int       `json:"goals" form:"goals" validate:"required"`
	Description string    `json:"description" form:"description" validate:"required"`
	Image       string    `json:"image" form:"image"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
