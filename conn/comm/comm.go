package comm

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Connectioner interface {
	New(client *mongo.Client, ctx context.Context)
	DB() *mongo.Database
	Ctx() *context.Context
}

type Collectioner interface {
	Mask()
}

// Creates a collection
func CreateCollection(collection string, conn Connectioner) *mongo.Collection {
	err := conn.DB().CreateCollection(*conn.Ctx(), collection)
	if err != nil {
		log.Fatal(err)
	}
	return conn.DB().Collection(collection)
}
