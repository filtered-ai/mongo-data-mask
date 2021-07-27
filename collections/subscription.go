package collections

import (
	"context"
	"log"

	"github.com/JRagone/mongo-data-gen/generators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Subscription struct {
	Id primitive.ObjectID `bson:"_id"`
}

// Populates `UserCollection` with `count` random users
func PopulateSubscriptions(db *mongo.Database, ctx context.Context, base generators.Base, count uint) {
	// Create collection
	collection := "subscriptionCollection"
	err := db.CreateCollection(ctx, collection)
	if err != nil {
		log.Fatal(err)
	}
	subscriptionCollection := db.Collection(collection)

	// Generate and insert data
	for i := uint(1); i <= count; i++ {
		insert := &Subscription{
			Id: primitive.NewObjectID(),
		}
		_, err = subscriptionCollection.InsertOne(ctx, insert)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Gets and returns a slice of all users
func GetSubscriptions(db *mongo.Database, ctx context.Context) []Subscription {
	subscriptionCollection := db.Collection("subscriptionCollection")
	cursor, err := subscriptionCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var subscriptions []Subscription
	for cursor.Next(ctx) {
		var subscription Subscription
		if err = cursor.Decode(&subscription); err != nil {
			log.Fatal(err)
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions
}
