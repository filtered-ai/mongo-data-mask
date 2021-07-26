package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	gofakeit.Seed(1)

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

	// Get database
	testDB := client.Database("testing")
	// Drop old collections
	collectionNames, _ := testDB.ListCollectionNames(ctx, bson.D{})
	for _, collectionName := range collectionNames {
		if err = testDB.Collection(collectionName).Drop(ctx); err != nil {
			log.Fatal(err)
		}
	}

	collection := "OrganizationCollection"
	err = testDB.CreateCollection(ctx, collection)
	if err != nil {
		log.Fatal(err)
	}
	moreTestingCol := testDB.Collection(collection)

	_, err = moreTestingCol.InsertOne(ctx, bson.D{
		{Key: "title", Value: gofakeit.Name()},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted one document into", collection+"!")
}
