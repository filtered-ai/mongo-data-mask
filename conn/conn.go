package conn

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/JRagone/mongo-data-gen/comm"
	"github.com/JRagone/mongo-data-gen/generators"
	"github.com/JRagone/mongo-data-gen/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Connection struct {
	DB   *mongo.Database
	Ctx  context.Context
	Coll map[string]comm.Collectioner
}

func (c *Connection) New(client *mongo.Client, ctx context.Context) {
	// Get database
	c.DB = client.Database("testing")
	// Drop old collections
	collNames, _ := c.DB.ListCollectionNames(ctx, bson.M{})
	for _, collName := range collNames {
		if err := c.DB.Collection(collName).Drop(ctx); err != nil {
			log.Fatal(err)
		}
	}
	// Object holding seeded random generator
	seed := uint64(64)
	base := generators.Base{
		Seed: seed,
	}
	rand.Seed(int64(base.Seed))
	// c.Coll[orgName] = newOrg(1000)
	c.Coll[user.Name] = user.New(1000)
	// c.Coll[subscription.Name] = subscription.New(1000)
}

func (c *Connection) Database() *mongo.Database {
	return c.DB
}

func (c *Connection) Context() *context.Context {
	return &c.Ctx
}

func (c *Connection) Collection(collName string) comm.Collectioner {
	return c.Coll[collName]
}

func (c *Connection) Populate() {
	// Generate data
	start := time.Now()
	for _, coll := range c.Coll {
		coll.Prepopulate(c)
	}
	for _, coll := range c.Coll {
		coll.Populate(c)
	}
	elapsed := time.Since(start)
	log.Printf("Generation took %s", elapsed)
}
