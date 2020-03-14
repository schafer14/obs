package observations

import (
	"fmt"
	"reflect"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

var translator *ut.UniversalTranslator

func init() {

	// Instantiate the english locale for the validator library.
	enLocale := en.New()

	// Create a value using English as the fallback locale (first argument).
	// Provide one or more arguments for additional supported locales.
	translator = ut.New(enLocale, enLocale)

	// Register the english error messages for validation errors.
	lang, _ := translator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, lang)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

var (
	ErrorNotFound = errors.New("observation not found")
)

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ValidationError struct {
	Err    string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error fulfills the error interface for validation error.
func (e *ValidationError) Error() string {
	s := e.Err
	for _, f := range e.Fields {
		s += fmt.Sprintf("\t%v: %v", f.Field, f.Error)
	}

	return s
}

func validationError(err error) error {
	// Use a type assertion to get the real error value.
	verrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	lang, _ := translator.GetTranslator("en")
	var fields []FieldError
	for _, verror := range verrors {
		field := FieldError{
			Field: verror.Namespace(),
			Error: verror.Translate(lang),
		}
		fields = append(fields, field)
	}

	return &ValidationError{
		Err:    "error validating observation",
		Fields: fields,
	}
}
