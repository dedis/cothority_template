//go:build vartime
// +build vartime

package main

import (
	"go.dedis.ch/cothority/v3"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/kyber/v3/pairing/bn256"
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
