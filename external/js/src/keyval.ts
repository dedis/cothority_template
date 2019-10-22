import Signer from "@dedis/cothority/darc/signer";
import ByzCoinRPC from "@dedis/cothority/byzcoin/byzcoin-rpc";
import ClientTransaction, { Argument, Instruction } from "@dedis/cothority/byzcoin/client-transaction";
import Instance, { InstanceID } from "@dedis/cothority/byzcoin/instance";
import { EMPTY_BUFFER, registerMessage } from "@dedis/cothority/protobuf";

import { addJSON } from '@dedis/cothority/protobuf';
import models from "./protobuf/models.json";

import { Message, Properties } from "protobufjs/light";

/**
 * This class offers a wrapper around the cothority library to spawn and update
 * a keyValue contract.
 */
export class KVInstance extends Instance {
    static readonly contractID = "keyValue";
    keyValueData: KeyValueData;

    static async spawn(bc: ByzCoinRPC, darcID: InstanceID, signers: Signer[], 
        key: string, value: Buffer): Promise<KVInstance> {
        
        const arg = new Argument({
            name: key, 
            value: value
        })
        const inst = Instruction.createSpawn(darcID, KVInstance.contractID, [arg]);

        const ctx = ClientTransaction.make(bc.getProtocolVersion(), inst);
        await ctx.updateCounters(bc, [signers])
        ctx.signWith([signers]);

        await bc.sendTransactionAndWait(ctx, 10);
        return KVInstance.fromByzcoin(bc, ctx.instructions[0].deriveId(), 10);
    }

    static create(bc: ByzCoinRPC, valueID: InstanceID, darcID: InstanceID, 
        data: Buffer): KVInstance {

        return new KVInstance(bc, new Instance({
            contractID: KVInstance.contractID,
            darcID,
            data: data,
            id: valueID,
        }));
    }

    static async fromByzcoin(bc: ByzCoinRPC, iid: InstanceID, 
        waitMatch: number = 0, interval: number = 1000):Promise<KVInstance> {

        return new KVInstance(bc, await Instance.fromByzcoin(bc, iid, 
            waitMatch, interval));
    }

    constructor(private rpc: ByzCoinRPC, inst: Instance) {
        super(inst);
        if (inst.contractID.toString() !== KVInstance.contractID) {
            throw new Error(`mismatch contract name: ${inst.contractID} vs 
            ${KVInstance.contractID}`);
        }

        this.keyValueData = KeyValueData.decode(inst.data)
    }


    async invokeUpdate(signers: Signer[], key: string, value: Buffer, 
        wait?: number): Promise<void> {

        const inst = Instruction.createInvoke(
            this.id,
            KVInstance.contractID,
            "update",
            [
                new Argument({
                    name: key, 
                    value: value
                }),
            ],
        );

        const ctx = ClientTransaction.make(this.rpc.getProtocolVersion(), inst);
        await ctx.updateCounters(this.rpc, [signers])
        ctx.signWith([signers]);

        await this.rpc.sendTransactionAndWait(ctx, 10);
    }

    async update(): Promise<KVInstance> {
        const p = await this.rpc.getProofFromLatest(this.id);
        if (!p.exists(this.id)) {
            throw new Error("fail to get a matching proof");
        }
        this.keyValueData = KeyValueData.decode(p.value)
        return this;
    }

    toString(): string {
        var res: string = "";
        res += "KeyValueInstance:\n";
        res += this.keyValueData.toString()
        return res;
    }
}

/**
 * This class declares a message to encode/decode the content of a keyValue
 * instance, namely a KeyValueData struct. It follows the definition of the
 * keyvalue.proto.
 */
export class KeyValueData extends Message<KeyValueData> {
    storage: KeyValue[];

    constructor(props?: Properties<KeyValueData>) {
        super(props);

        this.storage = this.storage || [];
    }

    static register() {
        registerMessage("KeyValueData", KeyValueData, KeyValue);
    }

    toString(): string {
        var res: string = "";
        res += "keyValueData:";
        var i = 1;
        this.storage.forEach(element => {
            res += "\n- KV " + i + ":"
            res += "\n-- key: " + element.key
            res += "\n-- value: " + element.value.toString("utf8")
            i++;
        });
        return res;
    }
}

export class KeyValue extends Message<KeyValue> {
    key: string;
    value: Buffer;

    constructor(props?: Properties<KeyValue>) {
        super(props);

        this.key = this.key || "default";
        this.value = Buffer.from(this.value || EMPTY_BUFFER)
    }

    static register() {
        registerMessage("KeyValue", KeyValue);
    }
}

// Here we update the model with our models.json containing the definition of
// the keyValueData and we register our messages classes, so that protobuf can
// encore and decode it.
addJSON(models)
KeyValue.register()
KeyValueData.register()