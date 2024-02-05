package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redhatinsights/mbop/internal/store"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type allowlistCreateRequest struct {
	IP string `json:"ip"`
}

type allowListResponse struct {
	IP        string    `json:"ip"`
	OrgID     string    `json:"org_id"`
	CreatedAt time.Time `json:"created_at"`
}

func AllowlistCreateHandler(w http.ResponseWriter, r *http.Request) {
	id := identity.Get(r.Context())
	if !id.Identity.User.OrgAdmin {
		doError(w, "user must be org admin to add addresses to allowlist", 403)
		return
	}

	var createReq allowlistCreateRequest
	err := json.NewDecoder(r.Body).Decode(&createReq)
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, "invalid json in body - only expected key is [ip]")
		return
	}

	db := store.GetStore()

	err = db.AllowAddress(&store.Address{IP: createReq.IP, OrgID: id.Identity.OrgID})
	if err != nil {
		w.WriteHeader(500)
		err = fmt.Errorf("error storing address: %w", err)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(201)
}

func AllowlistDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := identity.Get(r.Context())
	if !id.Identity.User.OrgAdmin {
		doError(w, "user must be org admin to add addresses to allowlist", 403)
		return
	}

	ip := chi.URLParam(r, "address")
	if ip == "" {
		w.WriteHeader(400)
		io.WriteString(w, "need address in path in the form `/v1/allowlist/{address}")
		return
	}

	db := store.GetStore()

	err := db.DenyAddress(&store.Address{IP: ip})
	if err != nil {
		if errors.Is(err, store.ErrAddressNotAllowListed) {
			w.WriteHeader(404)
			io.WriteString(w, "ip not allowlisted")
			return
		}

		w.WriteHeader(500)
		err = fmt.Errorf("error deleting addressaddress: %w", err)
		io.WriteString(w, err.Error())
		return
	}

	w.WriteHeader(204)
}

func AllowlistListHandler(w http.ResponseWriter, r *http.Request) {
	id := identity.Get(r.Context())
	if !id.Identity.User.OrgAdmin {
		doError(w, "user must be org admin to add addresses to allowlist", 403)
		return
	}

	db := store.GetStore()

	runtime.Breakpoint()
	addrs, err := db.AllowedAddresses(id.Identity.OrgID)
	if err != nil {
		w.WriteHeader(500)
		err = fmt.Errorf("error listing addresses: %w", err)
		io.WriteString(w, err.Error())
		return
	}

	out := make([]allowListResponse, len(addrs))
	for i, addr := range addrs {
		out[i] = allowListResponse{
			IP:        addr.IP,
			OrgID:     addr.OrgID,
			CreatedAt: addr.CreatedAt,
		}
	}

	json.NewEncoder(w).Encode(out)
}
