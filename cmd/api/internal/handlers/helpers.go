package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// WebError is a type of error that has a field that is safe to
// return to the client.
type webError struct {
	error

	message string
	status  int
	fields  []interface{}
}

// NewWebError creates an error that will have it's message returned to
// the end user.
func NewWebError(err error, msg string, status int, fields ...interface{}) webError {
	return webError{err, msg, status, fields}
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(val); err != nil {
		return newRequestError(fmt.Errorf("parsing JSON: %v", err), http.StatusUnprocessableEntity)
	}

	return nil
}

func newRequestError(err error, status int) error {
	return NewWebError(err, err.Error(), status)
}

// RespondError sends an error response back to the client.
func RespondError(ctx context.Context, w http.ResponseWriter, err error) {

	type respondError struct {
		Error  string        `json:"error"`
		Fields []interface{} `json:"fields,omitempty"`
	}

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := err.(webError); ok {
		er := respondError{
			Error:  webErr.message,
			Fields: webErr.fields,
		}
		Respond(ctx, w, er, webErr.status)
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
