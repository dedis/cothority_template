package main

import (
	"github.com/BurntSushi/toml"
	"github.com/dedis/cothority_template/service"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/simul/monitor"
)

/*
 * Defines the simulation for the service-template to be run with
 * `cothority/simul`.
 */

func init() {
	onet.SimulationRegister("TemplateService", NewSimulationService)
}

// Simulation only holds the BFTree simulation
type simulation struct {
	onet.SimulationBFTree
}

// NewSimulationService returns the new simulation, where all fields are
// initialised using the config-file
func NewSimulationService(config string) (onet.Simulation, error) {
	es := &simulation{}
	_, err := toml.Decode(config, es)
	if err != nil {
		return nil, err
	}
	return es, nil
}

// Setup creates the tree used for that simulation
func (e *simulation) Setup(dir string, hosts []string) (
	*onet.SimulationConfig, error) {
	sc := &onet.SimulationConfig{}
	e.CreateRoster(sc, hosts, 2000)
	err := e.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// Run is used on the destination machines and runs a number of
// rounds
func (e *simulation) Run(config *onet.SimulationConfig) error {
	size := config.Tree.Size()
	log.Lvl2("Size is:", size, "rounds:", e.Rounds)
	service, ok := config.GetService(template.Name).(*template.Service)
	if service == nil || !ok {
		log.Fatal("Didn't find service", template.Name)
	}
	for round := 0; round < e.Rounds; round++ {
		log.Lvl1("Starting round", round)
		round := monitor.NewTimeMeasure("round")
		ret, err := service.ClockRequest(&template.ClockRequest{Roster: config.Roster})
		if err != nil {
			log.Error(err)
		}
		resp, ok := ret.(*template.ClockResponse)
		if !ok {
			log.Fatal("Didn't get a ClockResponse")
		}
		if resp.Time <= 0 {
			log.Error("0 time elapsed")
		}
		round.Record()
	}
	return nil
}
