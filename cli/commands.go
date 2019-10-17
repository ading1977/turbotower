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
					Name: "application",
					ShortName: "a",
					Usage: "Display one or many application",
					Action: influx.GetApplication,
					ArgsUsage: "[NAME]",
				},
				/*
				{
					Name: "cluster",
					ShortName: "cl",
					Usage: "Display one or many clusters",
					Action: influx.GetVMCluster,
					ArgsUsage: "[NAME]",
				},
				{
					Name: "container",
					ShortName: "cnt",
					Usage: "Display one or many containers",
					Action: influx.GetContainer,
					ArgsUsage: "[NAME]",
				},
				{
					Name: "containerpod",
					ShortName: "pod",
					Usage: "Display one or many container pods",
					Action: influx.GetContainerPod,
					Flags: []cli.Flag{flClusterName},
					ArgsUsage: "[NAME]",
				},
				{
					Name: "service",
					ShortName: "s",
					Usage: "Display one or many services",
					Action: influx.GetService,
					ArgsUsage: "[NAME]",
				},
				{
					Name: "virtualmachine",
					ShortName: "vm",
					Usage: "Display one or many virtual machines that belong to a cluster",
					Action: influx.GetVirtualMachine,
					Flags: []cli.Flag{flClusterName},
					ArgsUsage: "[NAME]",
				},
 */
			},
		},
	}
)
