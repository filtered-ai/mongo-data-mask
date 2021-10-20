package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/JRagone/mongo-data-gen/conn"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load env variables
	godotenv.Load()

	seed, err := strconv.ParseInt(os.Getenv("SEED"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	gofakeit.Seed(seed)

	// Connect to cluster
	mongoURI := os.Getenv("MONGO_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Cancel context and disconnect last
	defer cancel()
	defer client.Disconnect(ctx)

	c := conn.Connection{}
	c.New(client, ctx)
	c.Mask()
}
