package models

type User struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email" gorm:"type:varchar(255);unique;not null"`
	Password string `json:"password" gorm:"type:varchar(255)"`
	Phone    int    `json:"phone" gorm:"type:varchar(255)"`
	Gender   string `json:"gender" gorm:"type:varchar(255)"`
	Address  string `json:"address" gorm:"type:varchar(255)"`
	Image    string `json:"image" gorm:"type:varchar(255)"`
	Role     string `json:"role" gorm:"type:varchar(255)"`
}
type UserProfile struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
	Image    string `json:"image"`
}

func (UserProfile) TableName() string {
	return "users"
}
