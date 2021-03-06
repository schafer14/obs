package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/schafer14/obs/internal/definitions"
	"github.com/schafer14/obs/internal/observations"
	"go.mongodb.org/mongo-driver/mongo"
)

type ObservationHandler struct {
	db *mongo.Collection
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
		if vError, ok := err.(*observations.ValidationError); ok {
			Respond(ctx, w, vError, http.StatusUnprocessableEntity)
			return
		}
		RespondError(ctx, w, errors.Wrap(err, "creating new observation"))
		return
	}

	err = observations.Save(ctx, o.db, obs)
	if err != nil {
		Respond(ctx, w, "unable to save observation", http.StatusInternalServerError)
		return
	}

	Respond(ctx, w, obs, http.StatusCreated)
}

type SearchParams struct {
	Filters []observations.Filter `json:"filters" validate:"omitempty,dive"`
}

// Get handles an http request for listing observations.
func (o *ObservationHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var filters SearchParams
	qs := r.URL.Query().Get("q")
	if err := DecodeAny(strings.NewReader(qs), &filters); err != nil && r.ContentLength > 0 {
		RespondError(ctx, w, err)
		return
	}

	obs, err := observations.Get(ctx, o.db, filters.Filters...)
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
		if err == observations.ErrorNotFound {
			Respond(ctx, w, map[string]string{"error": "observation not found"}, http.StatusNotFound)
			return
		}
		RespondError(ctx, w, errors.Wrap(err, "fetching observation"))
		return
	}

	Respond(ctx, w, obs, http.StatusOK)
}

// Generic makes an observation on a specific type based from the Definitions data store.
func (o *ObservationHandler) Generic(featureTypeSlug string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := definitions.Data
		ft, ok := data[featureTypeSlug]
		if !ok {
			Respond(ctx, w, fmt.Sprintf("Could not find feature type %v", featureTypeSlug), http.StatusUnprocessableEntity)
			return
		}

		propertySlug := chi.URLParam(r, "propertySlug")
		property, ok := ft.Properties[propertySlug]
		if !ok {
			Respond(ctx, w, fmt.Sprintf("Could not find property %v", propertySlug), http.StatusUnprocessableEntity)
			return
		}

		propertyTypeSlug := chi.URLParam(r, "propertyTypeSlug")
		propertyType, ok := property.PropertyTypes[propertyTypeSlug]
		if !ok {
			Respond(ctx, w, fmt.Sprintf("Could not find property type %v", propertyTypeSlug), http.StatusUnprocessableEntity)
			return
		}

		result, err := definitions.Validate(r.Body, propertyType)
		if err != nil {
			if validationError, ok := err.(definitions.ValidationError); ok {
				e := Error{
					errors.New("Validation error"),
					http.StatusUnprocessableEntity,
					[]FieldError{},
				}

				for _, msg := range validationError.Fields {
					e.Fields = append(e.Fields, FieldError{
						Field: msg.Field(),
						Error: msg.Description(),
					})
				}

				Respond(ctx, w, e, http.StatusUnprocessableEntity)
				return
			}

			Respond(ctx, w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		featureId := chi.URLParam(r, "id")
		newObs := observations.NewObservation{
			Feature:      observations.Referenceable{ID: featureId},
			FeatureType:  observations.Referenceable{ID: ft.ID},
			Property:     observations.Referenceable{ID: property.ID},
			PropertyType: observations.Referenceable{ID: propertyType.ID},
			Process:      observations.Referenceable{ID: "urn:matterable:generic-upload"},

			Result: result,
		}

		id := uuid.New().String()
		now := time.Now()
		obs, err := observations.New(newObs, id, now)
		if err != nil {
			if vError, ok := err.(*observations.ValidationError); ok {
				Respond(ctx, w, vError, http.StatusUnprocessableEntity)
				return
			}
			RespondError(ctx, w, errors.Wrap(err, "creating new observation"))
			return
		}

		err = observations.Save(ctx, o.db, obs)
		if err != nil {
			Respond(ctx, w, "unable to save observation", http.StatusInternalServerError)
			return
		}

		Respond(ctx, w, obs, http.StatusOK)
		return
	}
}
