package template

import (
	"testing"

	"github.com/dedis/cothority/log"
	"github.com/dedis/cothority/sda"
)

func TestMain(m *testing.M) {
	log.MainTest(m)
}

func TestServiceTemplate(t *testing.T) {
	local := sda.NewTCPTest()
	// generate 5 hosts, they don't connect, they process messages, and they
	// don't register the tree or entitylist
	_, el, _ := local.GenTree(5, true)
	defer local.CloseAll()

	// Send a request to the service
	client := NewClient()
	log.Lvl1("Sending request to service...")
	duration, err := client.Clock(el)
	log.ErrFatal(err, "Couldn't send")
	log.Lvl1("It took", duration, "to go through the tree.")
}
