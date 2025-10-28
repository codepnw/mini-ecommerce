package validate

import "github.com/go-playground/validator/v10"

func Struct(input any) error {
	v := validator.New()
	return v.Struct(input)
}
