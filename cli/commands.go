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
					Name:      "vmcluster",
					ShortName: "vmcl",
					Usage:     "Display one or many virtual machine clusters",
					Action:    command.GetVMCluster,
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "container",
					ShortName: "cnt",
					Usage:     "Display one or many containers",
					Action:    command.GetContainer,
					Flags:     []cli.Flag{flClusterName, flContainerSort, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "containerpod",
					ShortName: "pod",
					Usage:     "Display one or many container pods",
					Action:    command.GetContainerPod,
					Flags:     []cli.Flag{flClusterName, flContainerSort, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "physicalmachine",
					ShortName: "pm",
					Usage:     "Display one or many physical machines that belong to a cluster",
					Action:    command.GetPhysicalMachine,
					Flags:     []cli.Flag{flClusterName, flPhysicalMachineSort, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "pmcluster",
					ShortName: "pmcl",
					Usage:     "Display one or many physical machine clusters",
					Action:    command.GetPMCluster,
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "service",
					ShortName: "s",
					Usage:     "Display one or many services",
					Action:    command.GetService,
					ArgsUsage: "[NAME]",
				},
				{
					Name:      "virtualmachine",
					ShortName: "vm",
					Usage:     "Display one or many virtual machines that belong to a cluster",
					Action:    command.GetVirtualMachine,
					Flags:     []cli.Flag{flClusterName, flVirtualMachineSort, flSupplyChain},
					ArgsUsage: "[NAME]",
				},
			},
		},
	}
)
