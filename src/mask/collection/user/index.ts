import { faker } from "@faker-js/faker";
import { randomInt } from "crypto";
import { ObjectId } from "mongodb";
import Collection from "..";

interface User {
  _id: number;
  organization: number;
  displayName: string;
  email: string;
  phoneNumber: string;
  portfolio: string;
  photo: string;
  resumeURL: string;
  jobTitle: string;
  gender: string;
  experience: Experience[];
  education: Education[];
  profileURL: string;
}

interface Experience {
  id: ObjectId;
  position: string;
  company: string;
  startDate: string;
  endDate: string;
  isCurrentJob: boolean;
}

interface Education {
  id: ObjectId;
  major: string;
  school: string;
  startDate: string;
  endDate: string;
}

/**
 * User collection masker.
 */
export default class UserCollection extends Collection {
  public static readonly collName = "UserCollection";
  public static filteredUserIds: { [key: number]: boolean } = {};

  /**
   * Construct a new user collection object for masking.
   * @param dataFile Path to the JSON file to mask
   */
  constructor(dataFile: string) {
    super(dataFile);
  }

  protected genMaskedData(data: string): string {
    const user = <User>JSON.parse(data);
    if (user.organization === this.filteredOrgId) {
      UserCollection.filteredUserIds[user._id] = true;
      return JSON.stringify(user);
    }
    user.displayName = this.genDisplayName(10);
    user.email = this.genEmail(user.displayName);
    user.phoneNumber = faker.phone.number("tel:+1##########");
    user.portfolio = faker.internet.url();
    user.photo = faker.image.avatar();
    user.resumeURL =
      "https://ucarecdn.com/41db4370-b26f-4c8c-b912-bcb96dcece65/";
    user.jobTitle = faker.name.jobTitle();
    user.gender = this.genGender();
    user.experience = this.maskExperience(user.experience);
    user.education = this.maskEducation(user.education);
    user.profileURL = `https://github.com/${faker.internet.userName()}`;
    return JSON.stringify(user);
  }

  /**
   * Generate a display name with a one in `nameNotProvidedChance`
   * of it not being provided, i.e. the display name is a full name.
   * @param nameNotProvidedChance Chance of a name being provided
   * @returns Display name
   */
  private genDisplayName(nameNotProvidedChance: number) {
    const nameProvided = randomInt(nameNotProvidedChance);
    if (nameProvided == 0) {
      return faker.internet.userName();
    }
    return faker.name.findName();
  }

  /**
   * Generate an email address for a user. The display name is used in the
   * email.
   * @param displayName Display name
   * @returns Email
   */
  private genEmail(displayName: string) {
    const splitName = displayName.split(" ");
    if (splitName.length > 1) {
      return faker.internet.email(splitName[0], splitName[1]);
    }
    return faker.internet.email(displayName);
  }

  /**
   * Generate a random gender.
   * @returns Gender
   */
  private genGender() {
    const gender = randomInt(3);
    if (gender == 0) {
      return "female";
    }
    if (gender == 1) {
      return "male";
    }
    return "";
  }

  /**
   * Mask the position and company of a user's experiences.
   * @param experiencesIn User's experiences
   * @returns Masked experiences
   */
  private maskExperience(experiencesIn: Experience[]) {
    const experiences = experiencesIn;
    if (experiences) {
      for (const [i] of experiences.entries()) {
        experiences[i].position = faker.name.jobTitle();
        experiences[i].company = faker.company.companyName();
      }
    }
    return experiences;
  }

  /**
   * Mask the major and school of a user's education.
   * @param educationsIn User's education
   * @returns Masked education
   */
  private maskEducation(educationsIn: Education[]) {
    const educations = educationsIn;
    if (educations) {
      for (const [i] of educations.entries()) {
        educations[i].major = faker.name.jobArea();
        educations[i].school = `${faker.address.state()} University`;
      }
    }
    return educations;
  }
}
