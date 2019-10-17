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
	tp.PrintGraph()
	return err
}

func processCommoditySold(c *cli.Context, db *DBInstance, tp *topology.Topology) error {
	columns := append(commoditySoldFieldKeys,
		"oid", "entity_type", "display_name")
	row, err := db.query(newDBQuery(c).
		withColumns(columns...).
		withName("commodity_sold"))
	if err != nil {
		return err
	}
	index := len(columns)
	for _, value := range row.Values {
		// Parse OID
		oid, err := strconv.ParseInt(value[index-2].(string), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse OID %v", value[index-2])
		}
		// Parse entity type
		entityType, found := proto.EntityDTO_EntityType_value[value[index-1].(string)]
		if !found {
			return fmt.Errorf("failed to parse entity type %v", value[index-1])
		}
		// Parse display name
		displayName := value[index].(string)
		// Get or create the entity
		entity := tp.CreateEntityIfAbsent(displayName, oid, entityType)
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
	columns := append(commodityBoughtFieldKeys,
		"oid", "provider_id", "entity_type", "display_name")
	row, err := db.query(newDBQuery(c).
		withColumns(columns...).
		withName("commodity_bought"))
	if err != nil {
		return err
	}
	index := len(columns)
	for _, value := range row.Values {
		// Parse OID
		oid, err := strconv.ParseInt(value[index-3].(string), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse OID %v", value[index-3])
		}
		// Parse provider ID
		providerId, err := strconv.ParseInt(value[index-2].(string), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse provider ID %v", value[index-2])
		}
		// Parse entity type
		entityType, found := proto.EntityDTO_EntityType_value[value[index-1].(string)]
		if !found {
			return fmt.Errorf("failed to parse entity type %v", value[index-1])
		}
		// Parse display name
		displayName := value[index].(string)
		// Get or create the entity
		entity := tp.CreateEntityIfAbsent(displayName, oid, entityType)
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
