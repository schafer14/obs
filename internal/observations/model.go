package observations

import (
	"time"

	geojson "github.com/paulmach/go.geojson"
)

// NewObservation is an observation that has not yet been initalized.
// NewObservations cannot be saved, but can be initialized to a saveable
// observation with the New function.
type NewObservation struct {

	// Metadata about the observation
	PhenomenonTime      time.Time         `json:"phenomenonTime,omitempty" validate:"-"`
	ResultTime          time.Time         `json:"resultTime,omitempty" validate:"-"`
	ValidInterval       Interval          `json:"validInterval,omitempty" validate:"omitempty,dive"`
	PhenomenonLocation  *geojson.Geometry `json:"phenomenonLocation,omitempty" validate:"-"`
	ObservationLocation *geojson.Geometry `json:"observationLocation,omitempty" validate:"-"`

	// Type data about the observation
	Feature      Feature      `json:"feature" validate:"required,dive"`
	FeatureType  FeatureType  `json:"featureType" validate:"required,dive"`
	Property     Property     `json:"property" validate:"required,dive"`
	PropertyType PropertyType `json:"propertyType" validate:"required,dive"`
	Process      Process      `json:"process" validate:"required,dive"`

	Tags    map[string]string `json:"tags,omitempty" validate:"-"`
	Context []string          `json:"context,omitempty" validate:"-"`

	Result interface{} `json:"result" validate:"required"`
	Scale  string      `json:"scale,omitempty" validate:"-"`
}

// Observation is an observation that can be persisted to a database.
type Observation struct {

	// Identity
	ID string `json:"id"`

	// Metadata about the observation
	PhenomenonTime      time.Time         `json:"phenomenonTime"`
	ResultTime          time.Time         `json:"resultTime"`
	ValidInterval       Interval          `json:"validInterval"`
	PhenomenonLocation  *geojson.Geometry `json:"phenomenonLocation"`
	ObservationLocation *geojson.Geometry `json:"observationLocation"`

	// Type data about the observation
	Feature      Feature      `json:"feature"`
	FeatureType  FeatureType  `json:"featureType"`
	Property     Property     `json:"property"`
	PropertyType PropertyType `json:"propertyType"`
	Process      Process      `json:"process"`

	Tags    map[string]string `json:"tags"`
	Context []string          `json:"context"`

	Result interface{} `json:"result"`
	Scale  string      `json:"scale"`
}

// Feature is the entity being observed.
type Feature struct {
	ID          string `json:"id" validate:"required,uuid|uri"`
	Label       string `json:"label,omitempty" validate:"omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Reference   string `json:"reference,omitempty" validate:"omitempty,uri"`
}

// FeatureType is the type of the entity being observed.
type FeatureType struct {
	ID          string `json:"id" validate:"required,uuid|uri"`
	Label       string `json:"label,omitempty" validate:"omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Reference   string `json:"reference,omitempty" validate:"omitempty,uri"`
}

// Property is the attribute of the feature that is being observed.
type Property struct {
	ID          string `json:"id" validate:"required,uuid|uri"`
	Label       string `json:"label,omitempty" validate:"omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Reference   string `json:"reference,omitempty" validate:"omitempty,uri"`
}

// PropertyType is the type of the attribute being observed.
type PropertyType struct {
	ID          string `json:"id" validate:"required,uuid|uri"`
	Label       string `json:"label,omitempty" validate:"omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Reference   string `json:"reference,omitempty" validate:"omitempty,uri"`
}

// Process is the method that was used to make the observation.
type Process struct {
	ID          string `json:"id" validate:"required,uuid|uri"`
	Label       string `json:"label,omitempty" validate:"omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Reference   string `json:"reference,omitempty" validate:"omitempty,uri"`
}

// Interval is a period of a time with a start time and a duration.
type Interval struct {
	StartTime time.Time     `json:"startTime" validate:"-"`
	Duration  time.Duration `json:"duration" validate:"-"`
}
