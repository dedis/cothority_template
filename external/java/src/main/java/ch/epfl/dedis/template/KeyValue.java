package ch.epfl.dedis.template;

import ch.epfl.dedis.byzcoin.transaction.Argument;
import ch.epfl.dedis.template.proto.KeyValueProto;

import java.util.Arrays;

/**
 * KeyValue is one element of the KeyValue instance. It holds a key of type string
 * and a value of type byte[].
 */
public class KeyValue {
    private String key;
    private byte[] value;

    /**
     * Create a KeyValue object given its protobuf representation.
     *
     * @param kvp the protobuf representation of the KeyValue
     */
    public KeyValue(KeyValueProto.KeyValue kvp) {
        key = kvp.getKey();
        value = kvp.getValue().toByteArray();
    }

    /**
     * Create a KeyValue object given a key and a value
     *
     * @param key   the key for the object
     * @param value the value for the object
     */
    public KeyValue(String key, byte[] value) {
        this.key = key;
        this.value = value;
    }

    /**
     * @return the key of the object.
     */
    public String getKey() {
        return key;
    }

    /**
     * @return a copy of the value of the object.
     */
    public byte[] getValue() {
        return value.clone();
    }

    /**
     * @param key the new key
     */
    public void setKey(String key) {
        this.key = key;
    }

    /**
     * @param value the new value
     */
    public void setValue(byte[] value) {
        this.value = value;
    }

    /**
     * @return an argument representing the key/value pair.
     */
    public Argument toArgument() {
        return new Argument(getKey(), getValue());
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        KeyValue kv = (KeyValue)o;

        return key.equals(kv.getKey()) && Arrays.equals(value, kv.getValue());
    }
}
