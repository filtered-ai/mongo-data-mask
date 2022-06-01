package mongodatamask

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/filtered-ai/mongo-data-mask/internal/conn"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Mask(seed int64, mongoURI string) {
	// Seed fake data generator
	gofakeit.Seed(seed)
	// Connect to cluster
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
	//
	c := conn.Connection{}
	c.New(client, ctx)
	c.Mask()
}
