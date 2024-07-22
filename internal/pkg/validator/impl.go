package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validate struct {
	validate *validator.Validate
	trans    ut.Translator
}

func NewValidate() *Validate {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, ok := uni.GetTranslator("en")
	if !ok {
		panic("translator not found")
	}

	if err := enTranslations.
		RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}

	return &Validate{
		validate,
		trans,
	}
}

// Validate validates the data (struct)
// returning an error if the data is invalid.
func (v *Validate) Validate(
	data any,
) error {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	strErrs := make([]string, len(validationErrs))
	for i, validationErr := range validationErrs {
		strErrs[i] = validationErr.Translate(v.trans)
	}

	errMsg := strings.Join(
		strErrs,
		", ",
	)

	return errors.New(errMsg)
}
