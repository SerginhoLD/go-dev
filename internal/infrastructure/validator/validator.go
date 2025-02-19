package validator

import (
	"github.com/go-playground/validator/v10"
)

type ValidatorImpl struct {
	validate *validator.Validate
}

func New() *ValidatorImpl {
	return &ValidatorImpl{validator.New()}
}

func (v *ValidatorImpl) Validate(s interface{}) error {
	return v.validate.Struct(s)
}
