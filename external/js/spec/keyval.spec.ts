import ByzCoinRPC from "@dedis/cothority/byzcoin/byzcoin-rpc";
import { Rule } from "@dedis/cothority/darc/rules";
import { EMPTY_BUFFER } from "@dedis/cothority/protobuf";

import { KVInstance } from "../src/keyval";
import { BLOCK_INTERVAL, ROSTER, SIGNER, startConodes } from "./support/conondes";

describe("KVInstance Tests", () => {
    const roster = ROSTER.slice(0, 4);

    beforeAll(async () => {
        await startConodes();
    });

    it("should spawn and update a keyValue instance", async () => {
        const darc = ByzCoinRPC.makeGenesisDarc([SIGNER], roster);
        darc.addIdentity("spawn:keyValue", SIGNER, Rule.OR);
        darc.addIdentity("invoke:keyValue.update", SIGNER, Rule.OR);

        const rpc = await ByzCoinRPC.newByzCoinRPC(roster, darc, BLOCK_INTERVAL);
        const value = Buffer.from("keyValue instance");
        const vi = await KVInstance.spawn(rpc, darc.getBaseID(), [SIGNER], "key", value);

        expect(vi.keyValueData.storage.length).toEqual(1);
        expect(vi.keyValueData.storage[0].key).toEqual("key");
        expect(vi.keyValueData.storage[0].value).toEqual(value);

        // Let's add a new value
        const value2 = Buffer.from("another keyValue instance");
        await vi.invokeUpdate([SIGNER], "key2", value2);
        const vi2 = await vi.update();

        // The returned variable (vi2) should contain the updated data
        expect(vi2.keyValueData.storage.length).toEqual(2);
        expect(vi2.keyValueData.storage[0].key).toEqual("key");
        expect(vi2.keyValueData.storage[0].value).toEqual(value);
        expect(vi2.keyValueData.storage[1].key).toEqual("key2");
        expect(vi2.keyValueData.storage[1].value).toEqual(value2);

        // vi should have been updated
        expect(vi.keyValueData.storage.length).toEqual(2);
        expect(vi.keyValueData.storage[0].key).toEqual("key");
        expect(vi.keyValueData.storage[0].value).toEqual(value);
        expect(vi.keyValueData.storage[1].key).toEqual("key2");
        expect(vi.keyValueData.storage[1].value).toEqual(value2);

        // Let's remove the first value
        await vi.invokeUpdate([SIGNER], "key", EMPTY_BUFFER);
        await vi.update();

        expect(vi.keyValueData.storage.length).toEqual(1);
        expect(vi.keyValueData.storage[0].key).toEqual("key2");
        expect(vi.keyValueData.storage[0].value).toEqual(value2);
    });
});
