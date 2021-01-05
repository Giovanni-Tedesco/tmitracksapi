package main

import (
	"context"
	"time"

	api "github.com/Giovanni-Tedesco/tmitracksapi/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	a := &api.App{}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	a.Initialize(client)
	a.Run(":8080")
}
