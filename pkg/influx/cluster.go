package influx

import (
	"fmt"
	"github.com/urfave/cli"
)

func GetCluster(c *cli.Context) error {
	db, err := newDBInstance(c)
	if err != nil {
		return err
	}
	defer db.close()
	//	results, err := db.query(newDBQuery(c).
	//		withColumns("APPLICATION_USED", "display_name").
	//		withName("commodity_bought").
	//		withConditions("entity_type='VIRTUAL_APPLICATION'", "AND time>now()-10m"))
	row, err := db.query(newDBQuery(c).
		withQueryType("schema").
		withColumns("COMPUTE_CLUSTER").
		withName("commodity_bought").
		withConditions("entity_type='CONTAINER_POD'"))
	if err != nil {
		return err
	}
	for _, value := range row.Values {
		fmt.Println(value[1])
	}
	return nil
}
