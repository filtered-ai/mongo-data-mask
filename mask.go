package mongodatamask

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/filtered-ai/mongo-data-mask/internal/conn"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func Mask(seed int64, mongoUri string) {
	// Seed fake data generator
	gofakeit.Seed(seed)
	// Connect to database
	cs, err := connstring.ParseAndValidate(mongoUri)
	if err != nil {
		log.Fatal(err)
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
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
	// Mask data
	c := conn.Connection{}
	c.New(client, cs.Database, ctx)
	c.Mask()
}
