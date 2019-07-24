package protocol

/*
The `NewProtocol` method is used to define the protocol and to register
the handlers that will be called if a certain type of message is received.
The handlers will be treated according to their signature.

The protocol-file defines the actions that the protocol needs to do in each
step. The root-node will call the `Start`-method of the protocol. Each
node will only use the `Handle`-methods, and not call `Start` again.
*/

import (
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
)

func init() {
	network.RegisterMessages(Reply{}, Announce{})
	_, err := onet.GlobalProtocolRegister(Name, NewProtocol)
	if err != nil {
		panic(err)
	}
}

// TemplateProtocol holds the state of a given protocol.
//
// For this example, it defines a channel that will receive the number
// of children. Only the root-node will write to the channel.
type TemplateProtocol struct {
	*onet.TreeNodeInstance
	announceChan chan announceWrapper
	repliesChan  chan []replyWrapper
	ChildCount   chan int
}

// Check that *TemplateProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*TemplateProtocol)(nil)

// NewProtocol initialises the structure for use in one round
func NewProtocol(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {
	t := &TemplateProtocol{
		TreeNodeInstance: n,
		ChildCount:       make(chan int),
	}
	if err := n.RegisterChannels(&t.announceChan, &t.repliesChan); err != nil {
		return nil, err
	}
	return t, nil
}

// Start sends the Announce-message to all children
func (p *TemplateProtocol) Start() error {
	log.Lvl3(p.ServerIdentity(), "Starting TemplateProtocol")
	return p.SendTo(p.TreeNode(), &Announce{"cothority rulez!"})
}

// Dispatch implements the main logic of the protocol. The function is only
// called once. The protocol is considered finished when Dispatch returns and
// Done is called.
func (p *TemplateProtocol) Dispatch() error {
	defer p.Done()

	ann := <-p.announceChan
	if p.IsLeaf() {
		return p.SendToParent(&Reply{1})
	}
	p.SendToChildren(&ann.Announce)

	replies := <-p.repliesChan
	n := 1
	for _, c := range replies {
		n += c.ChildrenCount
	}

	if !p.IsRoot() {
		return p.SendToParent(&Reply{n})
	}

	p.ChildCount <- n
	return nil
}
