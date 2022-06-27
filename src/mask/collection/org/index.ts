import { faker } from "@faker-js/faker";
import { randomInt } from "crypto";
import Collection from "..";

interface Organization {
  _id: number;
  name: string;
  location: string;
  logoURL: string;
  domain: string;
  domainWhiteList: string[];
  industry: string;
  brandColor: string;
  brandBGImage: string;
  customCareerLanding: CustomCareerLanding;
  slackTeamName: string;
}

interface CustomCareerLanding {
  headerTitle: string;
  mainTitle: string;
  subTitle: string;
}

/**
 * Organization collection masker.
 */
export default class OrganizationCollection extends Collection {
  public static readonly collName = "OrganizationCollection";

  /**
   * Construct a new organization collection object for masking.
   * @param dataFile Path to the JSON file to mask
   */
  constructor(dataFile: string) {
    super(dataFile);
  }

  protected genMaskedData(data: string): string {
    const org = <Organization>JSON.parse(data);
    if (org._id === this.filteredOrgId) {
      return JSON.stringify(org);
    }
    org.name = faker.company.companyName();
    org.location = `${faker.address.city()}, ${faker.address.stateAbbr()}`;
    org.logoURL = faker.image.business(943, 368, true);
    org.domain = faker.internet.domainName();
    org.domainWhiteList = Array.from({ length: randomInt(3) }).map(() =>
      faker.internet.domainName()
    );
    org.industry = faker.name.jobArea();
    org.brandColor = faker.color.rgb();
    org.brandBGImage = faker.image.business(943, 668, true);
    org.customCareerLanding = <CustomCareerLanding>{
      headerTitle: faker.company.catchPhrase(),
      mainTitle: this.genMainTitle(),
      subTitle: this.genSubtitle(),
    };
    org.slackTeamName = org.name;
    return JSON.stringify(org);
  }

  /**
   * Generate a main title for an organization's custom career landing page.
   * @returns Main title
   */
  private genMainTitle() {
    const openingTags = '<h3 class="text-color-1B2B50"><b>';
    const closingTags = "</b></h3>";
    return `${openingTags}${faker.company.catchPhrase()}${closingTags}`;
  }

  /**
   * Generate a subtitle for an organization's custom career landing page.
   * @returns Subtitle
   */
  private genSubtitle() {
    const openingTags = "<h6>";
    const closingTags = "</h6>";
    return `${openingTags}${faker.company.catchPhrase()}${closingTags}`;
  }
}
