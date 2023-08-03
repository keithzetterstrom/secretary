package models

type UserRegistration struct {
	Email      string `json:"email" validate:"required,email"`
	Name       string `json:"name"`
	TgUserName string `json:"tg_user_name" validate:"required"`
}
