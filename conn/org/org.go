package org

import (
	"log"
	"math/rand"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
)

type Collection struct {
	comm.Collection
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
		comm.Collection{
			Conn: conn,
			Coll: conn.DB().Collection(Name),
		},
	}
}

func (c Collection) Mask(doc comm.Document) {
	companyName := gofakeit.Company()
	_, err := c.Coll.UpdateByID(*c.Conn.Ctx(), doc.Id, bson.D{{
		Key: "$set", Value: &Organization{
			Name:            companyName,
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
			SlackTeamName: companyName,
		},
	}})
	if err != nil {
		log.Fatal(err)
	}
}

// Generate `length` random domains
func genDomainWhiteList(length uint) []string {
	var domains []string
	for i := uint(1); i <= length; i++ {
		domains = append(domains, gofakeit.DomainName())
	}
	return domains
}

// Generate a random header title
func genHeaderTitle() string {
	return gofakeit.Word() + " " + gofakeit.Word() + " " + gofakeit.Word()
}

// Generate a random main title
func genMainTitle() string {
	const openingTags = "<h3 class=\"text-color-1B2B50\"><b>"
	const closingTags = "</b></h3>"
	return openingTags + gofakeit.HackerPhrase() + closingTags
}

// Generate a random subtitle
func genSubtitle() string {
	const openingTags = "<h6>"
	const closingTags = "</h6>"
	return openingTags + gofakeit.HackerPhrase() + closingTags
}
