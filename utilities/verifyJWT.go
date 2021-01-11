package utilities

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/Giovanni-Tedesco/tmitracksapi/internal/entity"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func VerifyJWT(w http.ResponseWriter, r *http.Request) (entity.Claims, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	jwtKey := os.Getenv("JWT_KEY")

	c, err := r.Cookie("token")

	if err != nil {
		if err != http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return entity.Claims{}, err
		}
		w.WriteHeader(http.StatusBadGateway)
		return entity.Claims{}, err
	}

	tokenString := c.Value

	claims := &entity.Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return entity.Claims{}, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return entity.Claims{}, err
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return entity.Claims{}, errors.New("Token is no longer valid")
	}

	return *claims, nil
}
