import { byzcoin } from "@dedis/cothority";
import Docker from "dockerode";

describe("Module import Tests", () => {
    it("should import the module", () => {
        expect(byzcoin.ByzCoinRPC).toBeDefined();
    });
});

describe("Docker should be available", () => {
    it("should not yield an error when getting docker infos", async () => {
        const docker = new Docker();
        await expectAsync(docker.info()).toBeResolved();
    });
});
