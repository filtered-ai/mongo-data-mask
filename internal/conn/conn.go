package conn

import (
	"context"
	"log"
	"time"

	"github.com/filtered-ai/mongo-data-mask/internal/conn/comm"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/iv"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/org"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/peopledata"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/question"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Connection struct {
	db   *mongo.Database
	ctx  context.Context
	coll map[string]comm.Collectioner
}

func (c *Connection) New(client *mongo.Client, database string, ctx context.Context) {
	// Get database
	c.db = client.Database(database)
	c.ctx = ctx
	c.coll = make(map[string]comm.Collectioner)
	c.coll[org.Name] = org.New(c)
	c.coll[user.Name] = user.New(c)
	c.coll[iv.Name] = iv.New(c)
	c.coll[peopledata.Name] = peopledata.New(c)
	c.coll[question.Name] = question.New(c)
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
