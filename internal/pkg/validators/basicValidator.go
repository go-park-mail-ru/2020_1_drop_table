package validators

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
)

//ToDo refactor function
func GetValidator() (*validator.Validate, ut.Translator, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	locale := "en"
	trans, found := uni.GetTranslator(locale)
	if !found {
		err := errors.New(fmt.Sprintf("translator for %s not found", locale))
		return nil, nil, err
	}

	v := validator.New()

	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, nil, err
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{{0} is a required field}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	return v, trans, nil
}
