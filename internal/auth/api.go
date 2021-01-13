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
	"github.com/Giovanni-Tedesco/tmitracksapi/utilities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	// "github.com/go-playground/validator/v10"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/options"
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

	err := json.NewDecoder(r.Body).Decode(&creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	Users := db.Collection("Users")

	insertDoc := User{Email: creds.Email, Password: string(hashedPassword)}

	results, err := Users.InsertOne(context.TODO(), insertDoc)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func SignIn(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

	// TODO: Request Input Validation

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	jwtKey := os.Getenv("JWT_KEY")

	var creds User

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&creds)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
		Email: creds.Email,
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

	// NOTE: This endpoint should be removed as it is a security vulnerability.

	claims, err := utilities.VerifyJWT(w, r)

	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}

	// json.NewEncoder(w).Encode(claims)
	fmt.Fprintf(w, "%v", claims.Email)

}
