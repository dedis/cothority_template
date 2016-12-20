/*
* This is a template for creating an app. It only has one command which
* prints out the name of the app.
 */
package main

import (
	"os"

	"github.com/dedis/cothority/libcothority"
	"github.com/dedis/cothority_template/service"
	"github.com/dedis/onet/app/config"
	"github.com/dedis/onet/log"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "Template"
	app.Usage = "Used for building other apps."
	app.Version = "0.1"
	groupsDef := "the group-definition-file"
	app.Commands = []cli.Command{
		{
			Name:      "time",
			Usage:     "measure the time to contact all nodes",
			Aliases:   []string{"t"},
			ArgsUsage: groupsDef,
			Action:    cmdTime,
		},
		{
			Name:      "counter",
			Usage:     "return the counter",
			Aliases:   []string{"t"},
			ArgsUsage: groupsDef,
			Action:    cmdCounter,
		},
	}
	app.Flags = []cli.Flag{
		libcothority.FlagDebug,
	}
	app.Before = func(c *cli.Context) error {
		log.SetDebugVisible(c.Int("debug"))
		return nil
	}
	app.Run(os.Args)
}

// Returns the time needed to contact all nodes.
func cmdTime(c *cli.Context) error {
	log.Info("Time command")
	group := readGroup(c)
	client := template.NewClient()
	resp, err := client.Clock(group.Roster)
	if err != nil {
		log.Fatal("When asking the time:", err)
	}
	log.Infof("Children: %d - Time spent: %f", resp.Children, resp.Time)
	return nil
}

// Returns the number of calls.
func cmdCounter(c *cli.Context) error {
	log.Info("Counter command")
	group := readGroup(c)
	client := template.NewClient()
	counter, err := client.Count(group.Roster.RandomServerIdentity())
	if err != nil {
		log.Fatal("When asking for counter:", err)
	}
	log.Info("Number of requests:", counter)
	return nil
}

func readGroup(c *cli.Context) *config.Group {
	if c.NArg() != 1 {
		log.Fatal("Please give the group-file as argument")
	}
	name := c.Args().First()
	f, err := os.Open(name)
	log.ErrFatal(err, "Couldn't open group definition file")
	group, err := config.ReadGroupDescToml(f)
	log.ErrFatal(err, "Error while reading group definition file", err)
	if len(group.Roster.List) == 0 {
		log.ErrFatalf(err, "Empty entity or invalid group defintion in: %s",
			name)
	}
	return group
}
