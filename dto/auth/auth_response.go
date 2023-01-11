package authdto

type LoginResponse struct {
	Email string `json:"email" gorm:"type: varchar(255)"`
	Token string `json:"token" gorm:"type: varchar(255)"`
	Role  string `json:"role" form:"role"`
}
type RegisterResponse struct {
	Email string `json:"email" gorm:"type: varchar(255)"`
	Token string `json:"token" gorm:"type: varchar(255)"`
}
type CheckAuthResponse struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Gender   string `json:"gender" form:"gender"`
	Phone    int    `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	Role     string `json:"role" form:"role"`
	Image    string `json:"image" form:"image"`
}
