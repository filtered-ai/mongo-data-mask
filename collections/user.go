package collections

import (
	"context"
	"log"

	"github.com/JRagone/mongo-data-gen/generators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id int32 `bson:"_id"`
}

// Populates `UserCollection` with `count` random users
func PopulateUsers(db *mongo.Database, ctx context.Context, base generators.Base, count uint) {
	// Create collection
	collection := "UserCollection"
	err := db.CreateCollection(ctx, collection)
	if err != nil {
		log.Fatal(err)
	}
	userCollection := db.Collection(collection)

	// Generate and insert data
	for i := uint(1); i <= count; i++ {
		insert := &User{
			Id: int32(i),
		}
		_, err = userCollection.InsertOne(ctx, insert)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Gets and returns a slice of all users
func GetUsers(db *mongo.Database, ctx context.Context) []User {
	userCollection := db.Collection("UserCollection")
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	var users []User
	for cursor.Next(ctx) {
		var user User
		if err = cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}
