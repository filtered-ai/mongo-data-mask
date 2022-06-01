package user

import (
	"math"
	"strings"
	"testing"
)

// Test generation of username and fullname from genDisplayName()
func TestGenDisplayName(t *testing.T) {
	username := genDisplayName(1)
	splitName := strings.Split(username, " ")
	if len(splitName) != 1 {
		t.Errorf("Name %s is not a username", username)
	}

	fullName := genDisplayName(math.MaxInt)
	splitName = strings.Split(fullName, " ")
	if !(len(splitName) >= 2) {
		t.Errorf("Name %s is not a full name", fullName)
	}
}

// Test phone number generation from genPhoneNumber()
func TestGenPhoneNumber(t *testing.T) {
	fullNumber := genPhoneNumber()
	if !strings.Contains(fullNumber, "tel:+1") {
		t.Errorf("%s is not a phone number", fullNumber)
	}
	number := strings.Split(fullNumber, "tel:+1")[1]
	expectedLength := 10
	if len(number) != expectedLength {
		t.Errorf("%s phone number length %d is less than expected %d", fullNumber, len(number), expectedLength)
	}
}

// Test photo generation from genPhoto()
func TestGenPhoto(t *testing.T) {
	defaultPhoto := genPhoto(true)
	if len(defaultPhoto) == 0 {
		t.Errorf("No default photo given")
	}

	noPhoto := genPhoto(false)
	if len(noPhoto) > 0 {
		t.Errorf("Photo %s erroneously generated", noPhoto)
	}
}

// Test gender generation from genGender()
// Gender can be male, female, or unspecified
func TestGenGender(t *testing.T) {
	gender := genGender()
	validValues := []string{"male", "female", ""}
	validGender := false
	for _, validValue := range validValues {
		if gender == validValue {
			validGender = true
		}
	}
	if !validGender {
		t.Errorf("Gender %s is invalid", gender)
	}
}

// Test for company from genExperience()
func TestGenExperience(t *testing.T) {
	experiences := genExperience()
	for _, experience := range experiences {
		if experience.Company == "" {
			t.Errorf("Experience %+v\n specified no company", experience)
		}
	}
}

// Test for school from genEducation()
func TestGenEducation(t *testing.T) {
	educations := genEducation()
	for _, education := range educations {
		if education.School == "" {
			t.Errorf("Education %+v\n specified no school", education)
		}
	}
}

// Test for skill name from genSkills
func TestSkills(t *testing.T) {
	skills := genSkills()
	for _, skill := range skills {
		if skill.SkillName == "" {
			t.Errorf("Skill %+v\n specified no skill name", skill)
		}
	}
}

// Test that profile URL from genProfileURL is a URL
func TestGenProfileURL(t *testing.T) {
	profileURL := genProfileURL()
	if !strings.Contains(profileURL, "https://") {
		t.Errorf("Profile URL %s is not a URL", profileURL)
	}
}
