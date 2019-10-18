package cli

import "github.com/urfave/cli"

var (
	flClusterName = &cli.StringFlag{
		Name: "cluster",
		Required: true,
		Usage: "Specify the `NAME` of the cluster to which the entities belong",
		EnvVar: "TURBO_CLUSTER",
	}
)
