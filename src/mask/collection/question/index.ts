import Collection from "..";
import InterviewCollection from "../iv";

interface Question {
  _id: number;
  organization: number;
  isVideoQuestion: boolean;
  videoID: string;
  expVideoID: string;
  codingHistory: CodingHistory[];
}

interface CodingHistory {
  snapshot: string;
}

/**
 * Question collection masker.
 */
export default class QuestionCollection extends Collection {
  public static readonly collName = "QuestionCollection";

  /**
   * Construct a new question collection object for masking.
   * @param dataFile Path to the JSON file to mask
   */
  constructor(dataFile: string) {
    const dependentCollections = [InterviewCollection.collName];
    super(dataFile, dependentCollections);
  }

  protected genMaskedData(data: string): string {
    const question = <Question>JSON.parse(data);
    if (InterviewCollection.filteredIvQuestions.has(question._id)) {
      return JSON.stringify(question);
    }
    if (question.isVideoQuestion) {
      question.videoID = "sample-candidate-video-02";
      return JSON.stringify(question);
    }
    if (question.expVideoID) {
      question.expVideoID = "sample-candidate-video-expl-01";
    }
    if (question.codingHistory) {
      for (const [i, codingHistory] of question.codingHistory.entries()) {
        if (codingHistory.snapshot) {
          question.codingHistory[i].snapshot =
            "https://ucarecdn.com/ea61ff48-1605-4b48-950c-358232c4fc8d/";
        }
      }
    }
    return JSON.stringify(question);
  }
}
