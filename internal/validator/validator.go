package validator

import "gopkg.in/go-playground/validator.v9"

type Validator struct {
	validate *validator.Validate
}

func InitValidate() (*Validator, error) {
	v := &Validator{}

	validate := validator.New()

	v.validate = validate
	return v, nil
}
