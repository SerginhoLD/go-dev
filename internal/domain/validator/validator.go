package validator

type Validator interface {
	Validate(s interface{}) error
}
