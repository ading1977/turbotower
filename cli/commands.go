package cli

import (
	"github.com/turbonomic/turbotower/pkg/influx"
	"github.com/urfave/cli"
)

var (
	commands = []cli.Command{
		{
			Name: "get",
			ShortName: "g",
			Usage: "Display one or many entities or groups of entities",
			Subcommands: []cli.Command{
				{
					Name: "cluster",
					ShortName: "c",
					Usage: "Display one or many clusters",
					Action: influx.GetCluster,
				},
				{
					Name: "service",
					ShortName: "s",
					Usage: "Display one or many services",
					Action: influx.GetService,
				},
			},
		},
	}
)
