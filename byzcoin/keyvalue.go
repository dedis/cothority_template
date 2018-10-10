package byzcoin

import (
	"errors"

	"github.com/dedis/cothority/byzcoin"
	"github.com/dedis/cothority/darc"
	"github.com/dedis/protobuf"
)

// The value contract can simply store a value in an instance and serves
// mainly as a template for other contracts. It helps show the possibilities
// of the contracts and how to use them at a very simple example.

// ContractKeyValueID denotes a contract that can store and update
// key/value pairs.
var ContractKeyValueID = "keyValue"

// ContractKeyValue is a simple key/value storage where you
// can put any data inside as wished.
// It can spawn new keyValue instances and will store all the arguments in
// the data field.
// Existing keyValue instances can be "update"d and deleted.
func ContractKeyValue(cdb byzcoin.CollectionView, inst byzcoin.Instruction, cIn []byzcoin.Coin) (scs []byzcoin.StateChange, cOut []byzcoin.Coin, err error) {
	cOut = cIn

	err = inst.VerifyDarcSignature(cdb)
	if err != nil {
		return
	}

	var darcID darc.ID
	_, _, darcID, err = cdb.GetValues(inst.InstanceID.Slice())
	if err != nil {
		return
	}

	switch inst.GetType() {
	case byzcoin.SpawnType:
		// Spawn a new instance of the KeyValue contract.
		// First create a new ContractStruct and encode it as a protobuf.
		cs := NewContractStruct(inst.Spawn.Args)
		var csBuf []byte
		csBuf, err = protobuf.Encode(&cs)
		if err != nil {
			return
		}

		instID := inst.DeriveID("")
		// Then create a StateChange request with the data of the instance. The
		// InstanceID is given by the DeriveID method of the instruction that allows
		// to create multiple instanceIDs out of a given instruction in a pseudo-
		// random way that will be the same for all nodes.
		scs = []byzcoin.StateChange{
			byzcoin.NewStateChange(byzcoin.Create, instID, ContractKeyValueID, csBuf, darcID),
		}
		return

	case byzcoin.InvokeType:
		if inst.Invoke.Command != "update" {
			return nil, nil, errors.New("Value contract can only update")
		}
		// The only command we can invoke is 'update' which will store the new values
		// given in the arguments in the data.
		//  1. decode the existing data
		//  2. update the data
		//  3. encode the data into protobuf again
		var csBuf []byte
		csBuf, _, _, err = cdb.GetValues(inst.InstanceID.Slice())
		cs := KeyValueData{}
		err = protobuf.Decode(csBuf, &cs)
		if err != nil {
			return
		}
		cs.Update(inst.Invoke.Args)
		csBuf, err = protobuf.Encode(&cs)
		if err != nil {
			return
		}
		scs = []byzcoin.StateChange{
			byzcoin.NewStateChange(byzcoin.Update, inst.InstanceID,
				ContractKeyValueID, csBuf, darcID),
		}
		return

	case byzcoin.DeleteType:
		// Delete removes all the data from the global state.
		scs = byzcoin.StateChanges{
			byzcoin.NewStateChange(byzcoin.Remove, inst.InstanceID, ContractKeyValueID, nil, darcID),
		}
		return
	}
	err = errors.New("didn't find any instruction")
	return
}

// NewContractStruct returns an initialised ContractStruct with all key/value
// pairs from the arguments.
func NewContractStruct(args byzcoin.Arguments) KeyValueData {
	cs := KeyValueData{}
	for _, kv := range args {
		cs.Storage = append(cs.Storage, KeyValue{kv.Name, kv.Value})
	}
	return cs
}

// Update goes through all the arguments and:
//  - updates the value if the key already exists
//  - deletes the keyvalue if the value is empty
//  - adds a new keyValue if the key does not exist yet
func (cs *KeyValueData) Update(args byzcoin.Arguments) {
	for _, kv := range args {
		var updated bool
		for i, stored := range cs.Storage {
			if stored.Key == kv.Name {
				updated = true
				if kv.Value == nil || len(kv.Value) == 0 {
					cs.Storage = append(cs.Storage[0:i], cs.Storage[i+1:]...)
					break
				}
				cs.Storage[i].Value = kv.Value
			}

		}
		if !updated {
			cs.Storage = append(cs.Storage, KeyValue{kv.Name, kv.Value})
		}
	}
}
