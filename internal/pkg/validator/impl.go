package validator

import (
	"errors"
	"fmt"
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

// Validate validates the data (struct or slice of struct)
// using the validator and returns an error if the data is invalid.
func (v *Validate) Validate(
	data any,
) error {
	strErrs := []string{}

	dataAsSlice, ok := data.([]any)
	if !ok {
		dataAsSlice = []any{data}
	}

	for _, item := range dataAsSlice {
		err := v.validate.Struct(item)

		if err == nil {
			continue
		}

		validatorErrs := err.(validator.ValidationErrors)

		for _, e := range validatorErrs {
			translatedErr := fmt.Errorf(
				e.Translate(v.trans),
			)
			strErrs = append(
				strErrs,
				translatedErr.Error(),
			)
		}
	}

	errMsg := strings.Join(
		strErrs,
		", ",
	)

	return errors.New(errMsg)
}
