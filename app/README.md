Navigation: [DEDIS](https://github.com/dedis/doc/tree/master/README.md) ::
[Cothority Template](../README.md) ::
App

# App Template

The template app is a good starting point to learn how to add a service into a
conode, and then talk to it.

We will look at the app, which is a client that talks to the server (one
conode). Then we'll see how to start the server and interact with it. The
protocol running inside the server is discussed in
[Protocol](../protocol/README.md). Finally we will see how to use the `test.sh`
script to do integration testing.

## The app

First look at the `app.go` file in this directory. It defines a command-line
interface that accepts a debug level. In this example, there are two commands:

* `time` - measures the time needed to contact all nodes in a tree
* `count` - returns how many requests have been made so far

Both commands take as an argument the group-definition file that defines which nodes
to contact.

We use [the urfave/cli package](https://godoc.org/gopkg.in/urfave/cli.v1) to
define the different commands.  In this template, a standard-flag of `-d` is
added that takes a debug-level  to change the output of our app during
debugging. If you need more commands,  you can simply add more of the
`Name`/`Action`. More options are available from the CLI package, see the GoDoc
documentation on it, linked above.

```go
func main() {
	app := cli.NewApp()
	...
	app.Commands = []cli.Command{
		{
			Name:      "time",
			Action:    cmdTime,
		},
		{
			Name:      "counter",
			Action:    cmdCounter,
		},
	}
	...
	app.Run(os.Args)
}
```

#### The command

The `cmdTime` and `cmdCounter` are very similar, so let's look at one, the `cmdTime`:

```
// Returns the time needed to contact all nodes.
func cmdTime(c *cli.Context) error {
	log.Info("Time command")
	roster := readGroup(c)
	client := template.NewClient()
	resp, err := client.Clock(roster)
	if err != nil {
		log.Fatal("When asking the time:", err)
	}
	log.Infof("Children: %d - Time spent: %f", resp.Children, resp.Time)
	return nil
}
```

Here, we make a new client and use it to execute the `Clock` method in the server.

A cothority server has one or more services loaded into it at compile time via
package import. In this case, we'll be compiling the server to include our one
service, with the `Clock` and `Count` methods in it.

For more information about how a service is written, see
[Service](../service/README.md).

## Running it

If you just try to build the app with `go build` and then run it with `./app
time`, you will find that you still need to give the group-file as an argument.

The easiest way to run the status app is to use the `test.sh` script. This script
takes care of incorporating our service and protocol into a `cothority`-skeleton
so that it will be run and can be contacted from the outside.

To see how this works, run `./test.sh -nt` (-nt: "no temp directory").
This will create a directory called `build`.

In the build directory, you will see that the `test.sh`-script created a
`conode`-binary that can be run to offer the service we just wrote. For
simulations, three conodes are defined in `co1`, `co2`, `co3` and their
respective definition in `coX/public.toml`. The `test.sh` script automatically
generates `public.toml` containing the first two conodes. It looks something
like:

```
[[servers]]
  Address = "tcp://127.0.0.1:2002"
  Public = "wxbIMiZ6eOdpYjL8K5xAwVQlCXGPvVYAsc5v8sWlxtI="
  Description = "New Cothority 1"
[[servers]]
  Address = "tcp://127.0.0.1:2004"
  Public = "dUKrAmIqz8WcDW5MRLf4+iVcKPl45hq1MdjBQiV/mok="
  Description = "New Cothority 2"
```

The public keys will differ in your case. If you want to run the cothorities
manually, you will have to give the correct configuration, like:

```
./conode -c co1/private.toml -d 3
```

Here `-d 3` specifies additional debugging-output. In a second window, run:

```
./conode -c co2/private.toml -d 3
```

To run a request, in a third window, type:

```
./app -d 3 time public.toml
```

The `app` client executable starts, chooses one of the servers from `public.toml`,
contacts it and requests it to run the `Clock` method. In response, the service in
the conode generates a n-ary tree with 2 children per node covering the entire
roster. It starts the protocol running on that tree. The protocol in our case
floods a message down to each child, and every child replies with either 1 (leaf
node) or the count of the replies it gets. Eventually the answer arrives at the
originating conode, and it is returned to the client, along with the amount of
time it took.

Voila!

## The test

As we do pre-test-driven development in DEDIS, we propose you start with writing
the integration-test in a bash-file that reflects what you really want to do
with the app.

The `test.sh` script will build your app and conode (with your services and
protocols included) every time you run it. If you want it to use a cached copy
from the previous build, give it the `-nt` argument. If you want to force a
build, give it the `-b` argument.

Looking at `test.sh`, you see the `main`-function:

```
main(){
    startTest
    buildConode
    test Count
    test Time
    stopTest
}
```

The most important lines here are the `test .*`-lines that indicate what tests
to run. `test Time` will run the cleanup function from libtest.sh, and then call
the local function called `testTime`.

Here is a good place to think about your app and reflect how a user should
interact with it. If it is difficult to write this integration-test for your
app, probably a user will also have difficulties with your app. So keep it
simple and clear!

Let's look at the `testTime`-function:

```
testTime(){
    testFail runTmpl time
    testOK runTmpl time public.toml
    testGrep Time runTmpl time public.toml
}
```

First we check that the app fails when asked to report the `time` without a
`public.toml` file. `testFail` runs the command and checks that the exit-code is
different from `0`. In bash, an exit-code of `0` is `success`, which is
different behavior than in C, where `0` often means `false`.

Then we check that the app runs when given a `public.toml` file using `testOK`.

Finally we check that the binary returns a `Time`-string when asked about the
time.

# Acknowledgement

This text is based on a first version written by Matthieu Girod - thank you very
much!
