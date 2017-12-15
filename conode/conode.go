// Conode is the main binary for running a Cothority server.
// A conode can participate in various distributed protocols using the
// *onet* library as a network and overlay library and the *dedis/crypto*
// library for all cryptographic primitives.
// Basically, you first need to setup a config file for the server by using:
//
//  ./conode setup
//
// Then you can launch the daemon with:
//
//  ./conode
//
package main

import (
	"github.com/dedis/cothority"
	"github.com/dedis/onet/app"

	// Import your service:
	_ "github.com/dedis/cothority_template/service"
	// Here you can import any other needed service for your conode.
	// For example, if your service needs cosi available in the server
	// as well, uncomment this:
	//_ "github.com/dedis/cothority/cosi/service"
)

func main() {
	app.Server(cothority.Suite)
}
