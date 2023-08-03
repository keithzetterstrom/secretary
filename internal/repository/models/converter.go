package models

import "github.com/keithzetterstrom/secretary/internal/pkg/models"

func ModelToUserRegistration(ur *models.UserRegistration) UserRegistration {
	return UserRegistration{
		Email:      ur.Email,
		Name:       ur.Name,
		TgUserName: ur.TgUserName,
	}
}
