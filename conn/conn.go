package conn

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"github.com/JRagone/mongo-data-gen/conn/org"
	"github.com/JRagone/mongo-data-gen/conn/orgconfig"
	"github.com/JRagone/mongo-data-gen/conn/orgusage"
	"github.com/JRagone/mongo-data-gen/conn/sub"
	"github.com/JRagone/mongo-data-gen/conn/user"
	"github.com/JRagone/mongo-data-gen/gens"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Connection struct {
	db   *mongo.Database
	ctx  context.Context
	coll map[string]comm.Collectioner
}

func (c *Connection) New(client *mongo.Client, ctx context.Context) {
	// Get database
	c.db = client.Database("testing")
	c.ctx = ctx
	c.coll = make(map[string]comm.Collectioner)
	// Drop old collections
	collNames, _ := c.db.ListCollectionNames(ctx, bson.M{})
	for _, collName := range collNames {
		if err := c.db.Collection(collName).Drop(ctx); err != nil {
			log.Fatal(err)
		}
	}
	// Object holding seeded random generator
	seed := uint64(64)
	base := gens.Base{
		Seed: seed,
	}
	rand.Seed(int64(base.Seed))
	orgCount := int32(1000)
	c.coll[org.Name] = org.New(orgCount)
	c.coll[user.Name] = user.New(1000)
	c.coll[sub.Name] = sub.New(1000)
	c.coll[orgconfig.Name] = orgconfig.New(orgCount)
	c.coll[orgusage.Name] = orgusage.New(orgCount)
}

func (c *Connection) DB() *mongo.Database {
	return c.db
}

func (c *Connection) Ctx() *context.Context {
	return &c.ctx
}

func (c *Connection) Coll(collName string) comm.Collectioner {
	return c.coll[collName]
}

func (c *Connection) Populate() {
	// Generate data
	start := time.Now()
	for _, coll := range c.coll {
		coll.Prepopulate()
	}
	for _, coll := range c.coll {
		coll.Populate(c)
	}
	elapsed := time.Since(start)
	log.Printf("Generation took %s", elapsed)
}
