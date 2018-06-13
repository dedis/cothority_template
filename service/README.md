Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[Cothority Template](../README.md) ::
Service

# Service Template

The template-service is an example of how to start a protocol from a service.
The source-code is very well documented and explains a lot of details, so
we will cover a more general overview here.

Each conode instantiates one service of each type. A service can load its
configuration and save it for later use. It can receive requests from a client
and can start new protocols. A service starts with the method defined in
`RegisterNewService`:

```go
// service.go
func init() {
	onet.RegisterNewService(Name, newService)
}

type Service struct {
	// We need to embed the ServiceProcessor, so that incoming messages
	// are correctly handled.
	*onet.ServiceProcessor
	path string
	// Count holds the number of calls to 'ClockRequest'
	Count int
}

// ...

func newService(c *onet.Context, path string) onet.Service {
	s := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
		path:             path,
	}
	if err := s.RegisterHandlers(s.ClockRequest, s.CountRequest); err != nil {
		log.ErrFatal(err, "Couldn't register messages")
	}
	return s
}
```

In our example it creates a new `Service`-structure and registers the handlers
that will be called by ONet if a client-request is received. The `Service`-
structure embeds the `onet.ServiceProcessor` which is responsible for the
correct routing of incoming client-requests to the handlers you register.

For convenience, the file `api.go` contains a `Client`-definition that clients
can use to communicate with the service. You have to keep in mind that if a client
instantiates a `NewClient`, this will not have a direct access to the service-
structure. It can only communicate to the service through the `SendProtobuf`
(or simple `Send`) call:

 ```go
 func (c *Client) Count(si *network.ServerIdentity) (int, error) {
 	reply := &CountResponse{}
 	err := c.SendProtobuf(si, &CountRequest{}, reply)
 	if err != nil {
 		return -1, err
 	}
 	return reply.Count, nil
 }
```

This method takes a `ServerIdentity` as a destination for the message. It calls
`SendProtobuf` which is a wrapper around marshalling and unmarshalling of the
messages. The `reply` has to be initialized beforehand, so that `SendProtobuf`
can correctly unmarshal the result received from the service.

Once a client calls `NewClient().Count(si)`, a message will be sent to the
service of conode `si` and the `Service.CountRequest`-method will be called.
In this method the counter is returned through the network, and `SendProtobuf`
unmarshals the response to the `reply` so that the count of calls to
`ClockRequest` can be returned.
