package main

import (
	_ "github.com/dedis/cothority_template/protocol"
	_ "github.com/dedis/cothority_template/service"
	"github.com/dedis/onet/simul"
)

func main() {
	simul.Start()
}
