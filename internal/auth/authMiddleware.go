package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/Giovanni-Tedesco/tmitracksapi/pkg/utilities"
)

func GeneralAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := VerifyJWT(w, r)

		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := VerifyJWT(w, r)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
			return
		}

		if claims.Role != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
