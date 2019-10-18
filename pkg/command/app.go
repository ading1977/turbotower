package command

import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbotower/pkg/influx"
	"github.com/turbonomic/turbotower/pkg/topology"
	"github.com/turbonomic/turbotower/utils"
	"github.com/urfave/cli"
	"strings"
)

var (
	format_list_all_header    = "%-60s%-10s%-15s%-10s%-10s\n"
	format_list_all_content   = "%-60s%-10.2f%-15.2f%-10s%-10s\n"
	format_show_one_header_1  = "%-25s%-25s%-30s%-30s\n"
	format_show_one_content_1 = "%-25s%-25d%-30s%-30s\n"
	format_show_one_header_2  = "%-50s%-30s%-30s\n"
	format_show_one_content_2 = "%-50s%-30s%-30s\n"
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
	appName := c.Args().Get(0)
	if appName == "" {
		listAll(c.String("cluster"), tp)
		return nil
	}
	showOne(appName, tp)
	return nil
}

func listAll(scope string, tp *topology.Topology) {
	containerPods := tp.GetContainerPodsInCluster(scope)
	nodes := topology.NewSupplyChainResolver().GetSupplyChainNodes(containerPods)
	for _, node := range nodes {
		if node.EntityType == int32(proto.EntityDTO_APPLICATION) {
			fmt.Printf(format_list_all_header,
				"Name", "VCPU", "VMEM", "QPS", "LATENCY")
			for entity := range node.Members.Iterator().C {
				app := entity.(*topology.Entity)
				avgValue := utils.GetAvgBoughtValues(app.CommodityBought)
				avgVCPU, _ := avgValue["VCPU_USED"]
				avgVMem, _ := avgValue["VMEM_USED"]
				fmt.Printf(format_list_all_content,
					utils.Truncate(app.Name, 55),
					avgVCPU, avgVMem, "-", "-")
			}
			return
		}
	}
}

func showOne(name string, tp *topology.Topology) {
	app := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_APPLICATION))
	if app == nil {
		return
	}
	nodes := topology.NewSupplyChainResolver().GetSupplyChainNodes([]*topology.Entity{app})
	for _, node := range nodes {
		fmt.Printf(format_show_one_header_1,
			"Type", "Count", "Providers", "Consumers")
		entityType, _ := proto.EntityDTO_EntityType_name[node.EntityType]
		count := node.Members.Cardinality()
		providers := strings.Join(node.GetProviderTypes(), ",")
		consumers := strings.Join(node.GetConsumerTypes(), ",")
		fmt.Printf(format_show_one_content_1,
			entityType, count, providers, consumers)
		fmt.Printf(format_show_one_header_2,
			"Name", "VCPU", "VMEM")
		for item := range node.Members.Iterator().C {
			entity := item.(*topology.Entity)
			name := utils.Truncate(entity.Name, 45)
			VCPU := "-"
			VMem := "-"
			soldValues := utils.GetSoldValues(entity.CommoditySold)
			if used, found := soldValues["VCPU_USED"]; found {
				VCPU = fmt.Sprintf("%.2f", used)
				if capacity, found := soldValues["VCPU_CAPACITY"]; found {
					VCPU += fmt.Sprintf(" (%.2f%%)", used/capacity*100)
				}
			}
			if used, found := soldValues["VMEM_USED"]; found {
				VMem = fmt.Sprintf("%.2f", used)
				if capacity, found := soldValues["VMEM_CAPACITY"]; found {
					VMem += fmt.Sprintf(" (%.2f%%)", used/capacity*100)
				}
			}
			fmt.Printf(format_show_one_content_2,
				name, VCPU, VMem)
		}
		fmt.Println()
	}
}
