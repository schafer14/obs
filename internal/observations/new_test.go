package observations_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	geojson "github.com/paulmach/go.geojson"
	"github.com/schafer14/observations/internal/observations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewObservationWithRequiredFields(t *testing.T) {

	// Arrange
	newObs := mkObs()
	id := uuid.New().String()
	now := time.Now()

	// Act
	obs, err := observations.New(newObs, id, now)

	// Assert
	require.Nil(t, err, "creating observation")

	assert.Equal(t, id, obs.ID, "invalid id")
	assert.Equal(t, now, obs.PhenomenonTime, "phenomenon time not generated")
	assert.Equal(t, now, obs.ResultTime, "result time not generated")
	assert.Equal(t, now, obs.ValidInterval.StartTime, "valid interval start time not generated")
	assert.Equal(t, time.Duration(0), obs.ValidInterval.Duration, "valid interval duration not set")

	assert.Equal(t, newObs.Feature, obs.Feature, "invalid feature")
	assert.Equal(t, newObs.FeatureType, obs.FeatureType, "invalid feature type")
	assert.Equal(t, newObs.Property, obs.Property, "invalid property")
	assert.Equal(t, newObs.PropertyType, obs.PropertyType, "invalid property type")
	assert.Equal(t, newObs.Process, obs.Process, "invalid process")

	assert.Equal(t, newObs.Result, obs.Result, "invalid result")
}

func TestReturnValidationError(t *testing.T) {

	// Arrange
	newObs := observations.NewObservation{}
	now := time.Now()
	id := uuid.New().String()

	// Act
	_, err := observations.New(newObs, id, now)

	// Assert
	require.Error(t, err, "no error for empty object")
}

func TestOptionalParamsMaybeProvided(t *testing.T) {

	// Arange
	newObs := mkObs()
	now := time.Now()
	id := uuid.New().String()
	newObs.Process = observations.Process{ID: "https://example.com/process", Label: "Process"}
	newObs.Context = []string{"These are my contexts", "They are okay"}
	newObs.Tags = map[string]string{"isCool": "true", "observed by": "Banner"}
	newObs.Scale = "inches"

	// Act
	obs, err := observations.New(newObs, id, now)

	// Assert
	require.Nil(t, err, "creating observation")

	assert.Equal(t, newObs.Process.ID, obs.Process.ID, "process id invalid")
	assert.Equal(t, newObs.Process.Label, obs.Process.Label, "process label invalid")
	assert.Equal(t, newObs.Context, obs.Context, "context invalid")
	assert.Equal(t, newObs.Tags, obs.Tags, "tag invalid")
	assert.Equal(t, newObs.Scale, obs.Scale, "scale invalid")
}

func TestObservationWithGeoJson(t *testing.T) {

	// Arange
	newObs := mkObs()
	now := time.Now()
	id := uuid.New().String()
	newObs.PhenomenonLocation = geojson.NewPointGeometry([]float64{1, 2})
	newObs.ObservationLocation = nil

	// Act
	obs, err := observations.New(newObs, id, now)

	// Assert
	require.Nil(t, err, "creating observation")

	assert.Equal(t, newObs.PhenomenonLocation, obs.PhenomenonLocation, "process id invalid")
	assert.Equal(t, newObs.ObservationLocation, obs.ObservationLocation, "process label invalid")
}

// mkObs is a text fixture. It creates a NewObservation
// with some hard coded data.
func mkObs() observations.NewObservation {
	return observations.NewObservation{
		Feature: observations.Feature{
			ID:    "https://example.com/banners-garden",
			Label: "Banner's Garden",
		},
		FeatureType: observations.FeatureType{
			ID:    "urn:example:garden",
			Label: "Garden",
		},
		Property: observations.Property{
			ID:    "urn:example:health",
			Label: "Health",
		},
		PropertyType: observations.PropertyType{
			ID:    "urn:example:scale-1-5",
			Label: "Scale 1 - 5",
		},
		Process: observations.Process{
			ID:    "urn:example:measurement:by-eye",
			Label: "Walk of garden",
		},
		Result: map[string]int{
			"wisteria": 5,
			"citrus":   4,
			"magnolia": 3,
		},
	}
}
