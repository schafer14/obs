package people

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"
)

// validate holds the settings and caches for validating request struct values.
var validate = validator.New()

// New initializes a new person so it is ready to persist.
func New(newPerson NewPerson, id string) (Person, error) {

	// Run validation.
	if err := validate.Struct(&newPerson); err != nil {
		return Person{}, errors.Wrap(err, "validation new person")
	}

	return Person{
		ID:    id,
		Name:  newPerson.Name,
		Email: newPerson.Email,
	}, nil
}

// Save persists a person to the database.
func Save(ctx context.Context, coll *mongo.Collection, person Person) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := coll.InsertOne(ctx, person)

	if err != nil {
		return errors.Wrap(err, "saving person")
	}

	return nil
}

// Find retrieves a single person from the database based on the observation id.
func Find(ctx context.Context, collection *mongo.Collection, id string) (Person, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var person Person
	err := collection.FindOne(ctx, bson.D{{"id", id}}).Decode(&person)
	if err != nil {
		return person, errors.Wrap(err, "finding person")
	}

	return person, nil
}

// Get retrieves a list of observations from the databse.
func Get(ctx context.Context, collection *mongo.Collection) ([]Person, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var people []Person
	cursor, err := collection.Find(ctx, bson.D{{}}, nil)
	if err != nil {
		return people, errors.Wrap(err, "fetching observations")
	}

	if err = cursor.All(ctx, &people); err != nil {
		return people, errors.Wrap(err, "decoding people")
	}

	return people, nil
}
