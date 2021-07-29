package collections

import (
	"context"
	"log"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
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

const orgCollectionName = "OrganizationCollection"
const maxSubOrgs = 5

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
func genSubOrgs(subOrgs []int32) []int32 {
	// Select number of sub orgs to return
	var selectedSubOrgs []int32
	numSubOrgs := rand.Intn(len(subOrgs))
	if len(subOrgs) == 0 {
		log.Fatal("There are zero sub orgs.")
		return selectedSubOrgs
	}
	if numSubOrgs > maxSubOrgs {
		numSubOrgs = maxSubOrgs
	}

	// Shuffle sub orgs so all calls don't get the same sub orgs
	rand.Shuffle(len(subOrgs), func(i, j int) {
		subOrgs[i], subOrgs[j] = subOrgs[j], subOrgs[i]
	})

	// Select sub orgs
	for i := 0; i < numSubOrgs; i++ {
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
func genCreator(users map[int32]User) int32 {
	// Return a random user
	index := rand.Int31n(int32(len(users)))
	return users[index].Id
}

// Generates a random subscription, which is a reference
func genSubscription(subs map[primitive.ObjectID]Subscription) primitive.ObjectID {
	// Return a random subscription
	keys := make([]primitive.ObjectID, 0, len(subs))
	for objectID := range subs {
		keys = append(keys, objectID)
	}
	index := rand.Intn(len(keys))
	return subs[keys[index]].Id
}

// Compose Organization object
func composeOrgs(preOrgs map[int32]Organization, users map[int32]User, subscriptions map[primitive.ObjectID]Subscription) []Organization {
	var orgs []Organization
	// Get pre-generation info
	subOrgs := GetSubOrgs(preOrgs)
	// Generate and insert data
	for Id, preOrg := range preOrgs {
		org := Organization{
			Id:                Id,
			IsClosed:          gofakeit.Bool(),
			Name:              gofakeit.Company(),
			Location:          gofakeit.City() + ", " + gofakeit.State(),
			LogoURL:           "https://ucarecdn.com/c76800e5-939f-43f6-a6de-750a2692e31e/",
			Domain:            gofakeit.DomainName(),
			DomainWhiteList:   genDomainWhiteList(3),
			Size:              genSize(),
			Team:              genTeam(),
			Industry:          gofakeit.BuzzWord(),
			IsSuperOrg:        preOrg.IsSuperOrg,
			IsSubOrg:          preOrg.IsSubOrg,
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
			Creator:      genCreator(users),
			Subscription: genSubscription(subscriptions),
		}
		// If org is a super org, add `subOrgs` field
		if org.IsSuperOrg {
			org.SubOrgs = genSubOrgs(subOrgs)
		}
		orgs = append(orgs, org)
	}
	return orgs
}

// Populates the database with `count` random orgs
func PopulateOrgs(preOrgs map[int32]Organization, users map[int32]User, subscriptions map[primitive.ObjectID]Subscription, db *mongo.Database, ctx context.Context, count uint) {
	// Create the collection
	collection := CreateCollection(orgCollectionName, db, ctx)
	// Get pre-generation info
	orgs := composeOrgs(preOrgs, users, subscriptions)
	// Convert []Organization to []interface{}
	var interfaceOrgs []interface{}
	for _, org := range orgs {
		interfaceOrgs = append(interfaceOrgs, org)
	}
	_, err := collection.InsertMany(ctx, interfaceOrgs)
	if err != nil {
		log.Fatal(err)
	}
}

// Populates random orgs containing preparatory organization data in map
func PrepopulateOrgs(preOrgs map[int32]Organization, db *mongo.Database, ctx context.Context, count uint) {
	// Generate and insert partial data
	for i := int32(1); i <= int32(count); i++ {
		isSuperOrg := gofakeit.Bool()
		isSubOrg := !isSuperOrg
		preOrg := Organization{
			Id:         i,
			IsSuperOrg: isSuperOrg,
			IsSubOrg:   isSubOrg,
		}
		preOrgs[i] = preOrg
	}
}

// Gets an organization by `Id`
func GetOrg(Id int32, db *mongo.Database, ctx context.Context) Organization {
	var org Organization
	collection := db.Collection(orgCollectionName)
	err := collection.FindOne(ctx, bson.M{"_id": Id}).Decode(&org)
	if err != nil {
		log.Fatal(err)
	}
	return org
}

// Gets a slice of all orgs that are sub orgs
func GetSubOrgs(preOrgs map[int32]Organization) []int32 {
	var subOrgs []int32
	for _, preOrg := range preOrgs {
		if preOrg.IsSubOrg {
			subOrgs = append(subOrgs, preOrg.Id)
		}
	}
	return subOrgs
}
