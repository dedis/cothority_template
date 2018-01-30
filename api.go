package template

/*
The api.go defines the methods that can be called from the outside. Most
of the methods will take a roster so that the service knows which nodes
it should work with.

This part of the service runs on the client or the app.
*/

import (
	"github.com/dedis/cothority"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/network"
)

// ServiceName is used for registration on the onet.
const ServiceName = "Template"

// Client is a structure to communicate with the CoSi
// service
type Client struct {
	*onet.Client
}

// NewClient instantiates a new cosi.Client
func NewClient() *Client {
	return &Client{Client: onet.NewClient(ServiceName, cothority.Suite)}
}

// Clock chooses one server from the Roster at random. It
// sends a ClockRequest to it, which is then processed on the server side
// via the code in the service package.
//
// Clock will return the time in seconds it took to run the protocol.
func (c *Client) Clock(r *onet.Roster) (*ClockResponse, error) {
	dst := r.RandomServerIdentity()
	log.Lvl4("Sending message to", dst)
	reply := &ClockResponse{}
	err := c.SendProtobuf(dst, &ClockRequest{r}, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// Count will return the number of times `Clock` has been called on this
// service-node.
func (c *Client) Count(si *network.ServerIdentity) (int, error) {
	reply := &CountResponse{}
	err := c.SendProtobuf(si, &CountRequest{}, reply)
	if err != nil {
		return -1, err
	}
	return reply.Count, nil
}
