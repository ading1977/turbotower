package cli

import (
	"github.com/turbonomic/turbotower/pkg/command"
	"github.com/urfave/cli"
)

var (
	commands = []cli.Command{
		{
			Name:      "get",
			ShortName: "g",
			Usage:     "Display one or many entities or groups of entities",
			Subcommands: []cli.Command{
				{
					Name:      "application",
					ShortName: "a",
					Usage:     "Display one or many application",
					Action:    command.GetApplication,
					Flags:     []cli.Flag{flClusterName, flAppSort, flSupplyChain},
					ArgsUsage: "[NAME]",
				},

					{
						Name: "cluster",
						ShortName: "cl",
						Usage: "Display one or many clusters",
						Action: command.GetVMCluster,
						ArgsUsage: "[NAME]",
					},
					{
						Name: "container",
						ShortName: "cnt",
						Usage: "Display one or many containers",
						Action: command.GetContainer,
						Flags:     []cli.Flag{flClusterName, flContainerSort, flSupplyChain},
						ArgsUsage: "[NAME]",
					},
					{
						Name: "containerpod",
						ShortName: "pod",
						Usage: "Display one or many container pods",
						Action: command.GetContainerPod,
						Flags:     []cli.Flag{flClusterName, flContainerSort, flSupplyChain},
						ArgsUsage: "[NAME]",
					},
				/*
					{
						Name: "service",
						ShortName: "s",
						Usage: "Display one or many services",
						Action: command.GetService,
						ArgsUsage: "[NAME]",
					},
					{
						Name: "virtualmachine",
						ShortName: "vm",
						Usage: "Display one or many virtual machines that belong to a cluster",
						Action: command.GetVirtualMachine,
						Flags: []cli.Flag{flClusterName},
						ArgsUsage: "[NAME]",
					},
				*/
			},
		},
	}
)
