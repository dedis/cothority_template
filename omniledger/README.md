Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[../README.md](Cothority Template) ::
OmniLedger Example

# OmniLedger Example

The files in this directory give an overview of how to use OmniLedger with
your own code - writing your own contract.
Before reading on, be sure that you read and understand the documentation found
here: [https://github.com/dedis/cothority/tree/master/omniledger/README.md](OmniLedger Documentation).

When writing a contract, you first have to think what you want to store in
OmniLedger and how you can verify if the data is valid or not. Your contract
will be called by `Instruction`s inside `ClientTransaction`s in two ways:

1. when a client requests a new instance by sending a `Spawn` to a darc that
has the appropriate permissions
2. when a client sends an `Instruction` directly to the already existing
instance

When the `ClientTransaction` arrives at the OmniLedger service, it will be queued,
and then batched together with other transactions to form a new block. This block
will be verified by all nodes, so every node will run _all_ the `ClientTransaction`s
and verify that the output is correct.
One very important caveat is that you should _not_ use any `random` method in
your contract, as this will be different on every node!

## Value storage and update

This is a copy of the `value.go` file in the cothority - until the coins work and
we can make something better... The contract is very simple - you have the following
options:

- `Spawn` takes the `value` argument and stores it in the instance
- `Invoke:update` stores the `value` argument as the new data in the instance

Both of these options are protected by the darc where the value will be stored.
A typical use case is:

1. an admin creates a `darc_user` for the new user with the two rules:
  - `Spawn:exampleValue`

## Running the example

To run the example, you can do:

## Files

The following files are in this directory:

- `service.go` only serves to register the contract with OmniLedger. If you
want to give more power to your service, be sure to look at the
[../service](OmniLedger service example).
- `value.go` defines the contract
