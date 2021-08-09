package orgconfig

import (
	"log"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	count int32
	data  Data
}

type Data map[primitive.ObjectID]OrgConfig

type OrgConfig struct {
	Id                     primitive.ObjectID `bson:"_id"`
	EnableDefaultOverrides bool               `bson:"enableDefaultOverrides"`
	EnableCandySocial      bool               `bson:"enableCandySocial"`
	EnableAnonymousMode    bool               `bson:"enableAnonymousMode"`
	AnonymousModeWhiteList []string           `bson:"anonymousModeWhiteList"`
	EnableTrackTabs        bool               `bson:"enableTrackTabs"`
	EnableEEOC             bool               `bson:"enableEEOC"`
	EnableResume           bool               `bson:"enableResume"`
	EnableAvailability     bool               `bson:"enableAvailability"`
	EnableAllowRetakes     bool               `bson:"enableAllowRetakes"`
	EnableAutoReminder     bool               `bson:"enableAutoReminder"`
	EnableAutoRecord       bool               `bson:"enableAutoRecord"`
	NumberOfReminders      int32              `bson:"numberOfReminders"`
	DefaultTimeLimit       int32              `bson:"defaultTimeLimit"`
	DefaultScoreThreshold  int32              `bson:"defaultScoreThreshold"`
	EnableSlackNoti        bool               `bson:"enableSlackNoti"`
	IvRedirectUrl          string             `bson:"ivRedirectUrl"`
	IsSnapshotOff          bool               `bson:"isSnapshotOff"`
}

const Name = "orgConfigCollection"

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
	var configs []interface{}
	for Id := range c.data {
		config := OrgConfig{
			Id:                     Id,
			EnableDefaultOverrides: gofakeit.Bool(),
			EnableCandySocial:      gofakeit.Bool(),
			EnableAnonymousMode:    gofakeit.Bool(),
			AnonymousModeWhiteList: genAnonymousModeWhiteList(),
			EnableTrackTabs:        gofakeit.Bool(),
			EnableEEOC:             gofakeit.Bool(),
			EnableResume:           gofakeit.Bool(),
			EnableAvailability:     gofakeit.Bool(),
			EnableAllowRetakes:     gofakeit.Bool(),
			EnableAutoReminder:     gofakeit.Bool(),
			EnableAutoRecord:       gofakeit.Bool(),
			NumberOfReminders:      3,
			DefaultTimeLimit:       180,
			DefaultScoreThreshold:  70,
			EnableSlackNoti:        gofakeit.Bool(),
			IvRedirectUrl:          "",
			IsSnapshotOff:          gofakeit.Bool(),
		}
		configs = append(configs, config)
	}
	_, err := collection.InsertMany(*conn.Ctx(), configs)
	if err != nil {
		log.Fatal(err)
	}
}

func genAnonymousModeWhiteList() []string {
	return []string{}
}

func (c Collection) Prepopulate() {
	for i := int32(1); i <= c.count; i++ {
		objectId := primitive.NewObjectID()
		config := OrgConfig{
			Id: objectId,
		}
		c.data[objectId] = config
	}
}
