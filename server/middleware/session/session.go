package session

import (
	"net/http"
)

func New() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionID := r.Header.Get("Mcp-Session-Id")
			if sessionID != "" {
				values := r.URL.Query()
				values.Add("sessionid", sessionID)
				r.URL.RawQuery = values.Encode()
			}

			next.ServeHTTP(w, r)
		})
	}
}
