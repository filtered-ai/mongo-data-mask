import fs from "fs";
import JSONStream from "JSONStream";
import { EventEmitter } from "events";
import path from "path";
import Logger from "../../logger";

export const collectionEvents = new EventEmitter();

/**
 * All maskable collections extend this class.
 * @abstract
 */
export default abstract class Collection {
  protected readonly filteredOrgId = 50;
  protected readonly dependentCollections: { [key: string]: boolean };
  private readonly dataFile: string;
  private readonly handleCollMaskedBoundFunc: (collName: string) => void;

  /**
   * Construct a new collection object for masking.
   * @param dataFile Path to the JSON file to mask
   */
  constructor(dataFile: string, dependentCollections?: string[]) {
    this.dataFile = dataFile;
    // Add dependent collections
    this.dependentCollections = {};
    if (dependentCollections) {
      for (const coll of dependentCollections) {
        // False means the dependent collection has not been masked
        this.dependentCollections[coll] = false;
      }
    }
    // Start masking on the `startMask` event
    collectionEvents.once("startMask", this.handleStartMask.bind(this));
    this.handleCollMaskedBoundFunc = this.handleCollectionMasked.bind(this);
  }

  /**
   * Handle the `startMask` collection event. If there are dependent
   * collections, listen for them to be masked.
   */
  private handleStartMask() {
    if (Object.keys(this.dependentCollections).length > 0) {
      collectionEvents.on("collectionMasked", this.handleCollMaskedBoundFunc);
    } else {
      this.mask();
    }
  }

  /**
   * Start masking once all dependent collections are masked.
   * @param collName The name of the collection that was masked
   */
  private handleCollectionMasked(collName: string): void {
    if (!this.dependentCollections[collName]) {
      this.dependentCollections[collName] = true;
      // If all dependent collections are masked, start masking
      if (Object.values(this.dependentCollections).every(Boolean)) {
        collectionEvents.off(
          "collectionMasked",
          this.handleCollMaskedBoundFunc
        );
        this.mask();
      }
    }
  }

  /**
   * Stream and mask content from the data file.
   */
  public mask() {
    const outFileParsedPath = path.parse(this.dataFile);
    const collName = path.parse(this.dataFile).name;
    Logger.progress(`Masking ${collName}...`);
    const start = Date.now();
    // Use JSONstream.stringify to write a JSON array
    const jsonWriter = JSONStream.stringify("[", ",", "]");
    const outFile = path.join(
      outFileParsedPath.dir,
      `${outFileParsedPath.name}.masked${outFileParsedPath.ext}`
    );
    const writeStream = fs.createWriteStream(outFile);
    // Pipe the JSON array to the write stream for the masked file
    jsonWriter.pipe(writeStream);
    fs.createReadStream(this.dataFile)
      .pipe(JSONStream.parse("*"))
      .on("data", (data) => {
        jsonWriter.write(JSON.parse(this.genMaskedData(JSON.stringify(data))));
      })
      .on("end", () => {
        jsonWriter.end();
        fs.renameSync(outFile, this.dataFile);
        const stop = Date.now();
        Logger.success(
          `${collName} masked in ${(stop - start) / 1000} seconds.`
        );
        collectionEvents.emit("collectionMasked", collName);
      });
  }

  /**
   * Generate masked data for a single document.
   * @param data JSON string of doc to mask
   * @returns JSON string of masked doc
   */
  protected abstract genMaskedData(data: string): string;
}
