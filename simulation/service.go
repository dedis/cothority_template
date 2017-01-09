package main

import (
	"github.com/BurntSushi/toml"
	"github.com/dedis/cothority_template/service"
	"github.com/dedis/onet"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/simul/monitor"
)

/*
 * Defines the simulation for the service-template
 */

func init() {
	onet.SimulationRegister("TemplateService", NewSimulationService)
}

// Simulation only holds the BFTree simulation
type SimulationService struct {
	onet.SimulationBFTree
}

// NewSimulationService returns the new simulation, where all fields are
// initialised using the config-file
func NewSimulationService(config string) (onet.Simulation, error) {
	es := &SimulationService{}
	_, err := toml.Decode(config, es)
	if err != nil {
		return nil, err
	}
	return es, nil
}

// Setup creates the tree used for that simulation
func (e *SimulationService) Setup(dir string, hosts []string) (
	*onet.SimulationConfig, error) {
	sc := &onet.SimulationConfig{}
	e.CreateRoster(sc, hosts, 2000)
	err := e.CreateTree(sc)
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// Node can be used to initialize each node before it will be run
// by the server. Here we call the 'Node'-method of the
// SimulationBFTree structure which will load the roster- and the
// tree-structure to speed up the first round.
func (e *SimulationService) Node(config *onet.SimulationConfig) error {
	index, _ := config.Roster.Search(config.Conode.ServerIdentity.ID)
	log.Lvl1("Initializing node-index", index)
	return e.SimulationBFTree.Node(config)
}

// Run is used on the destination machines and runs a number of
// rounds
func (e *SimulationService) Run(config *onet.SimulationConfig) error {
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
