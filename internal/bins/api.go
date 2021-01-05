package bins

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bin struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Number int64              `json:"Number,omitempty"`
}

type User struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"Name,omitempty"`
	Age  int64              `json:"Age,omitempty"`
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func TestDb(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	collection := db.Collection("Bins")

	curr := collection.FindOne(context.TODO(), bson.M{})

	var bin Bin

	err := curr.Decode(&bin)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%v", bin.Number)

}

func GetUsers(db *mongo.Database, w http.ResponseWriter, r *http.Request) {

	collection := db.Collection("users")

	curr := collection.FindOne(context.TODO(), bson.M{})

	var user User

	err := curr.Decode(&user)

	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, "Something went wrong")
	}

	fmt.Fprintf(w, "%v", user.Name)

}

func GetAllUsers(db *mongo.Database, w http.ResponseWriter, r *http.Request) {
	collection := db.Collection("users")

	curr, err := collection.Find(context.TODO(), bson.M{})

	ctx := context.Background()

	var users []User

	for curr.Next(ctx) {
		var u User
		err = curr.Decode(&u)

		if err != nil {
			log.Fatal(err)
		}

		users = append(users, u)
	}

	// res, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Fprintf(w, "%s", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}
