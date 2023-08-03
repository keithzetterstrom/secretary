package models

import "github.com/keithzetterstrom/secretary/internal/pkg/models"

func UserRegistrationToModel(ur *UserRegistration) models.UserRegistration {
	return models.UserRegistration{
		Email:      ur.Email,
		Name:       ur.Name,
		TgUserName: ur.TgUserName,
	}
}
