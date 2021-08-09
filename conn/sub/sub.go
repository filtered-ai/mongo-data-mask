package sub

import (
	"log"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	count int32
	data  Data
}

type Data map[primitive.ObjectID]Subscription

type Subscription struct {
	Id primitive.ObjectID `bson:"_id"`
}

const Name = "subscriptionCollection"

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

	var subscriptions []interface{}
	// Generate and insert data
	for objectId := range c.data {
		subscription := Subscription{
			Id: objectId,
		}
		subscriptions = append(subscriptions, subscription)
	}
	_, err := collection.InsertMany(*conn.Ctx(), subscriptions)
	if err != nil {
		log.Fatal(err)
	}
}

// Populates random subscriptions containing preparatory subscription data
// in map
func (c Collection) Prepopulate() {
	// Generate and insert partial data
	for i := int32(1); i <= c.count; i++ {
		objectId := primitive.NewObjectID()
		subscription := Subscription{
			Id: objectId,
		}
		c.data[objectId] = subscription
	}
}
