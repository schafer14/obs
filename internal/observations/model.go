package observations

import (
	"time"

	geojson "github.com/paulmach/go.geojson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	// Metadata about the observation
	PhenomenonTime      time.Time         `json:"phenomenonTime"`
	ResultTime          time.Time         `json:"resultTime"`
	ValidInterval       Interval          `json:"validInterval"`
	PhenomenonLocation  *geojson.Geometry `json:"phenomenonLocation,omitempty"`
	ObservationLocation *geojson.Geometry `json:"observationLocation,omitempty"`

	// Type data about the observation
	Feature      Referenceable `json:"feature"`
	FeatureType  Referenceable `json:"featureType"`
	Property     Referenceable `json:"property"`
	PropertyType Referenceable `json:"propertyType"`
	Process      Referenceable `json:"process"`

	// Additional fields for indexing and querying
	FeatureID      string `json:"-"`
	FeatureTypeID  string `json:"-"`
	PropertyID     string `json:"-"`
	PropertyTypeID string `json:"-"`
	ProcessID      string `json:"-"`

	Tags    map[string]string `json:"tags,omitempty"`
	Context []string          `json:"context,omitempty"`

	Result interface{} `json:"result"`
	Scale  string      `json:"scale,omitempty"`
}

// Referenceable field is a field that can be looked up with an ID and additionally has a human
// readable label and reference to an external url.
type Referenceable struct {
	ID          string `json:"id" validate:"required,uuid|uri"`
	Label       string `json:"label,omitempty" validate:"omitempty"`
	Description string `json:"description,omitempty" validate:"omitempty"`
	Reference   string `json:"reference,omitempty" validate:"omitempty,uri"`
}

// Interval is a period of a time with a start time and a duration.
type Interval struct {
	StartTime time.Time     `json:"startTime" validate:"-"`
	Duration  time.Duration `json:"duration,omitempty" validate:"-"`
}
