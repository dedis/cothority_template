Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[../README.md](Cothority Template) ::
ByzCoin Example

# ByzCoin Example

The files in this directory give an overview of how to use ByzCoin with
your own code and writing your own contract.
Before reading on, be sure that you read and understand the documentation found
here: [https://github.com/dedis/cothority/tree/master/byzcoin/README.md](ByzCoin Documentation).

When writing a contract, you first have to think what you want to store in
ByzCoin and how you can verify if the data is valid or not. Your contract
will be called by `Instruction`s inside `ClientTransaction`s in two ways:

1. when a client requests a new instance he will send a `Spawn` to a darc that
has the appropriate permissions and this darc will spawn your new contract
2. when a client sends an `Instruction` directly to the already existing
instance

When the `ClientTransaction` arrives at the ByzCoin service, it will be queued,
and then batched together with other transactions to form a new block. This block
will be verified by all nodes, so every node will run _all_ the `ClientTransaction`s
and verify that the output is correct.
One very important caveat is that you should _not_ use any `random` method in
your contract, as this will be different on every node!

## Key Value storage and update

This is a better version of the
[https://github.com/dedis/cothority/tree/master/byzcoin/contracts/value.go](value contract).
It holds a set of key/value pairs and lets you update, add and remove key/value
pairs.

- `Spawn` takes all the arguments and stores them as key/value pairs in the instance
- `Invoke:update` goes over all arguments and either updates the key/value pair
if the key already exists, or adds a new key/value pair, or deletes it, if the
value is empty.

Both of these options are protected by the darc where the value will be stored.
A typical use case is:

1. an admin creates a `darc_user` for the new user with the two rules:
  - `Spawn:keyValue`
  - `Invoke:update`

## Java API

There is a java-api with all the necessary definitions to interact with ByzCoin
and the new contract.
The java files are stored here: [../external/java](java files)

For testing, it uses a docker image. Before running the java
tests, you need to create the docker file:

```
cd (go env GOPATH)/github.com/dedis/cothority_template
make docker
```

Every time you change the go-code, you need to update the docker.

## Files

The following files are in this directory:

- `service.go` only serves to register the contract with ByzCoin. If you
want to give more power to your service, be sure to look at the
[../service](service example).
- `keyvalue.go` defines the contract
- `proto.go` has the definitions that will be translated into protobuf

### A word on ProtoBuf

Usually you start protobuf with a `.proto` file and then translate it to
different languages. Because we're using the `dedis/protobuf` library,
our definitions reside in the `proto.go` files and are translated using
`proto.awk` into `.proto` files. These files are then used to create the
java and javascript definitions.
