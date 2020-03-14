package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

// validate holds the settings and caches for validating request struct values.
var validate = validator.New()

// translator is a cache of locale and translation information.
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

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		fmt.Println("xyz", err)
		return Error{errors.Wrap(err, "decoding request"), http.StatusUnprocessableEntity, []FieldError{}}
	}

	if err := validate.Struct(val); err != nil {

		// Use a type assertion to get the real error value.
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// lang controls the language of the error messages. You could look at the
		// Accept-Language header if you intend to support multiple languages.
		lang, _ := translator.GetTranslator("en")

		var fields []FieldError
		for _, verror := range verrors {
			field := FieldError{
				Field: verror.Namespace(),
				Error: verror.Translate(lang),
			}
			fields = append(fields, field)
		}

		return Error{fmt.Errorf("unable to validate request"), http.StatusUnprocessableEntity, fields}
	}

	return nil
}

// RespondError sends an error response back to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) {

	type respondError struct {
		Error  string        `json:"error"`
		Fields []interface{} `json:"fields,omitempty"`
	}

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := err.(Error); ok {
		e := ErrorResponse{
			Fields: webErr.Fields,
			Error:  err.Error(),
		}
		Respond(ctx, w, e, webErr.Status)
		return
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := respondError{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	fmt.Printf("Unhandled error: %v\n", err)
	Respond(ctx, w, er, http.StatusInternalServerError)
	return
}

// Respond converts a Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) {

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		RespondError(ctx, w, errors.Wrap(err, "marshalling data"))
		return
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		w.WriteHeader(500)
		RespondError(ctx, w, errors.Wrap(err, "writing data"))
		return
	}

	return
}

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse is the form used for API responses from failures in the API.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error is used to pass an error during the request through the
// application with web specific context.
type Error struct {
	error
	Status int
	Fields []FieldError
}
