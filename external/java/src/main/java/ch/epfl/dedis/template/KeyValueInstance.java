package ch.epfl.dedis.template;

import ch.epfl.dedis.byzcoin.transaction.Argument;
import ch.epfl.dedis.byzcoin.transaction.ClientTransaction;
import ch.epfl.dedis.byzcoin.transaction.Instruction;
import ch.epfl.dedis.byzcoin.transaction.Invoke;
import ch.epfl.dedis.lib.Hex;
import ch.epfl.dedis.lib.exception.CothorityCryptoException;
import ch.epfl.dedis.lib.exception.CothorityException;
import ch.epfl.dedis.lib.exception.CothorityNotFoundException;
import ch.epfl.dedis.byzcoin.*;
import ch.epfl.dedis.byzcoin.contracts.DarcInstance;
import ch.epfl.dedis.lib.darc.Request;
import ch.epfl.dedis.lib.darc.Signature;
import ch.epfl.dedis.lib.darc.Signer;
import ch.epfl.dedis.template.proto.KeyValueProto;
import com.google.protobuf.InvalidProtocolBufferException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * KeyValueInstance represents a key/value store on ByzCoin. It can be initialised either by the
 * instanceID, in which case it will fetch itself the needed data from ByzCoin. Or it is initialized
 * with a proof, then it will simply copy the values stored in the proof to create a new KeyValueInstance.
 */
public class KeyValueInstance {
    private Instance instance;
    private ByzCoinRPC bc;
    private List<KeyValue> keyValues;

    private final static Logger logger = LoggerFactory.getLogger(KeyValueInstance.class);

    /**
     * Instantiates a new KeyValueInstance given a working ByzCoin instance and
     * an instanceId. This instantiator will contact ByzCoin and try to get
     * the current valueInstance. If the instance is not found, or is not of
     * contractId "Value", an exception will be thrown.
     *
     * @param ol is a link to a ByzCoin instance that is running
     * @param id of the value-instance to connect to
     * @throws CothorityException
     */
    public KeyValueInstance(ByzCoinRPC ol, InstanceId id) throws CothorityException {
        this.bc = ol;
        update(id);
    }

    /**
     * Instantiates a KeyValueInstance given a proof.
     *
     * @param bc is a link to a ByzCoin instance that is running
     * @param p is a proof for a valid KeyValue instance
     * @throws CothorityException
     */
    public KeyValueInstance(ByzCoinRPC bc, Proof p) throws CothorityException {
        this.bc = bc;
        update(p);
    }

    /**
     * Spawns a new KeyValueInstance on ByzCoin.
     *
     * @param bc a working ByzCoin instance
     * @param darcInstance a darcInstance with a rule "spawn:keyValue"
     * @param signer with the right to execute the "spawn:keyValue" rule from the darcInstance
     * @param kvs a list of KeyValues to include in the data of the instance
     * @throws CothorityException
     */
    public KeyValueInstance(ByzCoinRPC bc, DarcInstance darcInstance, Signer signer, List<KeyValue> kvs) throws CothorityException{
        this.bc = bc;
        List<Argument> args = new ArrayList<>();
        for (KeyValue kv: kvs){
            args.add(kv.toArgument());
        }
        Proof p = darcInstance.spawnInstanceAndWait("keyValue", signer, args, 10);
        update(p);
    }

    /**
     * Updates an existing KeyValueInstance in case it has been updated in ByzCoin.
     * @throws CothorityException
     */
    public void update() throws CothorityException {
        if (instance == null || bc == null || keyValues == null){
            throw new CothorityException("instance not initialized yet");
        }
        update(instance.getId());
    }

    /**
     * updates the keyvalue instance from a live ByzCoin.
     *
     * @throws CothorityException
     */
    public void update(InstanceId id) throws CothorityException {
        update(bc.getProof(id));
    }

    /**
     * Updates the keyvalue instance from a given proof.
     *
     * @param pr the proof to the KeyValue instance
     * @throws CothorityException
     */
    public void update(Proof pr) throws CothorityException{
        if (!pr.matches()){
            throw new CothorityException("cannot use non-matching proof for update");
        }
        instance = Instance.fromProof(pr);
        if (!instance.getContractId().equals("keyValue")) {
            logger.error("wrong instance: {}", instance.getContractId());
            throw new CothorityNotFoundException("this is not a keyValue instance");
        }
        try {
            KeyValueProto.KeyValueData kvd = KeyValueProto.KeyValueData.parseFrom(instance.getData());
            keyValues = new KeyValueData(kvd).getKeyValueList();
        } catch (InvalidProtocolBufferException e) {
            throw new CothorityException(e);
        }
    }

    /**
     * Creates an instruction to evolve the keyValue in ByzCoin. The signer must have its identity in the current
     * darc as "invoke:update" rule.
     * <p>
     * TODO: allow for evolution if the expression has more than one identity.
     *
     * @param keyValues the keyValues to replace/delete/add to the list.
     * @param owner     must have its identity in the "invoke:update" rule
     * @param pos       position of the instruction in the ClientTransaction
     * @param len       total number of instructions in the ClientTransaction
     * @return Instruction to be sent to ByzCoin
     * @throws CothorityCryptoException
     */
    public Instruction updateKeyValueInstruction(List<KeyValue> keyValues, Signer owner, int pos, int len) throws CothorityCryptoException {
        List<Argument> args = new ArrayList<>();
        for (KeyValue kv : keyValues) {
            args.add(new Argument(kv.getKey(), kv.getValue()));
        }
        Invoke inv = new Invoke("update", args);
        Instruction inst = new Instruction(instance.getId(), Instruction.genNonce(), pos, len, inv);
        try {
            Request r = new Request(instance.getDarcId(), "invoke:update", inst.hash(),
                    Arrays.asList(owner.getIdentity()), null);
            logger.info("Signing: {}", Hex.printHexBinary(r.hash()));
            Signature sign = new Signature(owner.sign(r.hash()), owner.getIdentity());
            inst.setSignatures(Arrays.asList(sign));
        } catch (Signer.SignRequestRejectedException e) {
            throw new CothorityCryptoException(e.getMessage());
        }
        return inst;
    }

    /**
     * Sends a request to update the keyvalue instance but doesn't wait for the request to be delivered.
     *
     * @param keyValues the keyValues to replace/delete/add to the list.
     * @param owner must have its identity in the "invoke:update" rule
     * @return a TransactionId pointing to the transaction that should be included
     * @throws CothorityException
     */
    public void updateKeyValue(List<KeyValue> keyValues, Signer owner) throws CothorityException {
        Instruction inst = updateKeyValueInstruction(keyValues, owner, 0, 1);
        ClientTransaction ct = new ClientTransaction(Arrays.asList(inst));
        bc.sendTransaction(ct);
    }

    /**
     * Asks ByzCoin to update the value and waits until the new value has
     * been stored in the global state.
     *
     * @param keyValues the value to replace the old value.
     * @param owner     is the owner that can sign to evolve the darc
     * @param wait      is the number of blocks to wait for an inclusion
     * @throws CothorityException
     */
    public void updateKeyValueAndWait(List<KeyValue> keyValues, Signer owner, int wait) throws CothorityException {
        Instruction inst = updateKeyValueInstruction(keyValues, owner, 0, 1);
        ClientTransaction ct = new ClientTransaction(Arrays.asList(inst));
        bc.sendTransactionAndWait(ct, wait);
        update();
    }

    /**
     * @return the id of the instance
     */
    public InstanceId getId() {
        return instance.getId();
    }

    /**
     * @return a copy of the key/values stored in this instance.
     */
    public List<KeyValue> getKeyValues() throws CothorityCryptoException {
        List<KeyValue> ret = new ArrayList<>();
        for (KeyValue kv : keyValues) {
            ret.add(kv);
        }
        return ret;
    }

    /**
     * @return the instance used.
     */
    public Instance getInstance() {
        return instance;
    }
}
