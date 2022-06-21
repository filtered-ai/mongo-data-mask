import { existsSync, mkdirSync, writeFileSync } from "fs";
import { join } from "path";
import { Masker } from "../src/mask";

export default function setup() {
  const exportDirPath = join(__dirname, "..", "export");
  if (!existsSync(exportDirPath)) {
    mkdirSync(exportDirPath);
  }
  for (const collName in Masker.maskableCollectionObjs) {
    const filePath = join(exportDirPath, `${collName}.json`);
    if (!existsSync(filePath)) {
      writeFileSync(filePath, "[{}]");
    }
  }
}
