package comm

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

type Connectioner interface {
	New(client *mongo.Client, ctx context.Context)
	Database() *mongo.Database
	Context() *context.Context
	Collection(string) Collectioner
}

type Collectioner interface {
	Count() int32
	Data() interface{}
	Populate(conn Connectioner)
	Prepopulate(conn Connectioner)
}

// Creates a collection
func CreateCollection(collection string, db Connectioner) *mongo.Collection {
	err := db.Database().CreateCollection(*db.Context(), collection)
	if err != nil {
		log.Fatal(err)
	}
	return db.Database().Collection(collection)
}
