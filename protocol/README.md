Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[Cothority Template](../README.md) ::
Protocol Template

# Protocol Template

In the this directory, you will find the implementation of a toy protocol where
nodes work together to count how many instances of the protocol there are. It
demonstrates how to send messages and how to handle incoming messages.

To implement a new protocol, you must create the `NewProtocol` function and
register it.

```go
func init() {
	_, err := onet.GlobalProtocolRegister(Name, NewProtocol)
	if err != nil {
		panic(err)
	}
}
```

## Writing `NewProtocol` - the protocol constructor

Inside `NewProtocol` you must register the channels that are used for receiving
messages, e.g., `announceChan chan announceWrapper`. Any state that is needed
by the protocol (for example, the ChildCount channel) should be initialized
here too.

The messages are defined in the file `struct.go`. For each message, you need to
define the message itself (e.g., `Announce`), and the message as it will arrive to
you from the cothority server (e.g., `announceWrapper`).


## Writing the protocol logic

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

Usually, an implementation for `Start` and `Dispatch` is needed, the others are
optional and the default implementation will be used if they are not
implemented. `Dispatch` is where the main protocol logic is implemented and it
is called automatically by onet when the protocol is initiated. `Start` is the
entry-point of the protocol, it needs to be called manually, typically by the
root node.
