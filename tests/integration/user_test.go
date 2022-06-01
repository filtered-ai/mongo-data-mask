package integration_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/JRagone/mongo-data-gen/conn"
	"github.com/JRagone/mongo-data-gen/conn/user"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNew(t *testing.T) {
	// Load env variables
	godotenv.Load()
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
	_ = user.New(&c)
}
