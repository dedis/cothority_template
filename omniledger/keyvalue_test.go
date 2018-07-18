package omniledger

import (
	"bytes"
	"testing"
	"time"

	"github.com/dedis/cothority"
	"github.com/dedis/cothority/omniledger/darc"
	"github.com/dedis/cothority/omniledger/service"
	"github.com/dedis/onet"
	"github.com/dedis/protobuf"
	"github.com/stretchr/testify/require"
)

func TestKeyValue_Spawn(t *testing.T) {
	// Create a new omniledger and prepare for proper closing
	olt := newOLTest(t)
	defer olt.Close()

	// Create a new instance with two key/values:
	//  "one": []byte{1}
	//  "two": []byte{2}
	args := service.Arguments{
		{
			Name:  "one",
			Value: []byte{1},
		},
		{
			Name:  "two",
			Value: []byte{2},
		},
	}
	// And send it to OmniLedger.
	instID := olt.createInstance(t, args)

	// Wait for the proof to be available.
	pr, err := olt.cl.WaitProof(instID, olt.gMsg.BlockInterval, nil)
	require.Nil(t, err)
	// Make sure the proof is a matching proof and not a proof of absence.
	require.True(t, pr.InclusionProof.Match())

	// Get the raw values of the proof.
	values, err := pr.InclusionProof.RawValues()
	require.Nil(t, err)
	// And decode the buffer to a ContractStruct.
	cs := KeyValueData{}
	err = protobuf.Decode(values[0], &cs)
	require.Nil(t, err)
	// Verify all values are in there.
	for i, s := range cs.Storage {
		require.Equal(t, args[i].Name, s.Key)
		require.Equal(t, args[i].Value, s.Value)
	}
}

func TestKeyValue_Invoke(t *testing.T) {
	// Create a new omniledger and prepare for proper closing
	olt := newOLTest(t)
	defer olt.Close()

	// Create a new instance with two key/values:
	//  "one": []byte{1}
	//  "two": []byte{2}
	args := service.Arguments{
		{
			Name:  "one",
			Value: []byte{1},
		},
		{
			Name:  "two",
			Value: []byte{2},
		},
	}
	// And send it to OmniLedger.
	instID := olt.createInstance(t, args)

	// Wait for the proof to be available.
	pr1, err := olt.cl.WaitProof(instID, olt.gMsg.BlockInterval, nil)
	require.Nil(t, err)

	// Delete the key "one", change "two" and add a "three"
	args = service.Arguments{
		{
			Name:  "one",
			Value: nil,
		},
		{
			Name:  "two",
			Value: []byte{22},
		},
		{
			Name:  "three",
			Value: []byte{3},
		},
	}
	olt.updateInstance(t, instID, args)

	// Wait for the new values to be written.
	// Store the values of the previous proof in 'values'
	_, values1, err := pr1.KeyValue()
	require.Nil(t, err)
	var values2 [][]byte
	// Try 10 times to get other values than that from OmniLedger.
	var i int
	for i = 0; i < 10; i++ {
		prRep2, err := olt.cl.GetProof(instID.Slice())
		require.Nil(t, err)
		_, values2, err = prRep2.Proof.KeyValue()
		if bytes.Compare(values1[0], values2[0]) != 0 {
			break
		}
		time.Sleep(olt.gMsg.BlockInterval)
	}
	require.NotEqual(t, 10, i, "didn't include new values in time")

	// Read the content of the instance back into a structure.
	var newArgs KeyValueData
	err = protobuf.Decode(values2[0], &newArgs)
	require.Nil(t, err)
	// Verify the content is as it is supposed to be.
	require.Equal(t, 2, len(newArgs.Storage))
	require.Equal(t, "two", newArgs.Storage[0].Key)
	require.Equal(t, []byte{22}, newArgs.Storage[0].Value)
	require.Equal(t, "three", newArgs.Storage[1].Key)
	require.Equal(t, []byte{3}, newArgs.Storage[1].Value)
}

func TestContractStruct_Update(t *testing.T) {
	cs := KeyValueData{
		Storage: []KeyValue{{
			Key:   "one",
			Value: []byte{1},
		}},
	}

	cs.Update(service.Arguments{{
		Name:  "one",
		Value: []byte{2},
	}})
	require.Equal(t, 1, len(cs.Storage))
	require.Equal(t, []byte{2}, cs.Storage[0].Value)

	cs.Update(service.Arguments{{
		Name:  "one",
		Value: nil,
	}})
	require.Equal(t, 0, len(cs.Storage))

	cs.Update(service.Arguments{{
		Name:  "two",
		Value: []byte{22},
	}})
	require.Equal(t, 1, len(cs.Storage))
	require.Equal(t, []byte{22}, cs.Storage[0].Value)

	cs.Update(service.Arguments{{
		Name:  "two",
		Value: []byte{},
	}})
	require.Equal(t, 0, len(cs.Storage))
}

// olTest is used here to provide some simple test structure for different
// omniledger tests.
type olTest struct {
	local   *onet.LocalTest
	signer  darc.Signer
	servers []*onet.Server
	roster  *onet.Roster
	cl      *service.Client
	gMsg    *service.CreateGenesisBlock
	gDarc   *darc.Darc
}

func newOLTest(t *testing.T) (olt *olTest) {
	olt = &olTest{}
	// First create a local test environment with three nodes.
	olt.local = onet.NewTCPTest(cothority.Suite)

	olt.signer = darc.NewSignerEd25519(nil, nil)
	olt.servers, olt.roster, _ = olt.local.GenTree(3, true)
	olt.cl = service.NewClient()

	// Then create a new omniledger with the genesis darc having the right
	// to create and update keyValue contracts.
	var err error
	olt.gMsg, err = service.DefaultGenesisMsg(service.CurrentVersion, olt.roster,
		[]string{"spawn:keyValue", "spawn:darc", "invoke:update"}, olt.signer.Identity())
	require.Nil(t, err)
	olt.gDarc = &olt.gMsg.GenesisDarc

	// This BlockInterval is good for testing, but in real world applications this
	// should be more like 5 seconds.
	olt.gMsg.BlockInterval = time.Second / 2

	_, err = olt.cl.CreateGenesisBlock(olt.gMsg)
	require.Nil(t, err)
	return olt
}

func (olt *olTest) Close() {
	olt.local.CloseAll()
}

func (olt *olTest) createInstance(t *testing.T, args service.Arguments) service.InstanceID {
	ctx := service.ClientTransaction{
		Instructions: []service.Instruction{{
			InstanceID: service.InstanceID{
				DarcID: olt.gDarc.GetBaseID(),
				SubID:  service.SubID{},
			},
			Nonce:  service.Nonce{},
			Index:  0,
			Length: 1,
			Spawn: &service.Spawn{
				ContractID: ContractKeyValueID,
				Args:       args,
			},
		}},
	}
	// And we need to sign the instruction with the signer that has his
	// public key stored in the darc.
	require.Nil(t, ctx.Instructions[0].SignBy(olt.signer))

	// Sending this transaction to OmniLedger does not directly include it in the
	// global state - first we must wait for the new block to be created.
	var err error
	_, err = olt.cl.AddTransaction(ctx)
	require.Nil(t, err)
	return ctx.Instructions[0].DeriveID(ContractKeyValueID)
}

func (olt *olTest) updateInstance(t *testing.T, instID service.InstanceID, args service.Arguments) {
	ctx := service.ClientTransaction{
		Instructions: []service.Instruction{{
			InstanceID: instID,
			Nonce:      service.Nonce{},
			Index:      0,
			Length:     1,
			Invoke: &service.Invoke{
				Command: "update",
				Args:    args,
			},
		}},
	}
	// And we need to sign the instruction with the signer that has his
	// public key stored in the darc.
	require.Nil(t, ctx.Instructions[0].SignBy(olt.signer))

	// Sending this transaction to OmniLedger does not directly include it in the
	// global state - first we must wait for the new block to be created.
	var err error
	_, err = olt.cl.AddTransaction(ctx)
	require.Nil(t, err)
}
