package userdto

type CreateUserRequest struct {
	Fullname string `json:"fullname" form:"fullname" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Phone    int    `json:"phone" form:"phone" validate:"required"`
	Gender   string `json:"gender" form:"gender" validate:"required"`
	Address  string `json:"address" form:"address" validate:"required"`
	Image    string `json:"image" form:"image"`
	Role     string `json:"role" form:"role"`
}

type UpdateUserRequest struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email" `
	Password string `json:"password" form:"password"`
	Phone    int    `json:"phone" form:"phone"`
	Gender   string `json:"gender" form:"gender"`
	Address  string `json:"address" form:"address"`
	Image    string `json:"image" form:"image"`
	Role     string `json:"role" form:"role"`
}
