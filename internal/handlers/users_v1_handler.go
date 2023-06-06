package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/service/ocm"
	"github.com/redhatinsights/mbop/internal/service/ocm/keycloak"

	"github.com/redhatinsights/mbop/internal/models"
)

func UsersV1Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch config.Get().UsersModule {
	case amsModule, mockModule:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			do500(w, "failed to read request body: "+err.Error())
			return
		}
		defer r.Body.Close()

		var usernames models.UserBody
		err = json.Unmarshal(body, &usernames)
		if err != nil {
			do400(w, "failed to parse request body: "+err.Error()+", request must include 'users': [] ")
			return
		}

		q, err := initV1UserQuery(r)
		if err != nil {
			do400(w, err.Error())
			return
		}

		// Create new SDK client
		client, err := ocm.NewOcmClient()
		if err != nil {
			do400(w, err.Error())
			return
		}

		err = client.InitSdkConnection(ctx)
		if err != nil {
			do500(w, "Can't build sdk connection: "+err.Error())
			return
		}

		u, err := client.GetUsers(usernames, q)
		if err != nil {
			do500(w, "Cant Retrieve Accounts: "+err.Error())
			return
		}

		// For each user see if it's an org_admin
		isOrgAdmin, err := client.GetOrgAdmin(u.Users)
		if err != nil {
			do500(w, "Cant Retrieve Role Bindings: "+err.Error())
			return
		}

		for i, user := range u.Users {
			response, ok := isOrgAdmin[user.ID]
			if ok {
				u.Users[i].IsOrgAdmin = response.IsOrgAdmin
			} else {
				user.IsOrgAdmin = false
			}
		}

		// Close SDK Connection
		client.CloseSdkConnection()

		sendJSON(w, u.Users)
	case keycloakModule:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			do500(w, "failed to read request body: "+err.Error())
			return
		}
		defer r.Body.Close()

		var usernames models.UserBody
		err = json.Unmarshal(body, &usernames)
		if err != nil {
			do400(w, "failed to parse request body: "+err.Error()+", request must include 'users': [] ")
			return
		}

		q, err := initV1UserQuery(r)
		if err != nil {
			do400(w, err.Error())
			return
		}

		client, err := keycloak.NewKeyCloakClient()
		if err != nil {
			do400(w, err.Error())
			return
		}

		err = client.InitKeycloakConnection()
		if err != nil {
			do500(w, "Can't build keycloak connection: "+err.Error())
			return
		}

		token, err := client.GetAccessToken()
		if err != nil {
			do500(w, "Can't fetch keycloak token: "+err.Error())
			return
		}

		u, err := client.GetUsers(token, usernames, q)
		if err != nil {
			do500(w, "Cant Retrieve Keycloak Accounts: "+err.Error())
			return
		}

		sendJSON(w, u.Users)
		return
	default:
		// mbop server instance injected somewhere
		// pass right through to the current handler
		CatchAll(w, r)
	}
}
