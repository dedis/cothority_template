package main

import (
	// Service needs to be imported here to be instantiated.
	_ "github.com/dedis/cothority_template/service"
	"go.dedis.ch/onet/v3/simul"
)

func main() {
	simul.Start()
}
