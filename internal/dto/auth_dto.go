package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserProfileResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Avatar    *string `json:"avatar"`
	CreatedAt string  `json:"createdAt"`
}

type AuthResponse struct {
	Success bool                `json:"success"`
	Token   string              `json:"token"`
	User    UserProfileResponse `json:"user"`
}
