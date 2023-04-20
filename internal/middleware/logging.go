package middleware

import (
	"fmt"
	"net/http"

	l "github.com/redhatinsights/mbop/internal/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Log.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL))

		next.ServeHTTP(w, r)
	})
}
