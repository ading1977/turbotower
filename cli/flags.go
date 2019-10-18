package cli

import "github.com/urfave/cli"

var (
	flClusterName = &cli.StringFlag{
		Name: "cluster, c",
		Required: true,
		Usage: "Specify the `NAME` of the cluster to which the entities belong",
		EnvVar: "TURBO_CLUSTER",
	}
	flAppSort = &cli.StringFlag{
		Name: "sort, s",
		Value: "VCPU",
		Usage: "Specify the `METRIC` to be used to sort the result in a descending order",
		EnvVar: "TURBO_SORT_METRIC",
	}
)
