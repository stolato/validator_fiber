package validator_fiber

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorsHandle struct {
	FailedField string
	Tag         string
	Value       interface{}
	Error       bool
}

func Validator(item interface{}) []ErrorsHandle {

	var validationErrors []ErrorsHandle
	errs := validate.Struct(item)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorsHandle

			elem.FailedField = strings.ToLower(err.Field())
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}
	return validationErrors
}
