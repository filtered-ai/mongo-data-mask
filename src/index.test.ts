import { mask } from ".";
import { join } from "path";
import { existsSync, mkdirSync, writeFileSync, readdirSync, rmSync } from "fs";
import { Masker } from "../src/mask";

jest.setTimeout(3000000);

beforeEach(() => {
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
});

describe('mask', () => {
  it('runs on data in the export dir', async () => {
    await mask(5, join(__dirname, "..", "export"));
  });
  it('runs on no data in the export dir', async () => {
    const exportDirPath = join(__dirname, "..", "export");
    readdirSync(exportDirPath).forEach((file) => {
      const filePath = join(exportDirPath, file);
      rmSync(filePath, { recursive: true });
    });
    await mask(5, join(__dirname, "..", "export"));
  });
});
