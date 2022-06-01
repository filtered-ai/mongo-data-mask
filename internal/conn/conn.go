package conn

import (
	"context"
	"log"
	"time"

	"github.com/JRagone/mongodatamask/internal/conn/comm"
	"github.com/JRagone/mongodatamask/internal/conn/org"
	"github.com/JRagone/mongodatamask/internal/conn/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Connection struct {
	db   *mongo.Database
	ctx  context.Context
	coll map[string]comm.Collectioner
}

func (c *Connection) New(client *mongo.Client, ctx context.Context) {
	// Get database
	c.db = client.Database("local_temp")
	c.ctx = ctx
	c.coll = make(map[string]comm.Collectioner)
	c.coll[org.Name] = org.New(c)
	c.coll[user.Name] = user.New(c)
}

func (c *Connection) Db() *mongo.Database {
	return c.db
}

func (c *Connection) Ctx() *context.Context {
	return &c.ctx
}

func (c *Connection) Mask() {
	for name, coll := range c.coll {
		start := time.Now()
		coll.IterateDocs(coll.Mask)
		elapsed := time.Since(start)
		log.Println(name, "masking took", elapsed)
	}
}
