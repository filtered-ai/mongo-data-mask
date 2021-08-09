package orgusage

import (
	"log"
	"math"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	count int32
	data  Data
}

type Data map[primitive.ObjectID]OrgUsage

type OrgUsage struct {
	Id                            primitive.ObjectID `bson:"_id"`
	Organization                  int32              `bson:"organization"`
	IvSentNum                     int32              `bson:"ivSentNum"`
	IvSentNumMax                  float64            `bson:"ivSentNumMax"`
	IvCompletedNum                int32              `bson:"ivCompletedNum"`
	IvCompletedNumMax             float64            `bson:"ivCompletedNumMax"`
	MspIvSentNum                  int32              `bson:"mspIvSentNum"`
	MspIvCompletedNum             float64            `bson:"mspIvCompletedNum"`
	AdminNum                      int32              `bson:"adminNum"`
	AdminNumMax                   float64            `bson:"adminNumMax"`
	UserNum                       int32              `bson:"userNum"`
	UserNumMax                    float64            `bson:"userNumMax"`
	ExternalShareNum              int32              `bson:"externalShareNum"`
	ClientPortalNum               int32              `bson:"clientPortalNum"`
	ClientPortalNumMax            float64            `bson:"clientPortalNumMax"`
	InternalPortalNum             int32              `bson:"internalPortalNum"`
	InternalPortalNumMax          float64            `bson:"internalPortalNumMax"`
	JobDesNum                     int32              `bson:"jobDesNum"`
	JobDesNumMax                  float64            `bson:"jobDesNumMax"`
	LiveRoomNum                   int32              `bson:"liveRoomNum"`
	LiveRoomNumMax                float64            `bson:"liveRoomNumMax"`
	PublicLiveRoomNum             int32              `bson:"publicLiveRoomNum"`
	PublicLiveRoomNumMax          float64            `bson:"publicLiveRoomNumMax"`
	PublicLiveRoomCandyEnteredNum int32              `bson:"publicLiveRoomCandyEnteredNum"`
}

const Name = "orgUsageCollection"

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

func (c Collection) Populate(conn comm.Connectioner) {
	collection := comm.CreateCollection(Name, conn)
	var usages []interface{}
	for objectID, preUsage := range c.data {
		usage := OrgUsage{
			Id:                            objectID,
			Organization:                  preUsage.Organization,
			IvSentNum:                     0,
			IvSentNumMax:                  math.Inf(0),
			IvCompletedNum:                0,
			IvCompletedNumMax:             math.Inf(0),
			MspIvSentNum:                  0,
			MspIvCompletedNum:             math.Inf(0),
			AdminNum:                      0,
			AdminNumMax:                   math.Inf(0),
			UserNum:                       0,
			UserNumMax:                    math.Inf(0),
			ExternalShareNum:              0,
			ClientPortalNum:               0,
			ClientPortalNumMax:            math.Inf(0),
			InternalPortalNum:             0,
			InternalPortalNumMax:          math.Inf(0),
			JobDesNum:                     0,
			JobDesNumMax:                  math.Inf(0),
			LiveRoomNum:                   0,
			LiveRoomNumMax:                math.Inf(0),
			PublicLiveRoomNum:             0,
			PublicLiveRoomNumMax:          math.Inf(0),
			PublicLiveRoomCandyEnteredNum: 0,
		}
		usages = append(usages, usage)
	}
	_, err := collection.InsertMany(*conn.Ctx(), usages)
	if err != nil {
		log.Fatal(err)
	}
}

func (c Collection) Prepopulate() {
	for i := int32(1); i <= c.count; i++ {
		objectID := primitive.NewObjectID()
		usage := OrgUsage{
			Id:           objectID,
			Organization: i,
		}
		c.data[objectID] = usage
	}
}
