package influx

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbotower/pkg/topology"
	"github.com/turbonomic/turbotower/utils"
	"github.com/urfave/cli"
)

const (
	format_get_app_header  = "%-60s%-10s%-15s%-10s%-10s\n"
	format_get_app_content = "%-60s%-10.2f%-15.2f%-10s%-10s\n"
)

func GetApplication(c *cli.Context) error {
	db, err := newDBInstance(c)
	if err != nil {
		return err
	}
	defer db.close()
	row, err := db.query(newDBQuery(c).
		withColumns("VCPU_USED", "VMEM_USED", "display_name", "oid", "provider_id").
		withName("commodity_bought").
		withConditions("entity_type='APPLICATION'",
			"AND display_name !~/GuestLoad*/",
			"AND time>now()-10m"))
	if err != nil {
		return err
	}
	tp := topology.NewTopology()
	for _, value := range row.Values {
		vcpuUsed, err := value[1].(json.Number).Float64()
		if err != nil {
			return fmt.Errorf("failed to convert VCPU_USED metric %v", value[1])
		}
		vmemUsed, err := value[2].(json.Number).Float64()
		if err != nil {
			return fmt.Errorf("failed to convert VMEM_USED metric %v", value[2])
		}
		entity := topology.
			NewEntity(value[3].(string), value[4].(string)).
			AddCommodityBought(topology.
				NewCommodityBought("VCPU_USED", vcpuUsed), value[5].(string)).
			AddCommodityBought(topology.
				NewCommodityBought("VMEM_USED", vmemUsed), value[5].(string))
		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("Entity created: %+v", entity)
		}
		tp.AddEntity(entity)
	}
	fmt.Printf(format_get_app_header,
		"Name", "VCPU", "VMEM", "QPS", "LATENCY")
	for _, entity := range tp.Entities {
		avgUsed := utils.GetAvgBoughtUsed(entity.CommodityBought)
		avgCPU, _ := avgUsed["VCPU_USED"]
		avgMem, _ := avgUsed["VMEM_USED"]
		fmt.Printf(format_get_app_content,
			utils.Truncate(entity.Name, 55),
			avgCPU, avgMem, "-", "-")
	}
	return nil
}
