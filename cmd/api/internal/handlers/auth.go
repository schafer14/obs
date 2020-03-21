package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/volatiletech/authboss"
)

type AuthHandler struct {
	ab *authboss.Authboss
}

// CurrentlyLoggedIn handles an http request for the currently logged in user.
func (a *AuthHandler) CurrentlyLoggedIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	u, err := a.ab.CurrentUser(r)
	if err != nil {
		RespondError(ctx, w, errors.Wrap(err, "fetching current user"))
	}

	Respond(ctx, w, u, http.StatusOK)
}
