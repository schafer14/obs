package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/volatiletech/authboss"
	"github.com/volatiletech/authboss/confirm"
	"github.com/volatiletech/authboss/expire"
	"github.com/volatiletech/authboss/lock"
	"github.com/volatiletech/authboss/remember"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Observations string
	People       string
}

func API(build string, db *mongo.Database, ab *authboss.Authboss, cfg Collections) chi.Router {
	r := chi.NewRouter()

	// Middleware
	// TODO: add no CSRF
	// TODO: add request throttling
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(ab.LoadClientStateMiddleware)
	r.Use(remember.Middleware(ab))

	// Define collections that will be used
	obsColl := db.Collection(cfg.Observations)
	personColl := db.Collection(cfg.People)
	// Define handlers
	authHandler := AuthHandler{ab}
	checkHandler := Check{build, db}
	oHandler := &ObservationHandler{obsColl}
	personHandler := &PersonHandler{personColl, obsColl}

	// ======================================
	// Protected routes
	// ======================================
	r.Group(func(r chi.Router) {
		r.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
		r.Use(lock.Middleware(ab))
		r.Use(confirm.Middleware(ab))
		r.Use(expire.Middleware(ab))

		// Information about currently logged in user
		r.MethodFunc("GET", "/v1/me", authHandler.CurrentlyLoggedIn)

		// Generic observation handler
		r.Route("/v1//observations", func(r chi.Router) {
			r.Get("/", oHandler.Get)
			r.Get("/{id}", oHandler.Find)
			r.Post("/", oHandler.Create)
		})

		// Person router
		r.Route("/v1/people", func(r chi.Router) {
			r.Post("/", personHandler.Create)
		})
	})

	// ======================================
	// Auth routes
	// ======================================
	r.Group(func(r chi.Router) {
		r.Use(authboss.ModuleListMiddleware(ab))
		r.Mount("/v1/auth", http.StripPrefix("/v1/auth", ab.Config.Core.Router))
	})

	// ======================================
	// Unprotected Routes
	// ======================================

	// Health Check
	r.Get("/health", checkHandler.Health)

	// Definitions route
	r.Get("/v1/definitions", GetDefinitions)

	return r
}
