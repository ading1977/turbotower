package influx

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"fmt"
	"github.com/turbonomic/turbotower/pkg/topology"
	"github.com/turbonomic/turbotower/utils"
	"github.com/urfave/cli"
)

const (
	format_get_container_header  = "%-60s%-10s%-10s%-15s%-10s\n"
	format_get_container_content = "%-60s%-10.2f%-10.2f%-15.2f%-10.2f\n"
)

func GetContainer(c *cli.Context) error {
	db, err := newDBInstance(c)
	if err != nil {
		return err
	}
	defer db.close()
	row, err := db.query(newDBQuery(c).
		withColumns("VCPU_USED", "VCPU_CAPACITY",
			"VMEM_USED", "VMEM_CAPACITY", "display_name", "oid").
		withName("commodity_sold").
		withConditions("entity_type='CONTAINER'",
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
		vcpuCapacity, err := value[2].(json.Number).Float64()
		if err != nil {
			return fmt.Errorf("failed to convert VCPU_CAPACITY metric %v", value[2])
		}
		vmemUsed, err := value[3].(json.Number).Float64()
		if err != nil {
			return fmt.Errorf("failed to convert VMEM_USED metric %v", value[3])
		}
		vmemCapacity, err := value[4].(json.Number).Float64()
		if err != nil {
			return fmt.Errorf("failed to convert VMEM_CAPACITY metric %v", value[4])
		}
		entity := topology.
			NewEntity(value[5].(string), value[6].(string)).
			AddCommoditySold(topology.
				NewCommoditySold("VCPU", vcpuUsed, vcpuCapacity)).
			AddCommoditySold(topology.
				NewCommoditySold("VMEM", vmemUsed, vmemCapacity))
		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("Entity created: %+v", entity)
		}
		tp.AddEntity(entity)
	}
	fmt.Printf(format_get_container_header,
		"Name", "VCPU", "VCPU_UTIL", "VMEM", "VMEM_UTIL")
	for _, entity := range tp.Entities {
		usedUtil := utils.GetSoldUsedUtil(entity.CommoditySold)
		cpuUsed, _ := usedUtil["VCPU_USED"]
		cpuUtil, _ := usedUtil["VCPU_UTIL"]
		memUsed, _ := usedUtil["VMEM_USED"]
		memUtil, _ := usedUtil["VMEM_UTIL"]
		fmt.Printf(format_get_container_content,
			utils.Truncate(entity.Name, 55),
			cpuUsed, cpuUtil, memUsed, memUtil)
	}
	return nil
}

