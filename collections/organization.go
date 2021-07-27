package collections

import (
	"context"
	"log"
	"math/rand"

	"github.com/JRagone/mongo-data-gen/generators"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomCareerLanding struct {
	HeaderTitle string `bson:"headerTitle"`
	MainTitle   string `bson:"mainTitle"`
	SubTitle    string `bson:"subTitle"`
}

type Organization struct {
	Id                  int32               `bson:"_id"`
	IsClosed            bool                `bson:"isClosed"`
	Name                string              `bson:"name"`
	Location            string              `bson:"location"`
	LogoURL             string              `bson:"logoURL"`
	Domain              string              `bson:"domain"`
	DomainWhiteList     []string            `bson:"domainWhiteList"`
	Size                string              `bson:"size"`
	Team                string              `bson:"team"`
	Industry            string              `bson:"industry"`
	IsSuperOrg          bool                `bson:"isSuperOrg"`
	IsSubOrg            bool                `bson:"isSubOrg"`
	ShowComment         bool                `bson:"showComment"`
	DisableCommentBox   bool                `bson:"disableCommentBox"`
	IsIVBranded         bool                `bson:"isIVBranded"`
	BrandColor          string              `bson:"brandColor"`
	BrandBGImage        string              `bson:"brandBGImage"`
	CustomCareerLanding CustomCareerLanding `bson:"customCareerLanding"`
	Creator             int32               `bson:"creator"`
	SubOrgs             []int32             `bson:"subOrgs,omitempty"`
	Subscription        primitive.ObjectID  `bson:"subscription"`
}

type preGen struct {
	isSuperOrg bool
	isSubOrg   bool
}

var orgSizes = [...]string{"1-100", "101-200", "201-1000", "1001-2000", "2001-4000", "4001+"}
var teams = [...]string{"Product", "Recruiting", "Sales", "Hiring team"}

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

// Generates a random creator, which is a reference Id to a user
func genCreator(users *[]User) int32 {
	// Return a random user
	index := rand.Intn(len(*users))
	return (*users)[index].Id
}

// Generates a random subscription, which is a reference
func genSubscription(subs *[]Subscription) primitive.ObjectID {
	// Return a random subscription
	index := rand.Intn(len(*subs))
	return (*subs)[index].Id
}

// Populates `OrganizationCollection` with `count` random orgs
func PopulateOrgs(db *mongo.Database, ctx context.Context, base generators.Base, count uint) {
	// Create collection
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
	users := GetUsers(db, ctx)
	subscriptions := GetSubscriptions(db, ctx)

	// Generate and insert data
	for i := uint(1); i <= count; i++ {
		preGen := preGenMap[i]
		insert := &Organization{
			Id:                int32(i),
			IsClosed:          gofakeit.Bool(),
			Name:              gofakeit.Company(),
			Location:          gofakeit.City() + ", " + gofakeit.State(),
			LogoURL:           "https://ucarecdn.com/c76800e5-939f-43f6-a6de-750a2692e31e/",
			Domain:            gofakeit.DomainName(),
			DomainWhiteList:   genDomainWhiteList(3),
			Size:              genSize(),
			Team:              genTeam(),
			Industry:          gofakeit.BuzzWord(),
			IsSuperOrg:        preGen.isSuperOrg,
			IsSubOrg:          preGen.isSubOrg,
			ShowComment:       gofakeit.Bool(),
			DisableCommentBox: gofakeit.Bool(),
			IsIVBranded:       gofakeit.Bool(),
			BrandColor:        gofakeit.HexColor(),
			BrandBGImage:      "https://ucarecdn.com/ca12d558-8553-478a-a415-4b75c42cf6ed/lyft_logo.png",
			CustomCareerLanding: CustomCareerLanding{
				HeaderTitle: genHeaderTitle(),
				MainTitle:   genMainTitle(),
				SubTitle:    genSubtitle(),
			},
			Creator:      genCreator(&users),
			Subscription: genSubscription(&subscriptions),
		}
		// If org is a super org, add `subOrgs` field
		if preGen.isSuperOrg {
			insert.SubOrgs = genSubOrgs(base, subOrgs)
		}
		_, err = organizationCollection.InsertOne(ctx, insert)
		if err != nil {
			log.Fatal(err)
		}
	}
}
