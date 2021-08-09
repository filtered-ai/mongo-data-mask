package user

import (
	"log"
	"math/rand"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection struct {
	conn comm.Connectioner
	coll *mongo.Collection
}

type User struct {
	Id          int32        `bson:"_id,omitempty"`
	DisplayName string       `bson:"displayName"`
	Email       string       `bson:"email"`
	PhoneNumber string       `bson:"phoneNumber"`
	Portfolio   string       `bson:"portfolio"`
	Photo       string       `bson:"photo,omitempty"`
	ResumeURL   string       `bson:"resumeURL"`
	Title       string       `bson:"jobTitle"`
	Gender      string       `bson:"gender"`
	Experience  []Experience `bson:"experience"`
	Education   []Education  `bson:"education"`
	Skills      []Skill      `bson:"skills"`
	ProfileURL  string       `bson:"profileUrl"`
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

type Skill struct {
	Id                primitive.ObjectID `bson:"_id"`
	SkillName         string             `bson:"skillName"`
	IsVerified        bool               `bson:"isVerified"`
	IsFromDataService bool               `bson:"isFromDataService"`
}

const Name = "UserCollection"

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
		var org User
		if err = cursor.Decode(&org); err != nil {
			log.Fatal(err)
		}
		_, err := c.coll.UpdateByID(*c.conn.Ctx(), org.Id, bson.D{{
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
				Skills:      genSkills(),
				ProfileURL:  genProfileURL(),
			},
		}})
		if err != nil {
			log.Fatal(err)
		}
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

func genPhoneNumber() string {
	return "tel:+1" + gofakeit.Phone()
}

func genPhoto(photoUploaded bool) string {
	if photoUploaded {
		return "https://ucarecdn.com/ea61ff48-1605-4b48-950c-358232c4fc8d/"
	}
	return ""
}

func genGender() string {
	gender := rand.Intn(3)
	if gender == 0 {
		return "male"
	} else if gender == 1 {
		return "female"
	}
	return ""
}

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

func genSkills() []Skill {
	skills := []Skill{}
	numSkills := rand.Intn(50)
	for i := 0; i < numSkills; i++ {
		objectID := primitive.NewObjectID()
		skill := Skill{
			Id:                objectID,
			SkillName:         gofakeit.HackerVerb(),
			IsVerified:        gofakeit.Bool(),
			IsFromDataService: gofakeit.Bool(),
		}
		skills = append(skills, skill)
	}
	return skills
}

func genProfileURL() string {
	return "https://github.com/" + gofakeit.Username()
}
