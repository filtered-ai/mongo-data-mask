package integration_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/filtered-ai/mongo-data-mask/internal/conn"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/user"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

func TestNew(t *testing.T) {
	// Load env variables
	godotenv.Load()
	// Connect to cluster
	mongoUri := os.Getenv("MONGO_URI")
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

	c := conn.Connection{}
	c.New(client, cs.Database, ctx)
	_ = user.New(&c)
}
