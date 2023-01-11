package fundingdto

type FundingRequest struct {
	Title       string `json:"title" form:"title" validate:"required"`
	Goals       string `json:"goals" form:"goals" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	Image       string `json:"image" form:"image"`
}
