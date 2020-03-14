package observations

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func Save(ctx context.Context, collection *firestore.CollectionRef, obs Observation) error {
	_, err := collection.Doc(obs.ID).Set(ctx, obs)

	if err != nil {
		return errors.Wrap(err, "saving observation")
	}

	return nil
}

// Find retrieves a single observation from the database based on the observation id.
func Find(ctx context.Context, collection *firestore.CollectionRef, id string) (Observation, error) {
	docsnap, err := collection.Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return Observation{}, ErrorNotFound
		}
		return Observation{}, errors.Wrap(err, "fetching observation")
	}

	var obs Observation
	if err := docsnap.DataTo(&obs); err != nil {
		return Observation{}, errors.Wrap(err, "parsing docsnap to observation")
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
func Get(ctx context.Context, collection *firestore.CollectionRef, filters ...Filter) ([]Observation, error) {
	q := collection.OrderBy("ResultTime", firestore.Desc)

	for _, f := range filters {
		firestoreField, fok := filterableFields[f.Path]
		firestoreOp, ook := queryOps[f.Op]
		if !fok || !ook {
			return []Observation{}, fmt.Errorf("invalid filter parameter")
		}

		q = q.Where(firestoreField, firestoreOp, f.Matcher)
	}

	var observations []Observation
	iter := q.Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// This could be an indexing error that we (developers) need to know about and add an index
			// to firestore.
			fmt.Println(err)
			return observations, errors.Wrap(err, "fetching document in iteration")
		}

		var observation Observation
		if err := doc.DataTo(&observation); err != nil {
			return observations, errors.Wrap(err, "parsing doc snapshot to observation")
		}
		observations = append(observations, observation)
	}

	return observations, nil
}
