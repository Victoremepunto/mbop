package handlers

import (
	"errors"
	"net/http"

	"github.com/redhatinsights/mbop/internal/store"
)

type AuthV1Response struct {
	Mechanism string `json:"mechanism"`
	User      User   `json:"user"`
}

type User struct {
	AccountNumber  string   `json:"account_number"`
	ActivationKeys []string `json:"activation_keys"`
	AddressString  string   `json:"address_string"`
	DisplayName    string   `json:"display_name"`
	Email          string   `json:"email"`
	FirstName      string   `json:"first_name"`
	ID             int      `json:"id"`
	IsActive       bool     `json:"is_active"`
	IsInternal     bool     `json:"is_internal"`
	IsOrgAdmin     bool     `json:"is_org_admin"`
	LastName       string   `json:"last_name"`
	Locale         string   `json:"locale"`
	OrgID          string   `json:"org_id"`
	Type           string   `json:"type"`
	Username       string   `json:"username"`
}

func AuthV1Handler(w http.ResponseWriter, r *http.Request) {
	gatewayCN, err := getCertCN(r.Header.Get(CertHeader))
	if err != nil {
		do400(w, err.Error())
		return
	}

	db := store.GetStore()

	reg, err := db.FindByUID(gatewayCN)
	if err != nil {
		if errors.Is(err, store.ErrRegistrationNotFound) {
			doError(w, err.Error(), 401)
		} else {
			do500(w, "failed to search for registration: "+err.Error())
		}
		return
	}

	sendJSON(w, AuthV1Response{
		Mechanism: "cert",
		User: User{
			OrgID:       reg.OrgID,
			DisplayName: reg.OrgID,
			ID:          -1,
			IsActive:    true,
			IsOrgAdmin:  true,
			Type:        "system",
		},
	})
}
