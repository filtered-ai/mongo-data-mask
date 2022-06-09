package peopledata

import (
	"log"

	"github.com/filtered-ai/mongo-data-mask/internal/conn/comm"
	myUser "github.com/filtered-ai/mongo-data-mask/internal/conn/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	comm.Collection
}

type PeopleData struct {
	comm.Document `bson:"inline"`
	Data          bson.M `bson:"data,omitempty"`
}

const Name = "peopleDataCollection"

func New(conn comm.Connectioner) *Collection {
	return &Collection{
		comm.Collection{
			Conn: conn,
			Coll: conn.Db().Collection(Name),
		},
	}
}

func (c Collection) Mask(doc comm.Document) {
	var data bson.M
	if doc.Mixed["data"] != nil {
		data = genData(c.Conn, doc)
	}
	_, err := c.Coll.UpdateByID(*c.Conn.Ctx(), doc.Mixed["_id"].(primitive.ObjectID), bson.D{{
		Key: "$set", Value: &PeopleData{
			Data: data,
		},
	}})
	if err != nil {
		log.Fatal(err)
	}
}

func genData(conn comm.Connectioner, doc comm.Document) bson.M {
	var user myUser.User
	if err := conn.Db().Collection(myUser.Name).FindOne(*conn.Ctx(), bson.M{"_id": doc.Mixed["submitter"], "organization": comm.FilteredOrgId}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return bson.M{
				"profiles": []bson.M{},
			}
		}
		log.Fatal(err)
	}
	return doc.Mixed["data"].(bson.M)
}
