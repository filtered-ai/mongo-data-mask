package collections

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Subscription struct {
	Id primitive.ObjectID `bson:"_id"`
}

// Populates `UserCollection` with `count` random users
func PopulateSubscriptions(preSubscriptions map[primitive.ObjectID]Subscription, db *mongo.Database, ctx context.Context, count uint) {
	// Create collection
	collection := "subscriptionCollection"
	err := db.CreateCollection(ctx, collection)
	if err != nil {
		log.Fatal(err)
	}
	subscriptionCollection := db.Collection(collection)

	var subscriptions []Subscription
	// Generate and insert data
	for objectId := range preSubscriptions {
		subscription := Subscription{
			Id: objectId,
		}
		subscriptions = append(subscriptions, subscription)
	}
	// Convert []Subscription to []interface{}
	var interfaceSubs []interface{}
	for _, subscription := range subscriptions {
		interfaceSubs = append(interfaceSubs, subscription)
	}
	_, err = subscriptionCollection.InsertMany(ctx, interfaceSubs)
	if err != nil {
		log.Fatal(err)
	}
}

// Populates random subscriptions containing preparatory subscription data
// in map
func PrepopulateSubscriptions(preSubscriptions map[primitive.ObjectID]Subscription, db *mongo.Database, ctx context.Context, count uint) {
	// Generate and insert partial data
	for i := int32(1); i <= int32(count); i++ {
		objectId := primitive.NewObjectID()
		preSubscription := Subscription{
			Id: objectId,
		}
		preSubscriptions[objectId] = preSubscription
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
