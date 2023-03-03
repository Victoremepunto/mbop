package handlers

import (
	"net/http"

	"github.com/redhatinsights/mbop/internal/config"
	"github.com/redhatinsights/mbop/internal/service/catchall"
)

// instantiate on startup - only if needed.
var mbop = catchall.MakeNewMBOPServer()

func CatchAll(w http.ResponseWriter, r *http.Request) {
	if config.Get().DisableCatchall {
		do404(w, "not found")
		return
	}

	mbop.MainHandler(w, r)
}
