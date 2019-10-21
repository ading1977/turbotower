package command

import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbotower/pkg/topology"
	"github.com/turbonomic/turbotower/utils"
	"strings"
)

func displaySupplyChain(seeds []*topology.Entity, summary bool) {
	nodes := topology.NewSupplyChainResolver().GetSupplyChainNodesFrom(seeds)
	if summary {
		fmt.Printf(format_show_one_header_1,
			"Type", "Count", "Providers", "Consumers")
	}
	for _, node := range nodes {
		if !summary {
			fmt.Printf(format_show_one_header_1,
				"Type", "Count", "Providers", "Consumers")
		}
		entityType, _ := proto.EntityDTO_EntityType_name[node.EntityType]
		count := node.Members.Cardinality()
		providers := strings.Join(node.GetProviderTypes(), ",")
		consumers := strings.Join(node.GetConsumerTypes(), ",")
		fmt.Printf(format_show_one_content_1,
			entityType, count, providers, consumers)
		if !summary {
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
}
