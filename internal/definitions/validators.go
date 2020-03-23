package definitions

import (
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

var (
	ErrorNoValidatorFound   = errors.New("no validator found for property type")
	ErrorParsingRequestBody = errors.New("error parsing request body")
)

// ValidationError is an error that shows validation results.
type ValidationError struct {
	Message string
	Fields  []gojsonschema.ResultError
}

func (v ValidationError) Error() string {
	return v.Message
}

// Validate takes an io reader. Reads the content into JSON
// and validates that json against a Property Type and returns
// en error if the validation fails.
func Validate(r io.Reader, propertyType PropertyType) (interface{}, error) {
	if propertyType.Schema == nil && propertyType.SchemaURL == "" {
		return nil, ErrorNoValidatorFound
	}

	// create a json loader from either the schema description or the schema url.
	var validator gojsonschema.JSONLoader
	if propertyType.Schema != nil {
		validator = gojsonschema.NewGoLoader(propertyType.Schema)
	} else {
		validator = gojsonschema.NewReferenceLoader(propertyType.SchemaURL)
	}

	// Read body into an interface.
	decoder := json.NewDecoder(r)
	var requestBody interface{}
	err := decoder.Decode(&requestBody)
	if err != nil {
		return nil, ErrorParsingRequestBody
	}
	data := gojsonschema.NewGoLoader(requestBody)

	// do the actual validation
	result, err := gojsonschema.Validate(validator, data)
	if err != nil {
		return nil, errors.Wrap(err, "validating schema")
	}
	if !result.Valid() {
		return nil, ValidationError{"validation error", result.Errors()}
	}

	return requestBody, nil
}

func validateDailyGoal(r io.Reader) (interface{}, error) {
	type StructuredGoal struct {
		Goal      []string      `json:"goals" validate:"required"`
		TimeFrame time.Duration `json:"timeFrame" validate:"-"`
	}

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var g StructuredGoal
	err := decoder.Decode(&g)
	if err != nil {
		return nil, err
	}

	// Run validation.
	if err := validate.Struct(&g); err != nil {
		return g, errors.Wrap(err, "Validating goal")
	}

	if g.TimeFrame == 0 {
		g.TimeFrame = 24 * time.Hour
	}

	return g, nil
}

func validateDailyGoalResult(r io.Reader) (interface{}, error) {
	type StructuredGoal struct {
		Goals []struct {
			Goal     string `json:"goal" validate:"required"`
			Achieved bool   `json:"achieved" validate:"-"`
		} `json:"goals" validate:"required,dive"`
		TimeFrame time.Duration `json:"timeFrame" validate:"-"`
	}

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	var g StructuredGoal
	err := decoder.Decode(&g)
	if err != nil {
		return nil, err
	}

	// Run validation.
	if err := validate.Struct(&g); err != nil {
		return g, errors.Wrap(err, "Validating goal")
	}

	if g.TimeFrame == 0 {
		g.TimeFrame = 24 * time.Hour
	}

	return g, nil
}
