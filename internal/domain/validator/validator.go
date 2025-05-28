package validator

type Validator interface {
	Validate(s any) error
}
