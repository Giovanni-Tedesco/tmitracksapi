package main

import (
	"context"
	"log"
	"os"
	"time"

	api "github.com/Giovanni-Tedesco/tmitracksapi/internal"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	a := &api.App{}

	err := godotenv.Load()
	uri := os.Getenv("MONGO_DB_URI")

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	a.Initialize(client)
	a.Run(":8080")
}
