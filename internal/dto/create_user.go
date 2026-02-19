package dto

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,len=60"`
}
