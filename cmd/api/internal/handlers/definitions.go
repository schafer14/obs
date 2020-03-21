package handlers

import (
	"net/http"

	"github.com/schafer14/observations/internal/definitions"
)

func GetDefinitions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	Respond(ctx, w, definitions.Data, http.StatusOK)
	return
}
