package main

import (
	"context"
	"log"
	"time"

	"github.com/JRagone/mongo-data-gen/comm"
	"github.com/JRagone/mongo-data-gen/conn"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	gofakeit.Seed(5)

	// Connect to cluster
	mongoURI := "mongodb://127.0.0.1:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Cancel context and disconnect last
	defer cancel()
	defer client.Disconnect(ctx)

	c := conn.Connection{Coll: make(map[string]comm.Collectioner), Ctx: ctx}
	c.New(client, ctx)
	c.Populate()
}
