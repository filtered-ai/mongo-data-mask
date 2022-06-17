import { mask } from ".";
import { join } from "path";

jest.setTimeout(300000);

test('mask', async () => {
    await mask(5, join(__dirname, "..", "export"));
});
