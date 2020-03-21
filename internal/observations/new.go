package observations

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/go-playground/validator.v9"
)

// validate holds the settings and caches for validating request struct values.
var validate = validator.New()

// New creates a new Observation from a NewObservation
func New(newObs NewObservation, id primitive.ObjectID, now time.Time) (Observation, error) {
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
		return Observation{}, validationError(err)
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
		Process:      newObs.Process,

		FeatureID:      newObs.Feature.ID,
		FeatureTypeID:  newObs.FeatureType.ID,
		PropertyID:     newObs.Property.ID,
		PropertyTypeID: newObs.PropertyType.ID,
		ProcessID:      newObs.Process.ID,

		Context: newObs.Context,
		Tags:    newObs.Tags,

		Scale:  newObs.Scale,
		Result: newObs.Result,
	}, nil
}

// Save persists an Observation to the database. It expects that the observation
// has been initiated with New and is not a Observation literal.
func Save(ctx context.Context, collection *mongo.Collection, obs Observation) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, obs)

	if err != nil {
		return errors.Wrap(err, "saving observation")
	}

	return nil
}

// Find retrieves a single observation from the database based on the observation id.
func Find(ctx context.Context, collection *mongo.Collection, id string) (Observation, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Observation{}, errors.Wrap(err, "parsing id")
	}

	var obs Observation
	err = collection.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&obs)
	if err != nil {
		return obs, errors.Wrap(err, "finding observation")
	}

	return obs, nil
}

type Filter struct {
	Path    string `json:"path" validate:"required"`
	Op      string `json:"op" validate:"required"`
	Matcher string `json:"match" validate:"required"`
}

var filterableFields map[string]string = map[string]string{
	"featureId": "FeatureID", "featureTypeId": "FeatureTypeID", "propertyId": "PropertyID",
	"propertyTypeId": "PropertyTypeID", "processId": "ProcessID", "id": "ID",
}

var queryOps map[string]string = map[string]string{
	"=": "==", "in": "in",
}

// Get retrieves a list of observations from the databse.
func Get(ctx context.Context, collection *mongo.Collection, filters ...Filter) ([]Observation, error) {

	var observations []Observation
	opts := options.Find().SetSort(bson.D{{"resultTime", 1}})
	cursor, err := collection.Find(context.TODO(), bson.D{{"name", "Bob"}}, opts)
	if err != nil {
		return observations, errors.Wrap(err, "fetching observations")
	}

	if err = cursor.All(context.TODO(), &observations); err != nil {
		return observations, errors.Wrap(err, "decoding observations")
	}

	return observations, nil
}
