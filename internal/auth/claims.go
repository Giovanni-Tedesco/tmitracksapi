package auth

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	email string `json:email`
	jwt.StandardClaims
}
