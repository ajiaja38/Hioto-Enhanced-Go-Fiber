package validators

import (
	"github.com/go-playground/validator/v10"
)

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("binary", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if value == "" {
			return false
		}

		for _, char := range value {
			if char != '0' && char != '1' {
				return false
			}
		}

		return true
	})
}
