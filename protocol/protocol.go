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
	"errors"

	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1/network"
)

func init() {
	network.RegisterMessage(Announce{})
	network.RegisterMessage(Reply{})
	onet.GlobalProtocolRegister(Name, NewProtocol)
}

// TemplateProtocol holds the state of a given protocol.
//
// For this example, it defines a channel that will receive the number
// of children. Only the root-node will write to the channel.
type TemplateProtocol struct {
	*onet.TreeNodeInstance
	ChildCount chan int
}

// Check that *TemplateProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*TemplateProtocol)(nil)

// NewProtocol initialises the structure for use in one round
func NewProtocol(n *onet.TreeNodeInstance) (onet.ProtocolInstance, error) {
	t := &TemplateProtocol{
		TreeNodeInstance: n,
		ChildCount:       make(chan int),
	}
	for _, handler := range []interface{}{t.HandleAnnounce, t.HandleReply} {
		if err := t.RegisterHandler(handler); err != nil {
			return nil, errors.New("couldn't register handler: " + err.Error())
		}
	}
	return t, nil
}

// Start sends the Announce-message to all children
func (p *TemplateProtocol) Start() error {
	log.Lvl3("Starting TemplateProtocol")
	return p.HandleAnnounce(StructAnnounce{p.TreeNode(),
		Announce{"cothority rulez!"}})
}

// HandleAnnounce is the first message and is used to send an ID that
// is stored in all nodes.
func (p *TemplateProtocol) HandleAnnounce(msg StructAnnounce) error {
	log.Lvl3("Parent announces:", msg.Message)
	if !p.IsLeaf() {
		// If we have children, send the same message to all of them
		p.SendToChildren(&msg.Announce)
	} else {
		// If we're the leaf, start to reply
		p.HandleReply(nil)
	}
	return nil
}

// HandleReply is the message going up the tree and holding a counter
// to verify the number of nodes.
func (p *TemplateProtocol) HandleReply(reply []StructReply) error {
	defer p.Done()

	children := 1
	for _, c := range reply {
		children += c.ChildrenCount
	}
	log.Lvl3(p.ServerIdentity().Address, "is done with total of", children)
	if !p.IsRoot() {
		log.Lvl3("Sending to parent")
		return p.SendTo(p.Parent(), &Reply{children})
	}
	log.Lvl3("Root-node is done - nbr of children found:", children)
	p.ChildCount <- children
	return nil
}
