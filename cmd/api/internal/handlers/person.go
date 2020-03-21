package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/schafer14/observations/internal/observations"
	"github.com/schafer14/observations/internal/people"
	"go.mongodb.org/mongo-driver/mongo"
)

type PersonHandler struct {
	personCollection      *mongo.Collection
	observationCollection *mongo.Collection
}

func (p *PersonHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newPerson people.NewPerson
	if err := Decode(r, &newPerson); err != nil {
		RespondError(ctx, w, err)
		return
	}

	id := uuid.New().String()
	person, err := people.New(newPerson, id)
	if err != nil {
		if vError, ok := err.(*observations.ValidationError); ok {
			Respond(ctx, w, vError, http.StatusUnprocessableEntity)
			return
		}
		RespondError(ctx, w, errors.Wrap(err, "creating new person"))
		return
	}

	err = people.Save(ctx, p.personCollection, person)
	if err != nil {
		Respond(ctx, w, "unable to save observation", http.StatusInternalServerError)
		return
	}

	Respond(ctx, w, person, http.StatusCreated)
}
