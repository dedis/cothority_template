Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[Cothority Template](../README.md) ::
Simulation

# Simulation

After you know that a new service and protocol work at all (unit testing and integration testing have passed) then you might want to know how that software will behave in larger networks of 10's, 100's or 1000's of nodes. This is the job of the simulation system.

References:
* The README.md: https://github.com/dedis/onet/blob/master/simul/README.md
* The interface you need to implement: https://godoc.org/github.com/dedis/onet#Simulation

Simulations in onet are a powerful way of making sure your code is well behaving
also in bigger settings, including on different servers, and of course to write
simulations used in research papers to make pretty graphs.
In order to write a simulation, you must make a struct that implements the [onet.Simulation](https://godoc.org/github.com/dedis/onet#Simulation) interface.

You will need to implement the `Setup` method to return the
`*onet.SimulationConfig` instance and to create a roster and a tree. The `Setup`
method is run at the beginning of the simulation, on your computer. It prepares
all the necessary structures and can also copy needed files for the actual simulation
run.

The `Node` method is run before the actual simulation is started and is called
once for every node. The simulation framework makes sure that all nodes have
finished their `Node` method before the `Run` is called.

The `Run` method is only called on the root node, which is the first node of
the Roster.

## Running your simulation

The simulation example is in the directory `simulation`. It demonstrates how to
write simulation drivers and configurations for simulating the behaviour of the
protocol alone (files `protocol.go` and `protocol.toml`) or for simulating a client
talking to a service (files `service.go` and `service.toml`). Of course, in our case
we know that this example service always starts one instance of the example
protocol, so the `TemplateService` will be exercising the protocol as well.

You can build the simulator executable with `go build`. If you try to run it
with no options (`./simulator`), it asks for a simulation to run. You must give
it one or more toml files on the commandline.

The protocol.toml file has:

```
Simulation = "TemplateProtocol"
Servers = 8
Bf = 4
Rounds = 10
CloseWait = 6000

Depth
1
2
```

When you run `./simulation protocol.toml` this is what you get:

```
$ ./simulation protocol.toml
1 : (                        simul.startBuild:  54) - Deploying to localhost
1 : (           platform.(*Localhost).Cleanup:  92) - Nothing to clean up
1 : (             platform.(*Localhost).Start: 161) - Starting 8 applications of simulation TemplateProtocol
1 : (                       platform.Simulate: 117) - Started counting children with timeout of 1000
1 : (                       platform.Simulate: 121) - Found all 5 children
1 : (                       platform.Simulate: 139) - Starting new node TemplateProtocol
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 0
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 1
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 2
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 3
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 4
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 5
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 6
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 7
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 8
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 9
1 : (           platform.(*Localhost).Cleanup:  92) - Nothing to clean up
1 : (             platform.(*Localhost).Start: 161) - Starting 8 applications of simulation TemplateProtocol
1 : (                       platform.Simulate: 117) - Started counting children with timeout of 1000
1 : (                       platform.Simulate: 121) - Found all 21 children
1 : (                       platform.Simulate: 139) - Starting new node TemplateProtocol
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 0
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 1
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 2
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 3
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 4
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 5
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 6
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 7
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 8
1 : (          main.(*SimulationProtocol).Run:  83) - Starting round 9
```

This shows that it started 8 conodes (they, however, just running inside of
goroutines; you will not see them as separate processes in a `ps -ef` listing).
This is set from "Servers = 8" in the protocol.toml file.

There are two runs, once for depth = 1 and once for depth = 3. For depth = 1,
and branching factor 4, there are 5 nodes expected (1 root + 4 children = 5).
Once all of the 10 rounds are finished, the conodes are shut down. The
simulation continues for depth = 3, with 1 + 4 + 4*4 = 21 children expected.
Because 21 children > 8 conodes, some of the conodes will be hosting 3 children.

You could add the `-debug 3` argument and get much, much more debug output from
each conode showing what it is seeing and doing.

## Running simulations with "go test"

The `simul_test.go` file shows that simulations can be launched from within
standard Go tests. This would make it possible to have different tests that use
different toml files in order to test different sizes of networks, etc.
