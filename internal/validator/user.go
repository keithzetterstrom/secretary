package validator

import "github.com/keithzetterstrom/secretary/internal/http/models"

func (v *Validator) ValidateUserRegistration(userRegistration models.UserRegistration) error {
	err := v.validate.Struct(userRegistration)
	if err != nil {
		return err
	}

	return nil
}
