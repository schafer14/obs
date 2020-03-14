package handlers

import (
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/schafer14/observations/internal/observations"
)

type ObservationHandler struct {
	db *firestore.CollectionRef
}

// Create handles an http request that creates a new observation.
func (o *ObservationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var newObs observations.NewObservation
	if err := Decode(r, &newObs); err != nil {
		RespondError(ctx, w, err)
		return
	}

	id := uuid.New().String()
	now := time.Now()
	obs, err := observations.New(newObs, id, now)
	if err != nil {
		Respond(ctx, w, "invalid observations", http.StatusUnprocessableEntity)
		return
	}

	err = observations.Save(ctx, o.db, obs)
	if err != nil {
		Respond(ctx, w, "unable to save observation", http.StatusInternalServerError)
		return
	}

	Respond(ctx, w, obs, http.StatusCreated)
}

// Get handles an http request for listing observations.
func (o *ObservationHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filters := []observations.Filter{}
	if err := Decode(r, &filters); err != nil && r.ContentLength > 0 {
		RespondError(ctx, w, errors.Wrap(err, "decoding filters"))
		return
	}

	obs, err := observations.Get(ctx, o.db, filters...)
	if err != nil {
		RespondError(ctx, w, errors.Wrap(err, "fetching observations"))
		return
	}

	if obs == nil {
		obs = []observations.Observation{}
	}

	Respond(ctx, w, obs, http.StatusOK)
}

// Find handles an http request for finding a single observation.
func (o *ObservationHandler) Find(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	obs, err := observations.Find(ctx, o.db, id)
	if err != nil {
		RespondError(ctx, w, errors.Wrap(err, "fetching observation"))
		return
	}

	Respond(ctx, w, obs, http.StatusOK)
}
