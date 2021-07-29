package collections

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// Creates a collection
func CreateCollection(collection string, db *mongo.Database, ctx context.Context) *mongo.Collection {
	err := db.CreateCollection(ctx, collection)
	if err != nil {
		log.Fatal(err)
	}
	return db.Collection(collection)
}
