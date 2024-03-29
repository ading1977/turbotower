package command

import (
	"fmt"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbotower/pkg/influx"
	"github.com/turbonomic/turbotower/pkg/topology"
	"github.com/urfave/cli"
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
	// Set sort strategy
	sortMetric := c.String("sort")
	sortType := topology.SortTypeCommoditySold
	if sortMetric == "VCPU" || sortMetric == "VMEM" {
		sortMetric += "_USED"
		sortType = topology.SortTypeCommodityBought
	}
	topology.SetEntityListSortStrategy(sortType, sortMetric)
	// Display
	scope := c.String("cluster")
	name := c.Args().Get(0)
	if c.Bool("supplychain") {
		return showApp(scope, name, tp)
	}
	return listApp(scope, name, tp)
}

func listApp(scope, name string, tp *topology.Topology) error {
	if name != "" {
		app := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_APPLICATION))
		if app == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_APPLICATION)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		displayEntities([]*topology.Entity{app}, proto.EntityDTO_APPLICATION)
		return nil
	}
	containerPods := tp.GetContainerPodsInCluster(scope)
	if containerPods == nil {
		return fmt.Errorf("failed to get entities in cluster scope %s", scope)
	}
	nodes := topology.NewSupplyChainResolver().
		WithSearchDirection(topology.Up).
		GetSupplyChainNodesFrom(containerPods)
	for _, node := range nodes {
		if node.EntityType == int32(proto.EntityDTO_APPLICATION) {
			if node.Members.Cardinality() < 1 {
				entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_APPLICATION)]
				return fmt.Errorf("failed to find any entity in the supply chain with type %s", entityType)
			}
			var entities []*topology.Entity
			for entity := range node.Members.Iterator().C {
				entities = append(entities, entity.(*topology.Entity))
			}
			sortedEntities := topology.SortEntities(entities)
			displayEntities(sortedEntities, proto.EntityDTO_APPLICATION)
		}
	}
	return nil
}

func showApp(scope, name string, tp *topology.Topology) error {
	if name != "" {
		app := tp.GetEntityByNameAndType(name, int32(proto.EntityDTO_APPLICATION))
		if app == nil {
			entityType, _ := proto.EntityDTO_EntityType_name[int32(proto.EntityDTO_APPLICATION)]
			return fmt.Errorf("failed to get entity by name %s and type %s", name, entityType)
		}
		displaySupplyChain([]*topology.Entity{app}, false)
		return nil
	}
	containerPods := tp.GetContainerPodsInCluster(scope)
	if containerPods == nil {
		return fmt.Errorf("failed to get entities in cluster scope %s", scope)
	}
	displaySupplyChain(containerPods, true)
	return nil
}
