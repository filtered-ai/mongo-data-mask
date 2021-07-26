package collections

import (
	"context"
	"log"
	"math/rand"

	"github.com/JRagone/mongo-data-gen/generators"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type preGen struct {
	isSuperOrg bool
	isSubOrg   bool
}

var orgSizes = []string{"1-100", "101-200", "201-1000", "1001-2000", "2001-4000", "4001+"}
var teams = []string{"Product", "Recruiting", "Sales", "Hiring team"}

// Generates `length` random domains
func genDomainWhiteList(length uint) []string {
	var domains []string
	for i := uint(1); i <= length; i++ {
		domains = append(domains, gofakeit.DomainName())
	}
	return domains
}

// Generates team size string
func genSize() string {
	index := rand.Intn(len(orgSizes))
	return orgSizes[index]
}

// Generates team name
func genTeam() string {
	index := rand.Intn(len(teams))
	return teams[index]
}

// Generates random slice of sub organizations
func genSubOrgs(base generators.Base, subOrgs []int32) []int32 {
	var selectedSubOrgs []int32
	numSubOrgs := rand.Intn(len(subOrgs))
	if len(subOrgs) == 0 {
		log.Fatal("There are zero sub orgs.")
		return selectedSubOrgs
	}

	rand.Shuffle(len(subOrgs), func(i, j int) {
		subOrgs[i], subOrgs[j] = subOrgs[j], subOrgs[i]
	})
	for i := 0; i <= numSubOrgs; i++ {
		selectedSubOrgs = append(selectedSubOrgs, subOrgs[i])
	}
	return selectedSubOrgs
}

// Populates `OrganizationCollection` with `count` random orgs
func PopulateOrgs(db *mongo.Database, ctx context.Context, base generators.Base, count uint) {
	collection := "OrganizationCollection"
	err := db.CreateCollection(ctx, collection)
	if err != nil {
		log.Fatal(err)
	}

	organizationCollection := db.Collection(collection)
	// boolGen := generators.NewBoolGenerator(base)

	// Generate pre-generation info
	preGenMap := make(map[uint]preGen)
	var subOrgs []int32
	for i := uint(1); i <= count; i++ {
		isSuperOrg := gofakeit.Bool()
		isSubOrg := !isSuperOrg
		if isSubOrg {
			subOrgs = append(subOrgs, int32(i))
		}
		newPreGen := preGen{isSuperOrg: isSuperOrg, isSubOrg: isSubOrg}
		preGenMap[i] = newPreGen
	}

	// Generate and insert data
	for i := uint(1); i <= count; i++ {
		preGen := preGenMap[i]
		insert := bson.D{
			{Key: "_id", Value: int32(i)},
			{Key: "isClosed", Value: gofakeit.Bool()},
			{Key: "name", Value: gofakeit.Company()},
			{Key: "location", Value: gofakeit.City() + ", " + gofakeit.State()},
			{Key: "logoURL", Value: "https://ucarecdn.com/c76800e5-939f-43f6-a6de-750a2692e31e/"},
			{Key: "domain", Value: gofakeit.DomainName()},
			{Key: "domainWhiteList", Value: genDomainWhiteList(3)},
			{Key: "size", Value: genSize()},
			{Key: "team", Value: genTeam()},
			{Key: "industry", Value: gofakeit.BuzzWord()},
			{Key: "isSuperOrg", Value: preGen.isSuperOrg},
			{Key: "isSubOrg", Value: preGen.isSubOrg},
			{Key: "showComment", Value: gofakeit.Bool()},
			{Key: "disableCommentBox", Value: gofakeit.Bool()},
			{Key: "isIVBranded", Value: gofakeit.Bool()},
			{Key: "brandColor", Value: gofakeit.HexColor()},
			{Key: "brandBGImage", Value: "https://ucarecdn.com/ca12d558-8553-478a-a415-4b75c42cf6ed/lyft_logo.png"},
			// {Key: "customCareerlanding", Value: }
		}
		// If org is a super org, add `subOrgs` field
		if preGen.isSuperOrg {
			insert = append(insert,
				bson.E{Key: "subOrgs", Value: genSubOrgs(base, subOrgs)})
		}
		_, err = organizationCollection.InsertOne(ctx, insert)
		if err != nil {
			log.Fatal(err)
		}
	}
}
