package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/mbop/internal/store"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type registationCreateRequest struct {
	UID *string `json:"uid,omitempty"`
}

func RegistrationCreateHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		do500(w, "failed to read body bytes: "+err.Error())
		return
	}

	var body registationCreateRequest
	err = json.Unmarshal(b, &body)
	if err != nil {
		do400(w, "failed to unmarshal body: "+err.Error())
		return
	}

	if body.UID == nil || *body.UID == "" {
		do400(w, "required parameter [uid] not found in body")
		return
	}

	id := identity.Get(r.Context())
	if !id.Identity.User.OrgAdmin {
		doError(w, "user must be org admin to register satellite", 403)
		return
	}

	gatewayCN, err := getCertCN(r.Header.Get(CertHeader))
	if err != nil {
		do400(w, err.Error())
		return
	}

	if gatewayCN != *body.UID {
		do400(w, "x-rh-certauth-cn does not match uid")
		return
	}

	db := store.GetStore()
	_, err = db.Find(id.Identity.OrgID, *body.UID)
	if err == nil {
		doError(w, "existing registration found", 409)
		return
	}

	_, err = db.Create(&store.Registration{
		OrgID: id.Identity.OrgID,
		UID:   *body.UID,
	})
	if err != nil {
		if errors.Is(err, store.ErrUIDAlreadyExists) {
			doError(w, "existing registration found", 409)
		} else {
			do500(w, "failed to create registration: "+err.Error())
		}
		return
	}

	sendJSONWithStatusCode(w, newResponse("Successfully registered"), 201)
}

func RegistrationDeleteHandler(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		do400(w, "invalid uid passed in path")
	}

	id := identity.Get(r.Context())
	if !id.Identity.User.OrgAdmin {
		doError(w, "user must be org admin to register satellite", 403)
		return
	}

	db := store.GetStore()

	err := db.Delete(id.Identity.OrgID, uid)
	if err != nil {
		if errors.Is(err, store.ErrRegistrationNotFound) {
			do404(w, err.Error())
		} else {
			do500(w, "error deleting registration: "+err.Error())
		}
		return
	}

	sendJSONWithStatusCode(w, newResponse("Successfully de-registered"), 204)
}
