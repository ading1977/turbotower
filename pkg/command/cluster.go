package command

import (
	"fmt"
	"github.com/turbonomic/turbotower/pkg/influx"
	"github.com/urfave/cli"
)

func GetVMCluster(c *cli.Context) error {
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	row, err := db.Query(influx.NewDBQuery(c).
		WithQueryType("schema").
		WithColumns("VM_CLUSTER").
		WithName("commodity_sold").
		WithConditions("entity_type='VIRTUAL_MACHINE'"))
	if err != nil {
		return err
	}
	for _, value := range row.Values {
		fmt.Println(value[1])
	}
	return nil
}

func GetPMCluster(c *cli.Context) error {
	db, err := influx.NewDBInstance(c)
	if err != nil {
		return err
	}
	defer db.Close()
	row, err := db.Query(influx.NewDBQuery(c).
		WithQueryType("schema").
		WithColumns("HOST_CLUSTER").
		WithName("commodity_sold").
		WithConditions("entity_type='PHYSICAL_MACHINE'"))
	if err != nil {
		return err
	}
	for _, value := range row.Values {
		fmt.Println(value[1])
	}
	return nil
}
