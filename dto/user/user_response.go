package userdto

type UserResponse struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
	Image    string `json:"image"`
	Role     string `json:"role"`
}
