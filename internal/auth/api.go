package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	. "github.com/Giovanni-Tedesco/tmitracksapi/internal/entity"
	"github.com/Giovanni-Tedesco/tmitracksapi/pkg/utilities"
	"github.com/go-playground/validator/v10"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func SignUp(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	// TODO: Request input validation

	var creds User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&creds)

	v := validator.New()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = v.Struct(creds)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	Users := db.Collection("Users")

	insertDoc := User{Email: creds.Email, Password: string(hashedPassword),
		Role: creds.Role}

	results, err := Users.InsertOne(context.TODO(), insertDoc)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

type SignInBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func SignIn(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

	// TODO: Request Input Validation

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	jwtKey := os.Getenv("JWT_KEY")

	var creds SignInBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	v := validator.New()
	err = v.Struct(creds)

	if err != nil {
		log.Fatal(err)
	}

	collection := db.Collection("Users")

	var storedCreds User

	err = collection.FindOne(context.TODO(), bson.M{"email": creds.Email}).Decode(&storedCreds)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "No documents")
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationtime := time.Now().Add(time.Minute * 5)

	claims := &Claims{
		Email: storedCreds.Email,
		Role:  storedCreds.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationtime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationtime,
		HttpOnly: true,
		// Secure:   true,
	})

	w.WriteHeader(http.StatusOK)
}

func TestSomething(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

	claims, err := utilities.VerifyJWT(w, r)

	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	// json.NewEncoder(w).Encode(claims)
	fmt.Fprintf(w, "The user role is: %v", claims.Role)

}
