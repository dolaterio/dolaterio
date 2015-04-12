package api

import "net/http"

// Authenticate gets the user from the request token
func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// token := req.Header.Get("Authentication")
		h.ServeHTTP(res, req)
	})
}
