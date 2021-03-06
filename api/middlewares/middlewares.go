package middlewares

import (
	"errors"
	"net/http"

	"github.com/Mstuart712/rm/api/auth"
	"github.com/Mstuart712/rm/api/responses"
	"github.com/Mstuart712/rm/api/token"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	token.checkingImport()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
