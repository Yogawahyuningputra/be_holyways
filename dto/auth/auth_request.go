package authdto

type RegisterRequest struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Phone    int    `json:"phone" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
