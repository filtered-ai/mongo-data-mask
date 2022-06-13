import { faker } from "@faker-js/faker";
import Collection from "..";

interface Interview {
  organization: number;
  candidateEmail: string;
  questions: number[];
}

/**
 * Interview collection masker.
 */
export default class InterviewCollection extends Collection {
  public static readonly collName = "InterviewCollection";
  public static filteredIvQuestions = new Set();

  /**
   * Construct a new interview collection object for masking.
   * @param dataFile Path to the JSON file to mask
   */
  constructor(dataFile: string) {
    super(dataFile);
  }

  protected genMaskedData(data: string): string {
    const iv = <Interview>JSON.parse(data);
    if (iv.organization === this.filteredOrgId) {
      InterviewCollection.filteredIvQuestions = new Set([
        ...iv.questions,
        ...InterviewCollection.filteredIvQuestions,
      ]);
      return JSON.stringify(iv);
    }
    iv.candidateEmail = faker.internet.email();
    return JSON.stringify(iv);
  }
}
