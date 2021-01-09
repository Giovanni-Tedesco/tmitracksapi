package auth

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/go-playground/validator/v10"
	// "github.com/Giovanni-Tedesco/tmitracksapi/internal/entity"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func AuthMiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func TestSomething(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This works 🚀")
}
