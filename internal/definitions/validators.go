package definitions

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"go.mongodb.org/mongo-driver/bson"
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
func Validate(r io.Reader, propertyType PropertyType) (bson.M, error) {
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
	var requestBody bson.M
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
