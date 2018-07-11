package omniledger

import (
	"errors"

	ol "github.com/dedis/cothority/omniledger/service"
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
func ContractKeyValue(cdb ol.CollectionView, inst ol.Instruction, cIn []ol.Coin) (scs []ol.StateChange, cOut []ol.Coin, err error) {
	cOut = cIn
	switch {
	case inst.Spawn != nil:
		// Spawn a new instance of the KeyValue contract.
		// First create a new ContractStruct and encode it as a protobuf.
		cs := NewContractStruct(inst.Spawn.Args)
		var csBuf []byte
		csBuf, err = protobuf.Encode(&cs)
		if err != nil {
			return
		}

		// Then create a StateChange request with the data of the instance. The
		// InstanceID is given by the DeriveID method of the instruction that allows
		// to create multiple instanceIDs out of a given instruction in a pseudo-
		// random way that will be the same for all nodes.
		scs = []ol.StateChange{
			ol.NewStateChange(ol.Create, inst.DeriveID(ContractKeyValueID), ContractKeyValueID, csBuf),
		}
		return

	case inst.Invoke != nil:
		if inst.Invoke.Command != "update" {
			return nil, nil, errors.New("Value contract can only update")
		}
		// The only command we can invoke is 'update' which will store the new values
		// given in the arguments in the data.
		//  1. decode the existing data
		//  2. update the data
		//  3. encode the data into protobuf again
		var csBuf []byte
		csBuf, _, err = cdb.GetValues(inst.InstanceID.Slice())
		cs := ContractStruct{}
		err = protobuf.Decode(csBuf, &cs)
		if err != nil {
			return
		}
		cs.Update(inst.Invoke.Args)
		csBuf, err = protobuf.Encode(&cs)
		if err != nil {
			return
		}
		scs = []ol.StateChange{
			ol.NewStateChange(ol.Update, inst.InstanceID,
				ContractKeyValueID, csBuf),
		}
		return

	case inst.Delete != nil:
		// Delete removes all the data from the global state.
		scs = ol.StateChanges{
			ol.NewStateChange(ol.Remove, inst.InstanceID, ContractKeyValueID, nil),
		}
		return
	}
	err = errors.New("didn't find any instruction")
	return
}

// NewContractStruct returns an initialised ContractStruct with all key/value
// pairs from the arguments.
func NewContractStruct(args ol.Arguments) ContractStruct {
	cs := ContractStruct{}
	for _, kv := range args {
		cs.Storage = append(cs.Storage, KeyValue{kv.Name, kv.Value})
	}
	return cs
}

// Update goes through all the arguments and:
//  - updates the value if the key already exists
//  - deletes the keyvalue if the value is empty
//  - adds a new keyValue if the key does not exist yet
func (cs *ContractStruct) Update(args ol.Arguments) {
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
