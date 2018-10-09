package ch.epfl.dedis.template;

import ch.epfl.dedis.byzcoin.Instance;
import ch.epfl.dedis.template.proto.KeyValueProto;
import com.google.protobuf.InvalidProtocolBufferException;

import java.util.ArrayList;
import java.util.List;

/**
 * KeyValueData represents the data stored in a KeyValue instance. A KeyValue instance
 * stores a list of key/value pairs and lets you add, update, or delete them.
 */
public class KeyValueData {
    private List<KeyValue> keyValueList;

    /**
     * Create a KeyValueData object given its protobuf representation.
     *
     * @param csProto the protobuf representation.
     */
    public KeyValueData(KeyValueProto.KeyValueData csProto) {
        keyValueList = new ArrayList<>();
        for (KeyValueProto.KeyValue kvp : csProto.getStorageList()) {
            keyValueList.add(new KeyValue(kvp));
        }
    }

    /**
     * Create a KeyValueData object given its binary representation of the protobuf.
     *
     * @param data binary representation of the protobuf
     * @throws InvalidProtocolBufferException
     */
    public KeyValueData(byte[] data) throws InvalidProtocolBufferException {
        this(KeyValueProto.KeyValueData.parseFrom(data));
    }

    /**
     * Create a KeyValueData object given an instance.
     *
     * @param inst the instance that holds the KeyValueData
     * @throws InvalidProtocolBufferException
     */
    public KeyValueData(Instance inst) throws InvalidProtocolBufferException {
        this(inst.getData());
    }

    /**
     * Returns a copy of the KeyValue list.
     *
     * @return a copy of the KeyValue list.
     */
    public List<KeyValue> getKeyValueList() {
        List<KeyValue> kvCopy = new ArrayList<>();
        for (KeyValue kv : keyValueList) {
            kvCopy.add(kv);
        }
        return kvCopy;
    }
}
