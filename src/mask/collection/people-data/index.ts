import Collection from "..";
import UserCollection from "../user";

interface PeopleData {
  submitter: number;
  data: {
    profiles: JSON[];
  };
}

/**
 * People data collection masker.
 */
export default class PeopleDataCollection extends Collection {
  public static readonly collName = "peopleDataCollection";

  /**
   * Construct a new people data collection object for masking.
   * @param dataFile Path to the JSON file to mask
   */
  constructor(dataFile: string) {
    const dependentCollections = [UserCollection.collName];
    super(dataFile, dependentCollections);
  }

  protected genMaskedData(data: string): string {
    const peopleData = <PeopleData>JSON.parse(data);
    // Only mask if the submitter is not in the Filtered organization
    if (UserCollection.filteredUserIds[peopleData.submitter]) {
      return JSON.stringify(peopleData);
    }
    if (peopleData.data) {
      peopleData.data.profiles = [];
    }
    return JSON.stringify(peopleData);
  }
}
