package people

import (
	"context"
	"time"

	"github.com/pkg/errors"
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
