package question

import (
	"log"

	"github.com/filtered-ai/mongo-data-mask/internal/conn/comm"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/iv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	comm.Collection
}

type Question struct {
	comm.Document `bson:"inline"`
	VideoId       string      `bson:"videoID,omitempty"`
	ExpVideoId    string      `bson:"expVideoID,omitempty"`
	CodingHistory primitive.A `bson:"codingHistory,omitempty"`
}

const Name = "QuestionCollection"

func New(conn comm.Connectioner) *Collection {
	return &Collection{
		comm.Collection{
			Conn: conn,
			Coll: conn.Db().Collection(Name),
		},
	}
}

func (c Collection) Mask(doc comm.Document) {
	var videoId string
	var expVideoId string
	var codingHistory primitive.A
	// Check if Filtered uses this question in any interviews
	filteredQuestion := true
	id := doc.Mixed["_id"].(int32)
	query := bson.M{
		"questions":    id,
		"organization": comm.FilteredOrgId,
	}
	var interview iv.Interview
	if err := c.Conn.Db().Collection(iv.Name).FindOne(*c.Conn.Ctx(), query).Decode(&interview); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Fatal(err)
		}
		filteredQuestion = false
	}
	if filteredQuestion {
		return
	}
	if doc.Mixed["isVideoQuestion"].(bool) {
		videoId = "sample-candidate-video-02"
	} else {
		if _, ok := doc.Mixed["expVideoID"]; ok {
			expVideoId = "sample-candidate-video-expl-01"
		}
		if doc.Mixed["codingHistory"] != nil {
			codingHistory = doc.Mixed["codingHistory"].(primitive.A)
			for _, history := range codingHistory {
				if history.(primitive.M)["snapshot"] != nil {
					history.(primitive.M)["snapshot"] = "https://ucarecdn.com/ea61ff48-1605-4b48-950c-358232c4fc8d/"
				}
			}
		}
	}
	_, err := c.Coll.UpdateByID(*c.Conn.Ctx(), doc.Mixed["_id"].(int32), bson.D{{
		Key: "$set", Value: &Question{
			VideoId:       videoId,
			ExpVideoId:    expVideoId,
			CodingHistory: codingHistory,
		},
	}})
	if err != nil {
		log.Fatal(err)
	}
}
