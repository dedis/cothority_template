import * as Cothority from "@dedis/cothority";
import {KVInstance, KeyValueData} from "./keyval";
import { EMPTY_BUFFER } from '@dedis/cothority/protobuf';

export {
    Cothority
};

// ----------------------------------------------------------------------------
// The following functions are called from the view and responsible for parsing
// the arguments and then calling the appropriate method on the Handler class.

export function initRoster(e: Event) {
    try {
        var handler = Handler.getInstance();
        const fr = new FileReader();
        fr.onload = function(evt) {
            handler.LoadRoster(evt.target.result.toString())
        }
        const target = e.target as HTMLInputElement;
        const file: File = (target.files as FileList)[0];
        fr.readAsText(file);
        // Needed so that we can reload a same file multiple times
        target.value = "";
    } catch (e) {
        Handler.prependLog("failed to initialize the roster: " + e)
    }
}

export function displayStatus() {
    try {
        var r: string;
        if ((r = Handler.checkRoster()) != "") {
            Handler.prependLog(r)
            return
        }
        var handler = Handler.getInstance();
        const div = document.createElement("div")
        if (Handler.roster === undefined) {
            Handler.prependLog("handler has not been initialized");
            return
        }
        Handler.roster.list.forEach(element => {
            var p = document.createElement("p")
            p.innerText = element.address + ", " + element.description
            div.appendChild(p);
        });
        Handler.prependLog(div);
    } catch (e) {
        Handler.prependLog("failed to display status: " + e)
    }
}

export function getDarc(scidID: string) {
    try {
        var r: string = Handler.checkRoster();
        if (r != "") {
            Handler.prependLog(r)
            return
        }
        const scidHolder = document.getElementById(scidID) as HTMLInputElement;
        const scidStr =  scidHolder.value;
        if (scidStr == "") {
            Handler.prependLog("please enter a skipchain id")
            return
        }
        Handler.getInstance().SetDarc(Buffer.from(hexStringToByte(scidStr)))
    } catch (e) {
        Handler.prependLog("failed to set DARC: " + e)
    }
}

export function loadSigner(iID: string) {
    try {
        var r: string = Handler.checkRoster() || Handler.checkDarc();
        if (r!= "") {
            Handler.prependLog(r)
            return
        }
        const signerHolder = document.getElementById(iID) as HTMLInputElement
        const signerStr = signerHolder.value
        if (signerStr == "") {
            Handler.prependLog("please provide a signer")
            return
        }

        Handler.getInstance().SetSigner(Buffer.from(hexStringToByte(signerStr)))
    } catch (e) {
        Handler.prependLog("failed to set the signer: " + e)
    }
}

export function addRule(rID: string) {
    try {
        var r: string = Handler.checkRoster() || Handler.checkDarc() || Handler.checkSigner();
        if (r != "") {
            Handler.prependLog(r)
            return
        }
        const ruleHolder = document.getElementById(rID) as HTMLInputElement
        const ruleStr = ruleHolder.value
        if (ruleStr == "") {
            Handler.prependLog("please provide a rule")
            return
        }
        Handler.getInstance().AddRule(ruleStr);
    } catch (e) {
        Handler.prependLog("failed to add rule on DARC: " + e)
    }
}

export function spawnKV(keyID: string, valueID: string) {
    try {
        var r: string = Handler.checkRoster() || Handler.checkDarc() || Handler.checkSigner();
        if (r != "") {
            Handler.prependLog(r)
            return
        }
        const keyHolder = document.getElementById(keyID) as HTMLInputElement
        const valueHolder = document.getElementById(valueID) as HTMLInputElement
        const keyStr = keyHolder.value
        if (keyStr == "") {
            Handler.prependLog("please provide a key")
            return
        }
        const valueStr = valueHolder.value
        if (valueStr == "") {
            Handler.prependLog("please provide a value. Empty value is not allowed in spawn")
            return
        }
        Handler.getInstance().SpawnKV(keyStr, valueStr);
    } catch (e) {
        Handler.prependLog("failed to spawn keyValue instance: " + e)
    }
}

export function printKV(instIDID: string) {
    try {
        var r: string = Handler.checkRoster() || Handler.checkDarc() || Handler.checkSigner();
        if (r != "") {
            Handler.prependLog(r)
            return
        }
        const instIDHolder = document.getElementById(instIDID) as HTMLInputElement
        const instIDStr = instIDHolder.value
        if (instIDStr == "") {
            Handler.prependLog("please provide an instance id")
            return
        }

        Handler.getInstance().PrintKV(Buffer.from(instIDStr, "hex"));
    } catch (e) {
        Handler.prependLog("failed to print keyValue instance: " + e)
    }
}

export function addKVpair(instIDID: string, keyID: string, valueID: string) {
    try { 
        var r: string = Handler.checkRoster() || Handler.checkDarc() || Handler.checkSigner();
        if (r != "") {
            Handler.prependLog(r)
            return
        }
        const instIDHolder = document.getElementById(instIDID) as HTMLInputElement
        const keyHolder = document.getElementById(keyID) as HTMLInputElement
        const valueHolder = document.getElementById(valueID) as HTMLInputElement
        const instIDStr = instIDHolder.value
        if (instIDStr == "") {
            Handler.prependLog("please provide an instance id")
            return
        }
        const keyStr = keyHolder.value
        if (keyStr == "") {
            Handler.prependLog("please provide a key")
            return
        }
        const valueStr = valueHolder.value
        if (valueStr == "") {
            var c = confirm("This action will remove the key '"+ keyStr + ". Are you sure?");
            if (!c) return
        }
        Handler.getInstance().AddKVpair(Buffer.from(instIDStr, "hex"), keyStr, valueStr);
    } catch (e) {
        Handler.prependLog("failed to add key/value pair: " + e);
    }
}

// ----------------------------------------------------------------------------
// The Handler class is a singleton that offers methods to talk to the cothority
// librairy.

class Handler {
    private static instance: Handler = new Handler();
    private static statusHolder: HTMLElement
    private static loaderHolder: HTMLElement

    static roster: Cothority.network.Roster
    static darc: Cothority.darc.Darc
    static signer: Cothority.darc.SignerEd25519
    static logCounter = 0
    static scid: Buffer // Skip Chain ID

    private constructor() {
        
    }

    static getInstance(): Handler {
        return Handler.instance;
    }

    static prependLog(...nodes: (Node | string)[]) {
        var wrapper = document.createElement("pre")
        var contentWrapper = document.createElement("div")
        var infos = document.createElement("div")
        infos.append(Handler.logCounter + "")
        contentWrapper.append(...nodes)
        wrapper.append(infos, contentWrapper)
        if (Handler.statusHolder == undefined) {
            Handler.statusHolder = document.getElementById("status");
        }
        Handler.statusHolder.prepend(wrapper);
        Handler.logCounter++;
    }

    static checkRoster(): string {
        if (Handler.roster === undefined) {
            return "Roster not set. Please load a roster first"
        }
        return ""
    }

    static checkDarc(): string {
        if (Handler.darc === undefined) {
            return "DARC not set. Please load a DARC first"
        }
        return ""
    }

    static checkSigner(): string {
        if (Handler.signer === undefined) {
            return "Signer not set. Please set a signer first"
        }
        return ""
    }

    static startLoader() {
        if (Handler.loaderHolder === undefined) {
            Handler.loaderHolder = document.getElementById("loader")
        }
        Handler.loaderHolder.classList.add("loading")
    }

    static stopLoader() {
        if (Handler.loaderHolder === undefined) {
            Handler.loaderHolder = document.getElementById("loader")
        }
        Handler.loaderHolder.classList.remove("loading")
    }

    LoadRoster(data: string) {
        Handler.startLoader()
        const roster = Cothority.network.Roster.fromTOML(data)
        const rpc = new Cothority.status.StatusRPC(roster)
        rpc.getStatus(0).then(
            (r) => {
                Handler.roster = roster
                Handler.prependLog("roster loaded!")
            },
            (e) => {
                Handler.prependLog("failed to load roster: " + e)
            }
        ).finally(
            () => Handler.stopLoader()
        )
    }

    SetDarc(scid: Buffer) {
        Handler.startLoader()
        Handler.prependLog("loading the genesis Darc and scid '" + scid.toString("hex") + "'...")
        const rpc = Cothority.byzcoin.ByzCoinRPC.fromByzcoin(Handler.roster, scid)
        rpc.then(
            (r) => {
                Handler.darc = r.getDarc()
                Handler.scid = scid
                Handler.prependLog("darc loaded:\n" + Handler.darc.toString())
            },
            (e) => {
                Handler.prependLog("failed to get the genesis darc: " + e)
            }
        ).finally(
            () => Handler.stopLoader()
        )
    }

    SetSigner(sid: Buffer) {
        Handler.prependLog("setting the signer with: '" + sid.toString("hex") + "'...")
        try {
            var signer = Cothority.darc.SignerEd25519.fromBytes(sid)
            Handler.signer = signer
        } catch(e) {
            Handler.prependLog("failed to create signer: " + e)
        }
        Handler.prependLog("signer '" + signer.toString() + "' set")
    }

    AddRule(ruleStr: string) {
        Handler.startLoader();
        Handler.prependLog("setting the rules " + ruleStr + "...")
        const rpc = Cothority.byzcoin.ByzCoinRPC.fromByzcoin(Handler.roster, Handler.scid)
        rpc.then(
            (r) => {
                var darc = r.getDarc()

                Handler.prependLog("RPC created, getting the darc...")
                Cothority.contracts.darc.DarcInstance.fromByzcoin(r, darc.getBaseID()).then(
                    (darcInstance) => {
                        const evolveDarc = darc.evolve();
                        evolveDarc.addIdentity(ruleStr, Handler.signer, Cothority.darc.Rule.OR)
                        Handler.prependLog("rule '" + ruleStr + "' added on temporary darc...")
                        const evolveInstance = darcInstance.evolveDarcAndWait(evolveDarc, [Handler.signer], 10).then(
                            (evolvedDarcInstance) => {
                                Handler.prependLog("darc instance evolved:\n" + evolvedDarcInstance.darc.toString())
                            },
                            (e) => {
                                Handler.prependLog("failed to evolve the darc instance: " + e)
                            }
                        ).finally(
                            () => Handler.stopLoader()
                        )
                    },
                    (e) => {
                        Handler.stopLoader()
                        Handler.prependLog("failed to get the darc instance")
                    }
                )
            },
            (e) => {
                Handler.stopLoader()
                Handler.prependLog("failed to create RPC: " + e)
            }
        )
    }

    SpawnKV(keyStr: string, valueStr: string) {
        Handler.startLoader()
        Handler.prependLog("creating an RPC to spawn a new key value instance...")
        const rpc = Cothority.byzcoin.ByzCoinRPC.fromByzcoin(Handler.roster, Handler.scid)
        rpc.then(
            (r) => {
                Handler.prependLog("RPC created, we now send a spawn:keyValue request...")
                KVInstance.spawn(r, Handler.darc.getBaseID(), [Handler.signer], keyStr, Buffer.from(valueStr)).then(
                    (kvInstance) => {
                        // Handler.prependLog("Key value instance spawned: " + kvInstance)
                        Handler.prependLog("Key value instance spawned: \n" + kvInstance.toString() + "\nInstance ID: " + kvInstance.id.toString("hex"))
                    },
                    (e) => {
                        console.error(e);
                        Handler.prependLog("failed to spawn the key value instance: " + e)
                    }
                ).finally(
                    () => Handler.stopLoader()
                )
            },
            (e) => {
                Handler.stopLoader()
                Handler.prependLog("failed to create RPC: " + e)
            }
        )
    }

    PrintKV(instIDStr: Buffer) {
        Handler.startLoader()
        Handler.prependLog("creating an RPC to get the key value instance...")
        const rpc = Cothority.byzcoin.ByzCoinRPC.fromByzcoin(Handler.roster, Handler.scid)
        rpc.then(
            (r) => {
                Handler.prependLog("RPC created, we now send a get proof request...")
                r.getProofFromLatest(instIDStr).then(
                    (proof) => {
                        Handler.prependLog("got the proof, let's check it...")
                        if (!proof.exists(instIDStr)) {
                            Handler.prependLog("this is not a proof of existence... aborting!")
                            return
                        }
                        if (!proof.matchContract(KVInstance.contractID)) {
                            Handler.prependLog("this is not a proof for the keyValue contrac... aborting!")
                            return
                        }
                        Handler.prependLog("ok, now let's decode it...")
                        var kvInstance = KeyValueData.decode(proof.value);
                        console.log(kvInstance)
                        Handler.prependLog("here is the key value instance: \n" + kvInstance.toString())
                    },
                    (e) => {
                        console.error(e)
                        Handler.prependLog("failed to get the key value instance: " + e)
                    }
                ).finally(
                    () => Handler.stopLoader()
                )
            },
            (e) => {
                Handler.stopLoader()
                Handler.prependLog("failed to create RPC: " + e)
            }
        )
    }

    AddKVpair(instID: Buffer, keyStr: string, valueStr: string) {
        Handler.startLoader()
        Handler.prependLog("updating the keyValue instance with a new key/value pair...")
        const rpc = Cothority.byzcoin.ByzCoinRPC.fromByzcoin(Handler.roster, Handler.scid)
        rpc.then(
            (r) => {
                Handler.prependLog("RPC created, we now send an invoke:keyValue.update request...")
                var buffValue = EMPTY_BUFFER
                if (valueStr != "") buffValue = Buffer.from(valueStr)
                KVInstance.fromByzcoin(r, instID).then(
                    (instance) => {
                        instance.invokeUpdate([Handler.signer], keyStr, Buffer.from(valueStr)).then(
                            () => {
                                Handler.prependLog("Key value instance updated, let's get the new data")
                                instance.update().then(
                                    (newInstance) => {
                                        Handler.prependLog(newInstance.toString() + "\nInstance ID: " + newInstance.id.toString("hex"))
                                    },
                                    (e) => {
                                        console.error(e);
                                        Handler.prependLog("failed to update the key value instance: " + e)
                                    }
                                ).finally(
                                    () => Handler.stopLoader()
                                )
                            },
                            (e) => {
                                Handler.stopLoader()
                                Handler.prependLog("failed call the invokeUpdate")
                            }
                        )
                    },
                    (e) => {
                        Handler.stopLoader()
                        Handler.prependLog("failed to get the keyValue instance for this instance ID")
                    }
                )
            },
            (e) => {
                Handler.stopLoader()
                Handler.prependLog("failed to create RPC: " + e)
            }
        )
    }

}

// Transforms an hexadecimal string to its byte representation. This is used
// when reading inputs given from the user, which generally come as hex strings.
function hexStringToByte(str: string) {
    if (!str) {
      return new Uint8Array();
    }
    
    var a = [];
    for (var i = 0, len = str.length; i < len; i+=2) {
      a.push(parseInt(str.substr(i,2),16));
    }
    
    return new Uint8Array(a);
}