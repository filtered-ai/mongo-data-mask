package collections

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollectionName = "UserCollection"

type UserCollection struct {
	Count int32
	Data  UserData
}

type UserData map[int32]User

type User struct {
	Id int32 `bson:"_id"`
}

// Populates `UserCollection` with `count` random users
func PopulateUsers(generated Collections, db *mongo.Database, ctx context.Context) {
	// Create collection
	err := db.CreateCollection(ctx, userCollectionName)
	if err != nil {
		log.Fatal(err)
	}
	collection := db.Collection(userCollectionName)

	var users []User
	// Generate and insert data
	for Id := range generated.Users.Data {
		user := User{
			Id: Id,
		}
		users = append(users, user)
	}
	// Convert []User to []interface{}
	var interfaceUsers []interface{}
	for _, user := range users {
		interfaceUsers = append(interfaceUsers, user)
	}
	_, err = collection.InsertMany(ctx, interfaceUsers)
	if err != nil {
		log.Fatal(err)
	}
}

func PrepopulateUsers(generated Collections, db *mongo.Database, ctx context.Context) {
	// Generate and insert partial data
	for i := int32(1); i <= generated.Users.Count; i++ {
		user := User{
			Id: i,
		}
		generated.Users.Data[i] = user
	}
}

// Gets and returns a slice of all users
func GetUsers(db *mongo.Database, ctx context.Context) []User {
	collection := db.Collection(userCollectionName)
	cursor, err := collection.Find(ctx, bson.M{})
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
