package command

import (
	"github.com/turbonomic/turbotower/pkg/influx"
	"github.com/turbonomic/turbotower/pkg/topology"
	"github.com/urfave/cli"
)

const (
	format_get_app_header  = "%-60s%-10s%-15s%-10s%-10s\n"
	format_get_app_content = "%-60s%-10.2f%-15.2f%-10s%-10s\n"
)

func GetApplication(c *cli.Context) error {
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	tp, err := topology.NewTopologyBuilder(db, c).Build()
	if err != nil {
		return err
	}
	containerPods := tp.GetContainerPodsInCluster(c.String("cluster"))
	topology.NewSupplyChainResolver().GetSupplyChainNodes(containerPods)
	return err
}

