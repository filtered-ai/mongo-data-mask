import { faker } from "@faker-js/faker";
import fs from "fs";
import path from "path";
import Logger from "../logger";
import Collection, { collectionEvents } from "./collection";
import InterviewCollection from "./collection/iv";
import OrganizationCollection from "./collection/org";
import PeopleDataCollection from "./collection/people-data";
import QuestionCollection from "./collection/question";
import UserCollection from "./collection/user";

interface MaskableCollectionObjs {
  [key: string]: {
    class: new (dataFile: string) => Collection;
    masked: boolean;
  };
}

/**
 * Data masker
 */
export class Masker {
  // Add new maskable collections here
  public static readonly maskableCollectionObjs: MaskableCollectionObjs = {
    [UserCollection.collName]: { class: UserCollection, masked: false },
    [OrganizationCollection.collName]: {
      class: OrganizationCollection,
      masked: false,
    },
    [InterviewCollection.collName]: {
      class: InterviewCollection,
      masked: false,
    },
    [PeopleDataCollection.collName]: {
      class: PeopleDataCollection,
      masked: false,
    },
    [QuestionCollection.collName]: {
      class: QuestionCollection,
      masked: false,
    },
  };
  private readonly dataDir: string;

  /**
   * Construct a new Masker.
   * @param dataDir Directory contianing exported MongoDB JSON files
   */
  constructor(dataDir: string) {
    this.dataDir = dataDir;
  }

  /**
   * Mask all maskable collections.
   */
  public async mask() {
    // Create all maskable collections if their JSON file exists
    for (const collName in Masker.maskableCollectionObjs) {
      const filename = collName;
      const filePath = path.join(this.dataDir, `${filename}.json`);
      if (!fs.existsSync(filePath)) {
        delete Masker.maskableCollectionObjs[collName];
        continue;
      }
      new Masker.maskableCollectionObjs[collName].class(filePath);
    }
    // Send event to start masking all collections
    collectionEvents.emit("startMask");
    try {
      await this.allCollectionsMasked();
    } catch (err) {
      if (typeof err === "string") Logger.error(err);
    }
  }

  /**
   * Mark collections as masked when they finish masking, and resolve
   * when all collections are masked.
   * @returns Resolved when all collections are masked.
   */
  private allCollectionsMasked(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (!Object.keys(Masker.maskableCollectionObjs).length) {
        resolve();
      }
      collectionEvents.on("collectionMasked", (collName: string) => {
        if (!(collName in Masker.maskableCollectionObjs)) {
          reject(`[ERROR] Unknown collection ${collName} masked.`);
        }
        Masker.maskableCollectionObjs[collName].masked = true;
        // If all collections are masked, resolve
        if (
          Object.values(Masker.maskableCollectionObjs).every(
            (coll) => coll.masked
          )
        ) {
          resolve();
        }
      });
    });
  }
}

/**
 * Mask data in `dataDir`.
 * @param seed Seed for generating fake data
 * @param dataDir Directory with exported JSON collection files
 */
export async function mask(seed: number, dataDir: string) {
  faker.mersenne.seed(seed);
  const masker = new Masker(dataDir);
  await masker.mask();
}
