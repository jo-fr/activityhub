package validate

import (
	"fmt"

	validatorPkg "github.com/go-playground/validator/v10"
)

var validate *validatorPkg.Validate

// Validator returns a validator instance as singleton
func Validator() *validatorPkg.Validate {
	if validate == nil {
		validate = validatorPkg.New()
	}
	return validate
}

func GetErrorMessage(fieldError validatorPkg.FieldError) string {

	field := fieldError.Field()
	if field == "" {
		field = "parameter"
	}
	switch fieldError.ActualTag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "uuid4":
		return fmt.Sprintf("%s is not a valid uuid", field)
	case "url":
		return fmt.Sprintf("%s is not a valid url", field)

	default:
		return fmt.Sprintf("%s is invalid", field)
	}

}
