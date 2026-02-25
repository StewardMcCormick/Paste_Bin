package dto

type UserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type LoginRequest struct {
	UserRequest
}
