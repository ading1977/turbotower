package influx

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbotower/pkg/topology"
	"github.com/urfave/cli"
	"strconv"
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
	tp := topology.NewTopology()
	err = processCommoditySold(c, db, tp)
	err = processCommodityBought(c, db, tp)
	tp.BuildGraph()
	tp.PrintGraph()
	tp.PrintEntityTypeIndex()
	containerEntities := tp.EntityTypeIndex[int32(proto.EntityDTO_CONTAINER)]
	topology.NewSupplyChainResolver().GetSupplyChainNodes(containerEntities)
	return err
}

func processCommoditySold(c *cli.Context, db *DBInstance, tp *topology.Topology) error {
	commoditySoldTagKeys := []string{"oid", "entity_type", "display_name", "VM_CLUSTER", "HOST_CLUSTER"}
	columns := append(commoditySoldFieldKeys, commoditySoldTagKeys...)
	row, err := db.query(newDBQuery(c).
		withColumns(columns...).
		withName("commodity_sold"))
	if err != nil {
		return err
	}
	index := len(columns) - len(commoditySoldTagKeys) + 1
	for _, value := range row.Values {
		// Parse OID
		oid, err := strconv.ParseInt(value[index].(string), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse OID %v", value[index])
		}
		// Parse entity type
		entityType, found := proto.EntityDTO_EntityType_value[value[index+1].(string)]
		if !found {
			return fmt.Errorf("failed to parse entity type %v", value[index+1])
		}
		// Parse display name
		displayName := value[index+2].(string)
		// Parse group names
		var groupNames []string
		for i := index+3; i < len(columns); i++ {
			valObj := value[i]
			if valObj != nil && valObj.(string) != "" {
				groupNames = append(groupNames, valObj.(string))
			}
		}
		// Get or create the entity
		entity := tp.CreateEntityIfAbsent(displayName, oid, entityType, groupNames...)
		// Parse commodity values
		for i, key := range commoditySoldFieldKeys {
			valObj := value[i+1]
			if valObj == nil {
				if log.GetLevel() >= log.DebugLevel {
					log.Debugf("Field value of %v is nil", key)
				}
				continue
			}
			val, err := value[i+1].(json.Number).Float64()
			if err != nil {
				log.Warningf("Failed to parse %v", value[i+1])
			}
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("Field value of %v is %v", key, val)
			}
			entity.CreateCommoditySoldIfAbsent(key, val)
		}
	}
	return nil
}

func processCommodityBought(c *cli.Context, db *DBInstance, tp *topology.Topology) error {
	commodityBoughtTagKeys := []string{"oid", "provider_id", "entity_type", "display_name", "VM_CLUSTER", "HOST_CLUSTER"}
	columns := append(commodityBoughtFieldKeys, commodityBoughtTagKeys...)
	row, err := db.query(newDBQuery(c).
		withColumns(columns...).
		withName("commodity_bought"))
	if err != nil {
		return err
	}
	index := len(columns) - len(commodityBoughtTagKeys) + 1
	for _, value := range row.Values {
		// Parse OID
		oid, err := strconv.ParseInt(value[index].(string), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse OID %v", value[index])
		}
		// Parse provider ID
		providerId, err := strconv.ParseInt(value[index+1].(string), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse provider ID %v", value[index+1])
		}
		// Parse entity type
		entityType, found := proto.EntityDTO_EntityType_value[value[index+2].(string)]
		if !found {
			return fmt.Errorf("failed to parse entity type %v", value[index+2])
		}
		// Parse display name
		displayName := value[index+3].(string)
		// Parse group names
		var groupNames []string
		for i := index+4; i < len(columns); i++ {
			valObj := value[i]
			if valObj != nil && valObj.(string) != "" {
				groupNames = append(groupNames, valObj.(string))
			}
		}
		// Get or create the entity
		entity := tp.CreateEntityIfAbsent(displayName, oid, entityType, groupNames...)
		// Parse commodity values
		for i, key := range commodityBoughtFieldKeys {
			valObj := value[i+1]
			if valObj == nil {
				if log.GetLevel() >= log.DebugLevel {
					log.Debugf("Field value of %v is nil", key)
				}
				continue
			}
			val, err := value[i+1].(json.Number).Float64()
			if err != nil {
				log.Warningf("Failed to parse %v", value[i+1])
			}
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("Field value of %v is %v", key, val)
			}
			entity.CreateCommodityBoughtIfAbsent(key, val, providerId)
		}
	}
	return nil
}
