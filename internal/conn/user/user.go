package user

import (
	"log"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/filtered-ai/mongo-data-mask/internal/conn/comm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	comm.Collection
}

type User struct {
	comm.Document `bson:"inline"`
	DisplayName   string       `bson:"displayName"`
	Email         string       `bson:"email"`
	PhoneNumber   string       `bson:"phoneNumber"`
	Portfolio     string       `bson:"portfolio"`
	Photo         string       `bson:"photo,omitempty"`
	ResumeURL     string       `bson:"resumeURL"`
	Title         string       `bson:"jobTitle"`
	Gender        string       `bson:"gender"`
	Experience    []Experience `bson:"experience"`
	Education     []Education  `bson:"education"`
	// Skills        []Skill      `bson:"skills"`
	ProfileURL string `bson:"profileUrl"`
}

type Experience struct {
	Id           primitive.ObjectID `bson:"_id"`
	Position     string             `bson:"position,omitempty"`
	Company      string             `bson:"company"`
	StartDate    string             `bson:"startDate,omitempty"`
	EndDate      string             `bson:"endDate,omitempty"`
	IsCurrentJob bool               `bson:"isCurrentJob,omitempty"`
}

type Education struct {
	Id        primitive.ObjectID `bson:"_id"`
	Major     string             `bson:"major,omitempty"`
	School    string             `bson:"school"`
	StartDate string             `bson:"startDate,omitempty"`
	EndDate   string             `bson:"endDate,omitempty"`
}

const Name = "UserCollection"

func New(conn comm.Connectioner) *Collection {
	return &Collection{
		comm.Collection{
			Conn: conn,
			Coll: conn.Db().Collection(Name),
		},
	}
}

func (c Collection) Mask(doc comm.Document) {
	if doc.Mixed["organization"] != nil && doc.Mixed["organization"].(int32) == comm.FilteredOrgId {
		return
	}
	_, err := c.Coll.UpdateByID(*c.Conn.Ctx(), doc.Mixed["_id"].(int32), bson.D{{
		Key: "$set", Value: &User{
			DisplayName: genDisplayName(10),
			Email:       gofakeit.Email(),
			PhoneNumber: genPhoneNumber(),
			Portfolio:   gofakeit.URL(),
			Photo:       genPhoto(true),
			ResumeURL:   "https://ucarecdn.com/41db4370-b26f-4c8c-b912-bcb96dcece65/",
			Title:       gofakeit.JobTitle(),
			Gender:      genGender(),
			Experience:  genExperience(),
			Education:   genEducation(),
			ProfileURL:  genProfileURL(),
		},
	}})
	if err != nil {
		log.Fatal(err)
	}
}

// Generate full name with 1 in `nameProvidedChance` odds of it being a
// username
func genDisplayName(nameProvidedChance int) string {
	randNum := rand.Intn(nameProvidedChance)
	if randNum%nameProvidedChance == 0 {
		return gofakeit.Username()
	}
	return gofakeit.Name()
}

// Generate phone number with 'tel:+1' prefix
func genPhoneNumber() string {
	return "tel:+1" + gofakeit.Phone()
}

// Generate a photo URL by providing a default photo
func genPhoto(photoUploaded bool) string {
	if photoUploaded {
		return "https://ucarecdn.com/ea61ff48-1605-4b48-950c-358232c4fc8d/"
	}
	return ""
}

// Generate a random gender
func genGender() string {
	gender := rand.Intn(3)
	if gender == 0 {
		return "male"
	} else if gender == 1 {
		return "female"
	}
	return ""
}

// Generate an array of random experiences
func genExperience() []Experience {
	experiences := []Experience{}
	numExperiences := rand.Intn(10)
	for i := 0; i < numExperiences; i++ {
		objectID := primitive.NewObjectID()
		experience := Experience{
			Id:      objectID,
			Company: gofakeit.Company(),
		}
		experiences = append(experiences, experience)
	}
	return experiences
}

// Generate an array of random educations
func genEducation() []Education {
	educations := []Education{}
	numEducations := rand.Intn(5)
	for i := 0; i < numEducations; i++ {
		objectID := primitive.NewObjectID()
		education := Education{
			Id:     objectID,
			School: gofakeit.State() + " University",
		}
		educations = append(educations, education)
	}
	return educations
}

// Generate a random profile URL
func genProfileURL() string {
	return "https://github.com/" + gofakeit.Username()
}
