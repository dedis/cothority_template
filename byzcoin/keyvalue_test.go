package byzcoin

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/cothority/v3"
	"go.dedis.ch/cothority/v3/byzcoin"
	"go.dedis.ch/cothority/v3/darc"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/protobuf"
)

func TestKeyValue_Spawn(t *testing.T) {
	// Create a new ledger and prepare for proper closing
	bct := newBCTest(t)
	defer bct.Close()

	// Create a new instance with two key/values:
	//  "one": []byte{1}
	//  "two": []byte{2}
	args := byzcoin.Arguments{
		{
			Name:  "one",
			Value: []byte{1},
		},
		{
			Name:  "two",
			Value: []byte{2},
		},
	}
	// And send it to the ledger.
	instID := bct.createInstance(t, args)

	// Get the proof from byzcoin
	reply, err := bct.cl.GetProof(instID.Slice())
	require.Nil(t, err)
	// Make sure the proof is a matching proof and not a proof of absence.
	pr := reply.Proof
	require.True(t, pr.InclusionProof.Match(instID.Slice()))

	// Get the raw values of the proof.
	_, val, _, _, err := pr.KeyValue()
	require.Nil(t, err)
	// And decode the buffer to a KeyValueData
	cs := KeyValueData{}
	err = protobuf.Decode(val, &cs)
	require.Nil(t, err)
	// Verify all values are in there.
	for i, s := range cs.Storage {
		require.Equal(t, args[i].Name, s.Key)
		require.Equal(t, args[i].Value, s.Value)
	}
}

func TestKeyValue_Invoke(t *testing.T) {
	// Create a new ledger and prepare for proper closing
	bct := newBCTest(t)
	defer bct.Close()

	// Create a new instance with two key/values:
	//  "one": []byte{1}
	//  "two": []byte{2}
	args := byzcoin.Arguments{
		{
			Name:  "one",
			Value: []byte{1},
		},
		{
			Name:  "two",
			Value: []byte{2},
		},
	}
	// And send it to the ledger.
	instID := bct.createInstance(t, args)

	// Get the proof from byzcoin
	reply, err := bct.cl.GetProof(instID.Slice())
	require.Nil(t, err)
	pr1 := reply.Proof

	// Delete the key "one", change "two" and add a "three"
	args = byzcoin.Arguments{
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
	bct.updateInstance(t, instID, args)

	// Store the values of the previous proof in 'values'
	_, v1, _, _, err := pr1.KeyValue()
	require.Nil(t, err)
	var v2 []byte
	prRep2, err := bct.cl.GetProof(instID.Slice())
	require.Nil(t, err)
	_, v2, _, _, err = prRep2.Proof.KeyValue()
	require.NotEqual(t, 0, bytes.Compare(v1, v2), "didn't include new values")

	// Read the content of the instance back into a structure.
	var newArgs KeyValueData
	err = protobuf.Decode(v2, &newArgs)
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

	cs.Update(byzcoin.Arguments{{
		Name:  "one",
		Value: []byte{2},
	}})
	require.Equal(t, 1, len(cs.Storage))
	require.Equal(t, []byte{2}, cs.Storage[0].Value)

	cs.Update(byzcoin.Arguments{{
		Name:  "one",
		Value: nil,
	}})
	require.Equal(t, 0, len(cs.Storage))

	cs.Update(byzcoin.Arguments{{
		Name:  "two",
		Value: []byte{22},
	}})
	require.Equal(t, 1, len(cs.Storage))
	require.Equal(t, []byte{22}, cs.Storage[0].Value)

	cs.Update(byzcoin.Arguments{{
		Name:  "two",
		Value: []byte{},
	}})
	require.Equal(t, 0, len(cs.Storage))
}

// bcTest is used here to provide some simple test structure for different
// tests.
type bcTest struct {
	local   *onet.LocalTest
	signer  darc.Signer
	servers []*onet.Server
	roster  *onet.Roster
	cl      *byzcoin.Client
	gMsg    *byzcoin.CreateGenesisBlock
	gDarc   *darc.Darc
	ct      uint64
}

func newBCTest(t *testing.T) (out *bcTest) {
	out = &bcTest{}
	// First create a local test environment with three nodes.
	out.local = onet.NewTCPTest(cothority.Suite)

	out.signer = darc.NewSignerEd25519(nil, nil)
	out.servers, out.roster, _ = out.local.GenTree(3, true)

	// Then create a new ledger with the genesis darc having the right
	// to create and update keyValue contracts.
	var err error
	out.gMsg, err = byzcoin.DefaultGenesisMsg(byzcoin.CurrentVersion, out.roster,
		[]string{"spawn:keyValue", "invoke:keyValue.update"}, out.signer.Identity())
	require.Nil(t, err)
	out.gDarc = &out.gMsg.GenesisDarc

	// This BlockInterval is good for testing, but in real world applications this
	// should be more like 5 seconds.
	out.gMsg.BlockInterval = time.Second / 2

	out.cl, _, err = byzcoin.NewLedger(out.gMsg, false)
	require.Nil(t, err)
	out.ct = 1

	return out
}

func (bct *bcTest) Close() {
	bct.local.CloseAll()
}

func (bct *bcTest) createInstance(t *testing.T, args byzcoin.Arguments) byzcoin.InstanceID {
	ctx := byzcoin.ClientTransaction{
		Instructions: []byzcoin.Instruction{{
			InstanceID:    byzcoin.NewInstanceID(bct.gDarc.GetBaseID()),
			SignerCounter: []uint64{bct.ct},
			Spawn: &byzcoin.Spawn{
				ContractID: ContractKeyValueID,
				Args:       args,
			},
		}},
	}
	bct.ct++
	// And we need to sign the instruction with the signer that has his
	// public key stored in the darc.
	require.NoError(t, ctx.FillSignersAndSignWith(bct.signer))

	// Sending this transaction to ByzCoin does not directly include it in the
	// global state - first we must wait for the new block to be created.
	var err error
	_, err = bct.cl.AddTransactionAndWait(ctx, 5)
	require.Nil(t, err)
	return ctx.Instructions[0].DeriveID("")
}

func (bct *bcTest) updateInstance(t *testing.T, instID byzcoin.InstanceID, args byzcoin.Arguments) {
	ctx := byzcoin.ClientTransaction{
		Instructions: []byzcoin.Instruction{{
			InstanceID:    instID,
			SignerCounter: []uint64{bct.ct},
			Invoke: &byzcoin.Invoke{
				ContractID: ContractKeyValueID,
				Command:    "update",
				Args:       args,
			},
		}},
	}
	bct.ct++
	// And we need to sign the instruction with the signer that has his
	// public key stored in the darc.
	require.NoError(t, ctx.FillSignersAndSignWith(bct.signer))

	// Sending this transaction to ByzCoin does not directly include it in the
	// global state - first we must wait for the new block to be created.
	var err error
	_, err = bct.cl.AddTransactionAndWait(ctx, 5)
	require.Nil(t, err)
}
