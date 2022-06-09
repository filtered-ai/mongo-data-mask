package comm

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const FilteredOrgId = 50

type Connectioner interface {
	New(client *mongo.Client, database string, ctx context.Context)
	Db() *mongo.Database
	Ctx() *context.Context
}

type Collectioner interface {
	IterateDocs(handleDoc func(doc Document))
	Mask(doc Document)
}

type Collection struct {
	Conn Connectioner
	Coll *mongo.Collection
}

type Document struct {
	Mixed bson.M `bson:",inline"`
}

func (c Collection) IterateDocs(handleDoc func(doc Document)) {
	cursor, err := c.Coll.Find(*c.Conn.Ctx(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(*c.Conn.Ctx())
	var wg sync.WaitGroup
	for cursor.Next(*c.Conn.Ctx()) {
		var doc Document
		if err = cursor.Decode(&doc); err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleDoc(doc)
		}()
	}
	wg.Wait()
}

// Create a collection
func CreateCollection(collection string, conn Connectioner) *mongo.Collection {
	err := conn.Db().CreateCollection(*conn.Ctx(), collection)
	if err != nil {
		log.Fatal(err)
	}
	return conn.Db().Collection(collection)
}
