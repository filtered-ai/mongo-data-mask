package org

import (
	"log"
	"math/rand"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	conn comm.Connectioner
	coll *mongo.Collection
}

type Organization struct {
	Id                  int32               `bson:"_id,omitempty"`
	Name                string              `bson:"name,omitempty"`
	Location            string              `bson:"location,omitempty"`
	LogoURL             string              `bson:"logoURL,omitempty"`
	Domain              string              `bson:"domain,omitempty"`
	DomainWhiteList     []string            `bson:"domainWhiteList,omitempty"`
	Industry            string              `bson:"industry,omitempty"`
	BrandColor          string              `bson:"brandColor,omitempty"`
	BrandBGImage        string              `bson:"brandBGImage,omitempty"`
	CustomCareerLanding CustomCareerLanding `bson:"customCareerLanding,omitempty"`
	SlackTeamName       string              `bson:"slackTeamName,omitempty"`
}

type CustomCareerLanding struct {
	HeaderTitle string `bson:"headerTitle,omitempty"`
	MainTitle   string `bson:"mainTitle,omitempty"`
	SubTitle    string `bson:"subTitle,omitempty"`
}

const Name = "OrganizationCollection"

func New(conn comm.Connectioner) *Collection {
	return &Collection{
		conn: conn,
		coll: conn.DB().Collection(Name),
	}
}

func (c Collection) Mask() {
	cursor, err := c.coll.Find(*c.conn.Ctx(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(*c.conn.Ctx())
	for cursor.Next(*c.conn.Ctx()) {
		var org Organization
		if err = cursor.Decode(&org); err != nil {
			log.Fatal(err)
		}
		_, err := c.coll.UpdateByID(*c.conn.Ctx(), org.Id, bson.D{{
			Key: "$set", Value: &Organization{
				Name:            gofakeit.Company(),
				Location:        gofakeit.City() + ", " + gofakeit.State(),
				LogoURL:         "https://ucarecdn.com/e6de77b8-aea5-4007-b42b-2f716e8734f8/",
				Domain:          gofakeit.DomainName(),
				DomainWhiteList: genDomainWhiteList(uint(rand.Intn(3))),
				Industry:        "Technology",
				BrandColor:      gofakeit.HexColor(),
				BrandBGImage:    "https://ucarecdn.com/ca12d558-8553-478a-a415-4b75c42cf6ed/lyft_logo.png",
				CustomCareerLanding: CustomCareerLanding{
					HeaderTitle: genHeaderTitle(),
					MainTitle:   genMainTitle(),
					SubTitle:    genSubtitle(),
				},
				SlackTeamName: gofakeit.Company(),
			},
		}})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Generates `length` random domains
func genDomainWhiteList(length uint) []string {
	var domains []string
	for i := uint(1); i <= length; i++ {
		domains = append(domains, gofakeit.DomainName())
	}
	return domains
}

// Generates a random header title
func genHeaderTitle() string {
	return gofakeit.Word() + " " + gofakeit.Word() + " " + gofakeit.Word()
}

// Generates a random main title
func genMainTitle() string {
	const openingTags = "<h3 class=\"text-color-1B2B50\"><b>"
	const closingTags = "</b></h3>"
	return openingTags + gofakeit.HackerPhrase() + closingTags
}

// Generates a random subtitle
func genSubtitle() string {
	const openingTags = "<h6>"
	const closingTags = "</h6>"
	return openingTags + gofakeit.HackerPhrase() + closingTags
}
