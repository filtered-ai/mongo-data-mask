package user

import (
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/JRagone/mongo-data-gen/conn/comm"
	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	count int32
	data  Data
}

type Data map[int32]User

type User struct {
	Id               int32        `bson:"_id"`
	IsClosed         bool         `bson:"isClosed"`
	DisplayName      string       `bson:"displayName"`
	Provider         string       `bson:"provider"`
	ProviderID       string       `bson:"providerId"`
	PJoinDate        time.Time    `bson:"pJoinDate"`
	PRepoCount       int32        `bson:"pRepoCount"`
	LoginCount       int32        `bson:"loginCount"`
	Team             interface{}  `bson:"team,omitempty"`
	Email            string       `bson:"email"`
	PhoneNumber      string       `bson:"phoneNumber"`
	Portfolio        string       `bson:"portfolio"`
	Photo            string       `bson:"photo,omitempty"`
	PhotoUploaded    bool         `bson:"photoUploaded"`
	NameProvided     bool         `bson:"nameProvided"`
	EmailProvided    bool         `bson:"emailProvided"`
	ResumeURL        string       `bson:"resumeURL"`
	Title            string       `bson:"jobTitle"`
	Gender           string       `bson:"gender"`
	Experience       []Experience `bson:"experience"`
	Education        []Education  `bson:"education"`
	Skills           []Skill      `bson:"skills"`
	ProfileURL       string       `bson:"profileUrl"`
	SignupSource     string       `bson:"signupSource"`
	IsEmailVerified  bool         `bson:"isEmailVerified"`
	Role             string       `bson:"role"`
	Position         string       `bson:"position"`
	IsAdmin          bool         `bson:"isAdmin"`
	IsDevAdmin       bool         `bson:"IsDevAdmin"`
	IsOrgsManager    bool         `bson:"IsOrgsManager"`
	ManagingOrgs     []int32      `bson:"managingOrgs"`
	NoManageOrgsPage bool         `bson:"noManageOrgsPage"`
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

var providers = [...]string{"filtered", "github", "slack", "linkedin", "bitbucket"}
var signupSources = [...]string{"candidate-app", "interview-room"}
var roles = [...]string{"recruiter", "candidate"}

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

// Populates `UserCollection` with `count` random users
func (c Collection) Populate(conn comm.Connectioner) {
	// Create collection
	collection := comm.CreateCollection(Name, conn)

	var users []interface{}
	// Generate and insert data
	for Id := range c.data {
		email := gofakeit.Email()
		nameProvided := gofakeit.Bool()
		photoUploaded := gofakeit.Bool()
		role := genRole()
		user := User{
			Id:          Id,
			IsClosed:    gofakeit.Bool(),
			DisplayName: genDisplayName(nameProvided, email),
			Provider:    genProvider(),
			ProviderID:  gofakeit.UUID(),
			PJoinDate:   genPJoinDate(),
			PRepoCount:  int32(gofakeit.Float32Range(0, 50)),
			LoginCount:  int32(gofakeit.Float32Range(0, 1000)),
			// TODO: Team
			Email:         email,
			PhoneNumber:   genPhoneNumber(),
			Portfolio:     gofakeit.URL(),
			Photo:         genPhoto(photoUploaded),
			PhotoUploaded: photoUploaded,
			NameProvided:  nameProvided,
			EmailProvided: true,
			ResumeURL:     "https://ucarecdn.com/41db4370-b26f-4c8c-b912-bcb96dcece65/",
			Title:         gofakeit.JobTitle(),
			Gender:        genGender(),
			Experience:    genExperience(),
			Education:     genEducation(),
			Skills:        genSkills(),
			ProfileURL:    genProfileURL(),
			SignupSource:  genSignupSource(),
			Role:          role,
			IsAdmin:       false,
			IsDevAdmin:    false,
			IsOrgsManager: gofakeit.Bool(),
		}
		if role == "recruiter" {
			user.Position = gofakeit.JobTitle()
			user.IsOrgsManager = gofakeit.Bool()
		}
		users = append(users, user)
	}
	_, err := collection.InsertMany(*conn.Ctx(), users)
	if err != nil {
		log.Fatal(err)
	}
}

func genDisplayName(nameProvided bool, email string) string {
	if nameProvided {
		return gofakeit.Name()
	}
	at := strings.LastIndex(email, "@")
	if at >= 0 {
		name := email[:at]
		return name
	} else {
		log.Fatal("Email ", email, " is invalid")
		return ""
	}
}

func genProvider() string {
	index := rand.Intn(len(providers))
	return providers[index]
}

func genPJoinDate() time.Time {
	newYork, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
	return gofakeit.DateRange(time.Date(2017, 1, 1, 1, 1, 1, 1, newYork), time.Now())
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

func genSignupSource() string {
	index := rand.Intn(len(signupSources))
	return signupSources[index]
}

func genRole() string {
	index := rand.Intn(len(roles))
	return roles[index]
}

func (c Collection) Prepopulate() {
	// Generate and insert partial data
	for i := int32(0); i < c.count; i++ {
		user := User{
			Id: i,
		}
		c.data[i] = user
	}
}
