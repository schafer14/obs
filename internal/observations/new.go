package observations

import (
	"time"

	"gopkg.in/go-playground/validator.v9"
)

// validate holds the settings and caches for validating request struct values.
var validate = validator.New()

// New creates a new Observation from a NewObservation
func New(newObs NewObservation, id string, now time.Time) (Observation, error) {
	phenomenonTime, startTime, resultTime := now, now, now

	// Set times to now if they are not provided.
	if !newObs.PhenomenonTime.IsZero() {
		phenomenonTime = newObs.PhenomenonTime
	}

	if !newObs.ResultTime.IsZero() {
		resultTime = newObs.ResultTime
	}

	if !newObs.ValidInterval.StartTime.IsZero() {
		startTime = newObs.ValidInterval.StartTime
	}

	// Run validation.
	if err := validate.Struct(&newObs); err != nil {
		return Observation{}, err
	}

	return Observation{

		ID: id,

		PhenomenonTime: phenomenonTime,
		ResultTime:     resultTime,
		ValidInterval: Interval{
			StartTime: startTime,
			Duration:  newObs.ValidInterval.Duration,
		},
		PhenomenonLocation:  newObs.PhenomenonLocation,
		ObservationLocation: newObs.ObservationLocation,

		Feature:      newObs.Feature,
		FeatureType:  newObs.FeatureType,
		Property:     newObs.Property,
		PropertyType: newObs.PropertyType,

		Process: newObs.Process,
		Context: newObs.Context,
		Tags:    newObs.Tags,

		Scale:  newObs.Scale,
		Result: newObs.Result,
	}, nil
}
