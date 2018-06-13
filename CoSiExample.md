Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[Cothority Template](../README.md) ::
CoSi Explanation

# CoSi Explanation

For an introduction to the CoSi protocol and how it works from a theoretical
point of view, see [CoSi
app](https://github.com/dedis/cothority/tree/master/cosi). This explanation is
about how the ftCoSi protocol is implemented in the cothority-tree. Here we don't
describe the CoSi-app, but only the CoSi-protocol and the CoSi-service.

CoSi has been replaced by ftCoSi, which is the fault-tolerant CoSi
implementation that allows for nodes to fail during the commit phase.
However, the goal of this document is to show how the messages are passed
from one node to another, and for this the CoSi-protocol is simpler to
understand.

## CoSi Protocol

You can find the whole source-code of the CoSi-protocol at
https://github.com/dedis/cothority/tree/master/cosi/protocol.
The first part of the file `cosi.go` registers the protocol with ONet under
the name `CoSi`. In the structure of the protocol we find the four
channels that are used to receive messages from ONet:

```go
	// The channel waiting for Announcement message
	announce chan chanAnnouncement
	// the channel waiting for Commitment message
	commit chan chanCommitment
	// the channel waiting for Challenge message
	challenge chan chanChallenge
	// the channel waiting for Response message
	response chan chanResponse
```

These four channels are initialized during the creation of the protocol:

```go
	if err := node.RegisterChannels(&c.announce, &c.commit, &c.challenge,
		&c.response); err != nil {
		return c, err
	}
```

By giving the address to the channel `&c.announce`, ONet automatically
calls `make` to create the appropriate channel. Once the protocol is ready
to be used, ONet starts a new goroutine with the `Dispatch`-method of the
protocol:

```go
// Dispatch will listen on the four channels we use (i.e. four steps)
func (c *CoSi) Dispatch() error {
	nbrChild := len(c.Children())
	if !c.IsRoot() {
		log.Lvl3(c.Name(), "Waiting for announcement")
		ann := (<-c.announce).Announcement
		err := c.handleAnnouncement(&ann)
		if err != nil {
			return err
		}
	}
	for n := 0; n < nbrChild; n++ {
		log.Lvlf3("%s Waiting for commitment of child %d/%d",
			c.Name(), n+1, nbrChild)
		commit := (<-c.commit).Commitment
		err := c.handleCommitment(&commit)
		if err != nil {
			return err
		}
	}
	if !c.IsRoot() {
		log.Lvl3(c.Name(), "Waiting for Challenge")
		challenge := (<-c.challenge).Challenge
		err := c.handleChallenge(&challenge)
		if err != nil {
			return err
		}
	}
	for n := 0; n < nbrChild; n++ {
		log.Lvlf3("%s Waiting for response of child %d/%d", c.Name(), n+1, nbrChild)
		response := (<-c.response).Response
		err := c.handleResponse(&response)
		if err != nil {
			return err
		}
	}
	<-c.done
	return nil
}
```

The `Dispatch` method waits in term for the appropriate messages in order. This
means that if a message comes in too early, it will have to wait before it is
treated. We enforce an order of the messages. What is missing here is the
possibility to abort the protocol during its runtime. For this we would have to
do something like this for every step:

```go
		ann, ok := (<-c.announce).Announcement
		if !ok{
			// The protocol has been aborted and the c.announce-channel closed
			return error.New("Protocol aborted")
		}
```

The rest of the protocol is how each message is handled once it arrives at the
protocol-level, and a way to add "hooks" (callbacks in JavaScript) to the protocol.

## CoSi service

You can find the CoSi-service source-code here:
https://github.com/dedis/cothority/cosi/service It follows the same structure as
the template-service in [Service](../service/README.md).
