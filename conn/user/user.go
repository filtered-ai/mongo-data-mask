package user

import (
	"log"

	"github.com/JRagone/mongo-data-gen/comm"
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

	var users []interface{}
	// Generate and insert data
	for Id := range c.data {
		user := User{
			Id: Id,
		}
		users = append(users, user)
	}
	_, err := collection.InsertMany(*conn.Ctx(), users)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Collection) Prepopulate() {
	// Generate and insert partial data
	for i := int32(1); i <= c.count; i++ {
		user := User{
			Id: i,
		}
		c.data[i] = user
	}
}
