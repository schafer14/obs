package handlers

import (
	"net/http"

	"github.com/schafer14/observations/internal/platform/database"
	"go.mongodb.org/mongo-driver/mongo"
)

// Check provides support for orchestration health checks.
type Check struct {
	build string
	db    *mongo.Database
}

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	health := struct {
		Version string `json:"version"`
		Status  string `json:"status"`
	}{
		Version: c.build,
	}

	// Check if the database is ready.
	if err := database.Check(ctx, c.db.Client()); err != nil {

		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "db not ready"
		Respond(ctx, w, health, http.StatusInternalServerError)
		return
	}

	health.Status = "ok"
	Respond(ctx, w, health, http.StatusOK)
	return
}
