package oauth

import (
	"net/http"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
)

// NewRoutes returns routes slice for the main app
func NewRoutes(cnf *config.Config, db *gorm.DB) []*rest.Route {
	return []*rest.Route{
		rest.Post("/oauth2/api/v1/tokens", func(w rest.ResponseWriter, r *rest.Request) {
			tokensHandler(w, r, cnf, db)
		}),
	}
}

// POST /oauth2/api/v1/tokens (handles all OAuth 2.0 grant types)
func tokensHandler(w rest.ResponseWriter, r *rest.Request, cnf *config.Config, db *gorm.DB) {
	// Check the grant type
	grantTypes := map[string]bool{
		"authorization_code": true,
		"implicit":           true,
		"password":           true,
		"client_credentials": true,
		"refresh_token":      true,
	}
	if !grantTypes[r.FormValue("grant_type")] {
		api.Error(w, "Invalid grant type", http.StatusBadRequest)
		return
	}

	// Authenticate the client
	client, err := authClient(r.Request, db)
	if err != nil {
		api.UnauthorizedError(w, err.Error())
		return
	}

	grants := map[string]func(){
		"authorization_code": func() { authorizationCodeGrant(w, r, cnf, db, client) },
		"implicit":           func() { implicitGrant(w, r, cnf, db, client) },
		"password":           func() { passwordGrant(w, r, cnf, db, client) },
		"client_credentials": func() { clientCredentialsGrant(w, r, cnf, db, client) },
		"refresh_token":      func() { refreshTokenGrant(w, r, cnf, db, client) },
	}
	grants[r.FormValue("grant_type")]()
}
