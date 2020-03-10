package observations_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
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

// mkObs is a text fixture. It creates a NewObservation
// with some hard coded data.
func mkObs() observations.NewObservation {
	return observations.NewObservation{
		Feature: &observations.Feature{
			ID:    "https://example.com/banners-garden",
			Label: "Banner's Garden",
		},
		FeatureType: &observations.FeatureType{
			ID:    "urn:example:garden",
			Label: "Garden",
		},
		Property: &observations.Property{
			ID:    "urn:example:health",
			Label: "Health",
		},
		PropertyType: &observations.PropertyType{
			ID:    "urn:example:scale-1-5",
			Label: "Scale 1 - 5",
		},
		Process: &observations.Process{
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
