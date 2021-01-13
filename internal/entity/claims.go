package entity

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}
