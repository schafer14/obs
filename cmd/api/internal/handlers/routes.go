package handlers

import (
	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	observationCollection = "observations"
)

func API(db *firestore.Client) chi.Router {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	coll := db.Collection(observationCollection)
	oHandler := &ObservationHandler{coll}
	r.Route("/v1/observations", func(r chi.Router) {
		r.Get("/", oHandler.Get)
		r.Get("/{id}", oHandler.Find)
		r.Post("/", oHandler.Create)
	})

	return r
}
