package auth

import (

	// "log"
	"fmt"
	"net/http"

	. "github.com/Giovanni-Tedesco/tmitracksapi/utilities"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := VerifyJWT(w, r)

		if err != nil {
			fmt.Fprintf(w, "%v", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
