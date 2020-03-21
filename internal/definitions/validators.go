package definitions

import (
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"
)

func validateTextGoal(r io.Reader) (interface{}, error) {
	type Goal struct {
		Goal string `json:"goal" validate:"required"`
	}

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	var g Goal
	err := decoder.Decode(&g)
	if err != nil {
		return nil, err
	}

	// Run validation.
	if err := validate.Struct(&g); err != nil {
		return g, errors.Wrap(err, "Validating goal")
	}

	return g.Goal, nil
}

func validateStructuredGoal(r io.Reader) (interface{}, error) {
	type StructuredGoal struct {
		Goal  string `json:"goal" validate:"required"`
		Steps []struct {
			Description string    `json:"accomplish" validate:"required"`
			By          time.Time `json:"by,omitempty" validate:"-"`
		} `json:"steps" validate:"required,dive"`
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

	return g, nil
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
