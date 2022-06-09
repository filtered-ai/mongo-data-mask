package iv

import (
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/comm"
	"go.mongodb.org/mongo-driver/bson"
)

type Collection struct {
	comm.Collection
}

type Interview struct {
	comm.Document  `bson:"inline"`
	CandidateEmail string `bson:"candidateEmail,omitempty"`
}

const Name = "InterviewCollection"

func New(conn comm.Connectioner) *Collection {
	return &Collection{
		comm.Collection{
			Conn: conn,
			Coll: conn.Db().Collection(Name),
		},
	}
}

func (c Collection) Mask(doc comm.Document) {
	if doc.Mixed["organization"].(int32) != comm.FilteredOrgId {
		return
	}
	_, err := c.Coll.UpdateByID(*c.Conn.Ctx(), doc.Mixed["_id"].(int32), bson.D{{
		Key: "$set", Value: &Interview{
			CandidateEmail: gofakeit.Email(),
		},
	}})
	if err != nil {
		log.Fatal(err)
	}
}
