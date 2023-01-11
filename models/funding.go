package models

type Funding struct {
	ID          int    `json:"id"`
	Title       string `json:"title" gorm:"type: varchar(255)"`
	Goals       string `json:"goals" gorm:"type: varchar(255)"`
	Description string `json:"description" gorm:"type: varchar(255)"`
	Image       string `json:"image" gorm:"type: varchar(255)"`
}

type FundingResponses struct {
	ID          int    `json:"id"`
	Title       string `json:"title" gorm:"type: varchar(255)"`
	Goals       string `json:"goals" gorm:"type: varchar(255)"`
	Description string `json:"description" gorm:"type: varchar(255)"`
	Image       string `json:"image" gorm:"type: varchar(255)"`
}

func (FundingResponses) TableName() string {
	return "fundings"
}
