[![Build Status](https://travis-ci.org/dedis/cothority_template.svg?branch=master)](https://travis-ci.org/dedis/cothority_template)

Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
Cothority Template

# Template for a new cothority protocol/service/app

The following paragraphs give a short introduction in the onet framework and
should be your starting point if this is a new project for you.

## Terminology

Onet has many concepts and components which are described below:

- Server members of the Onet are called Nodes or Conodes. They are started by the
conode app from the cothority-repository.
- Services are run by nodes. They keep a persistent state and need to
be restarted if the node crashes.
	- Can start protocols
	- Communicate with other service-instances
	- Communicate with the clients through an API
	- Have state that is kept over restart of the server
- Protocols are an exchange of messages started by a service and joined by the
others.
	- Have a well-defined number of steps
	- Have an entry and an exit point
	- Are usable for different purposes
- Apps are the end-user programs that communicate with the
services to initiate actions (such as requesting node status for the status app).
- Cothority is a collective authority formed by any group of two or more conodes.
- Conode - or a Cothority-node: `conode` running on a server.
- Roster - the list of conodes present in a cothority
- ServerIdentity - the information needed to identity a conode

## Directory Overview

Building on the ONet-library available at
https://github.com/dedis/onet, this
repo holds templates to build the different parts necessary for a cothority
addition:

- [protocol](protocol) - define an ephemeral, distributed, decentralized protocol
- [service](service) - create a long-term service that can spawn any number of protocols
- [app](app) - write an app that will interact with, or spawn, a cothority
- [simulation](simulation) - how to create a simulation of a protocol or service
- [byzcoin](byzcoin) - how to write a contract using ByzCoin (early alpha!)

This repo is geared towards PhD-students who want to add a new functionality to
the cothority by creating their own protocols, services, simulations or apps.

## Testing and Simulating

You can test your code at 3 different levels, from smallest to biggest:

- go-test - using `LocalTest`, protocols and services can be tested using to golang-framework
- integration testing - a small bash-testing framework is available to write full integration tests for the applications to make sure that everything will work for the users
- [simulation](TemplateSimulation.md) - when running tests with a bigger number of nodes (more than one hundred), a simulation can launch the required nodes on a simulation platform like Deterlab, servers with Mininet or any cloud-platform with an SSH-access

The Go tests should be written for protocols, services and apps alike, while the simulation is only necessary if you want to measure the performance of your protocol. Integration testing is only used for apps.

## Setting up your own repository

Just for testing you can `go get github.com/dedis/cothority_template`. For setting
up a new protocol/service/simulation, we propose that you create a new personal
repository in your account and then copy over the necessary files. Then you
will need to replace all the `github.com/dedis/cothority_tempate` references
with the path of your repository, e.g. `github.com/foo/super_protocol` if your
account is `foo` and the repository is `super_protocol`.
If you happen to do a semester project for DEDIS, please ask your responsible to
set up a `github.com/dedis/student_yy_name` repository for you.

The Perl pie to the rescue (or `sed -i` if you prefer...):

```bash
find . -name "*go" | xargs \
perl -pi -e "s:github.com/dedis/cothority_template:github.com/foo/super_protocol:"
```

**Note:** Everywhere you see the word "template" in the code, you should imagine
that you'll be substituting in your own application name when you are ready to
fork this repository and start your project!

## Documentation

Privacy preserving, decentralized, distributed, blockchain-related, and lots of
other buzzwords are covered with our Cothority-framework. Different projects are
done using our framework in EPFL and other Universities. Here is some overview
of what you can do and what not.

- Template descriptions and overviews of the different parts of this repository.
  - [Protocol](protocol/README.md) - what is in a protocol
  - [Service](service/README.md) - the basics of a service
	- [Simulation](simulation/README.md) - how to run the protocol on different platforms
  - [App](app/README.md) - how to create an app for your service
- [CoSiExample](CoSiExample.md) shows you how the ideas of a paper have been implemented
- [Coding](https://github.com/dedis/Coding) technical aspects of programming in Cothority

Some more specific subjects that might help you:

- [Intercepting messages](Intercepting-messages.md) how to test a protocol by intercepting and
eventually dropping messages - very useful in tests and simulations.

## To cothority and beyond

More documentation and examples can be found in the cothority-repository:
- To run and use a conode, have a look at
	[Cothority Node](https://github.com/dedis/cothority)
	with examples of protocols, services and apps
- To participate as a core-developer, go to
	[Cothority Network Library](https://github.com/dedis/onet)

## License

All repositories for the cothority are double-licensed under a
GNU/AGPL 3.0 and a commercial license. If you want to have more information,
contact us at dedis@epfl.ch.

## Contribution

If you want to contribute to Cothority-ONet, please have a look at
[CONTRIBUTION](https://github.com/dedis/cothority/blob/master/CONTRIBUTION) for
licensing details. Once you are OK with those, you can have a look at our
coding-guidelines in
[Coding](https://github.com/dedis/Coding). In short, we use the github-issues
to communicate and pull-requests to do code-review. Travis makes sure that
everything goes smoothly. And we'd like to have good code-coverage.

# Contact

You can contact us at https://groups.google.com/forum/#!forum/cothority
