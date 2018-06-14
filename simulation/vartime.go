// +build vartime

package main

import (
	"github.com/dedis/cothority"
	"github.com/dedis/kyber"
	"github.com/dedis/kyber/pairing"
	"github.com/dedis/kyber/pairing/bn256"
)

func init() {
	cothority.Suite = struct {
		pairing.Suite
		kyber.Group
	}{
		Suite: bn256.NewSuite(),
		Group: bn256.NewSuiteG2(),
	}
}
