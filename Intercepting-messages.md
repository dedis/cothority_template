Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[Cothority Template](../README.md) ::
Intercepting Messages

# Intercepting Messages

When testing protocols it is sometimes useful to intercept messages and decide
if you want to drop a message to test a timeout/exception mechanism. Using the
`LocalTest` structure, you can set up something like the following in your test:

```go
local := onet.NewLocalTest()
nbrNodes := 2
servers, _, tree := local.GenTree(nbrNodes, true)
for _, s := range servers {
	serv := s
	serv.RegisterProcessorFunc(onet.ProtocolMsgID, func(e *network.Envelope) {
		// protoMsg holds also To and From fields that can help decide
		// whether a message should be sent or not.
		protoMsg := e.Msg.(*onet.ProtocolMsg)
		_, protoMsg, err := network.Unmarshal(protoMsg.MsgSlice)
		if err != nil {
			log.Error(err)
		} else {
			log.Lvlf1("Got message %#v", protoMsg)
		}
		// Finally give the message back to onet. If this last call is not
		// made, the message is dropped.
		local.Overlays[serv.ServerIdentity.ID].Process(e)
	})
}
```
