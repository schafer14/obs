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
	Feature      Referenceable `json:"feature" validate:"required,dive"`
	FeatureType  Referenceable `json:"featureType" validate:"required,dive"`
	Property     Referenceable `json:"property" validate:"required,dive"`
	PropertyType Referenceable `json:"propertyType" validate:"required,dive"`
	Process      Referenceable `json:"process" validate:"required,dive"`

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
	PhenomenonTime      time.Time         `json:"phenomenonTime" firestore="phenomenonTime"`
	ResultTime          time.Time         `json:"resultTime" firestore="resultTime"`
	ValidInterval       Interval          `json:"validInterval" firestore="validInterval"`
	PhenomenonLocation  *geojson.Geometry `json:"phenomenonLocation" firestore="phenomenonLocation"`
	ObservationLocation *geojson.Geometry `json:"observationLocation" firestore="observationLocation"`

	// Type data about the observation
	Feature      Referenceable `json:"feature" firestore="feature"`
	FeatureType  Referenceable `json:"featureType" firestore="featureType"`
	Property     Referenceable `json:"property" firestore="property"`
	PropertyType Referenceable `json:"propertyType" firestore="propertyType"`
	Process      Referenceable `json:"process" firestore="process"`

	// Additional fields for indexing and querying
	FeatureID      string `json:"-" firestore="featureId"`
	FeatureTypeID  string `json:"-" firestore="featureTypeId"`
	PropertyID     string `json:"-" firestore="propertyId"`
	PropertyTypeID string `json:"-" firestore="propertyTypeId`
	ProcessID      string `json:"-" firestore="processId"`

	Tags    map[string]string `json:"tags" firestore="tags"`
	Context []string          `json:"context" firestore="context"`

	Result interface{} `json:"result" firestore="result"`
	Scale  string      `json:"scale" firestore="scale"`
}

// Referenceable field is a field that can be looked up with an ID and additionally has a human
// readable label and reference to an external url.
type Referenceable struct {
	ID          string `json:"id" validate:"required,uuid|uri" firestore="id"`
	Label       string `json:"label,omitempty" validate:"omitempty" firestore="label"`
	Description string `json:"description,omitempty" validate:"omitempty" firestore="description"`
	Reference   string `json:"reference,omitempty" validate:"omitempty,uri" firestore="reference"`
}

// Interval is a period of a time with a start time and a duration.
type Interval struct {
	StartTime time.Time     `json:"startTime" validate:"-"`
	Duration  time.Duration `json:"duration" validate:"-"`
}
