package user

import (
	"context"
	"log"

	"github.com/JRagone/mongo-data-gen/comm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	count int32
	data  Data
}

type Data map[int32]User

type User struct {
	Id int32 `bson:"_id"`
}

const Name = "UserCollection"

func New(count int32) *Collection {
	return &Collection{
		count: count,
		data:  make(Data),
	}
}

func (c Collection) Count() int32 {
	return c.count
}

func (c Collection) Data() interface{} {
	return c.data
}

// Populates `UserCollection` with `count` random users
func (c Collection) Populate(conn comm.Connectioner) {
	// Create collection
	collection := comm.CreateCollection(Name, conn)

	var users []User
	// Generate and insert data
	for Id := range conn.Collection(Name).Data().(Data) {
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
	_, err := collection.InsertMany(*conn.Context(), interfaceUsers)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Collection) Prepopulate(conn comm.Connectioner) {
	// Generate and insert partial data
	for i := int32(1); i <= conn.Collection(Name).Count(); i++ {
		user := User{
			Id: i,
		}
		conn.Collection(Name).Data().(Data)[i] = user
	}
}

// Gets and returns a slice of all users
func Get(db *mongo.Database, ctx context.Context) []User {
	collection := db.Collection(Name)
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
