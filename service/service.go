package template

/*
The service.go defines what to do for each API-call. This part of the service
runs on the node.
*/

import (
	"time"

	"github.com/dedis/cothority_template/protocol"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/log"
	"gopkg.in/dedis/onet.v1/network"
)

// Name is the name to refer to the Template service from another
// package.
const Name = "Template"

func init() {
	onet.RegisterNewService(Name, newService)
}

// Service is our template-service
type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*onet.ServiceProcessor
	path string
	// Count holds the number of calls to 'ClockRequest'
	Count int
}

// ClockRequest starts a template-protocol and returns the run-time.
func (s *Service) ClockRequest(req *ClockRequest) (network.Message, onet.ClientError) {
	s.Count++
	tree := req.Roster.GenerateBinaryTree()
	pi, err := s.CreateProtocol(template.Name, tree)
	if err != nil {
		return nil, onet.NewClientError(err)
	}
	start := time.Now()
	pi.Start()
	<-pi.(*template.ProtocolTemplate).ChildCount
	elapsed := time.Now().Sub(start).Seconds()
	return &ClockResponse{elapsed}, nil
}

// CountRequest returns the number of instantiations of the protocol.
func (s *Service) CountRequest(req *CountRequest) (network.Message, onet.ClientError) {
	return &CountResponse{s.Count}, nil
}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
// If you use CreateProtocolSDA, this will not be called, as the SDA will
// instantiate the protocol on its own. If you need more control at the
// instantiation of the protocol, use CreateProtocolService, and you can
// give some extra-configuration to your protocol in here.
func (s *Service) NewProtocol(tn *onet.TreeNodeInstance, conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	log.Lvl3("Not templated yet")
	return nil, nil
}

// newTemplate receives the context and a path where it can write its
// configuration, if desired. As we don't know when the service will exit,
// we need to save the configuration on our own from time to time.
func newService(c *onet.Context) onet.Service {
	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	if err := s.RegisterHandlers(s.ClockRequest, s.CountRequest); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	return s
}
