Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[Cothority Template](../README.md) ::
Protocol Template

# Protocol Template

In the this directory, you will find the implementation of a toy protocol where
nodes work together to count how many instances of the protocol there are. It
demonstrates how to send messages and how to handle incoming messages.

To implement a new protocol, you must register the kinds of messages it passes,
and also the protocol itself:

```go
func init() {
	network.RegisterMessage(Announce{})
	network.RegisterMessage(Reply{})
	onet.GlobalProtocolRegister(Name, NewProtocol)
}
```

The messages are defined in the file `struct.go`. For each message, you need
to define the message itself, and the message as it will arrive to you from the
cothority server.

After registering, define a struct that implements the
[onet.ProtocolInstance](https://godoc.org/github.com/dedis/onet#ProtocolInstance)
interface:

```go
type TemplateProtocol struct {
  *onet.TreeNodeInstance
  ...
}

// Check that *TemplateProtocol implements onet.ProtocolInstance
var _ onet.ProtocolInstance = (*TemplateProtocol)(nil)
```

Next, define the function that generates a new protocol instance. Using the
newly created `onet.ProtocolInstance`, you can call RegisterHandler in order to
register handlers for each of the message types. Later, the server will choose
the correct handler to call based on which message type arrives.

Any state that is needed by the protocol (for example, the ChildCount channel)
should be initialized here.

Finally, define the `Start` function that will initiate the protocol instance:

```go
func (p *ProtocolExampleHandlers) Start() error {...}
```
