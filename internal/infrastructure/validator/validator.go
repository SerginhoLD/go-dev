package validator

import (
	"exampleapp/internal/domain/validator"

	playground_validator "github.com/go-playground/validator/v10"
)

type ValidatorImpl struct {
	validate *playground_validator.Validate
}

func New() validator.Validator {
	return &ValidatorImpl{playground_validator.New()}
}

func (v *ValidatorImpl) Validate(s any) error {
	return v.validate.Struct(s)
}
