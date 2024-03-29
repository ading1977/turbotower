package topology

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"github.com/turbonomic/turbotower/pkg/influx"
	"github.com/urfave/cli"
	"strconv"
)

type TopologyBuilder struct {
	db      *influx.DBInstance
	context *cli.Context
}

func NewTopologyBuilder(db *influx.DBInstance, context *cli.Context) *TopologyBuilder {
	return &TopologyBuilder{
		db:      db,
		context: context,
	}
}

func (t *TopologyBuilder) Build() (*Topology, error) {
	tp := newTopology()
	if err := t.processCommoditySold(tp); err != nil {
		return nil, err
	}
	if err := t.processCommodityBought(tp); err != nil {
		return nil, err
	}
	tp.buildGraph()
	//tp.PrintGraph()
	//tp.PrintEntityTypeIndex()
	return tp, nil
}

func (t *TopologyBuilder) processCommoditySold(tp *Topology) error {
	commoditySoldTagKeys := []string{"oid", "entity_type", "display_name", "VM_CLUSTER", "HOST_CLUSTER"}
	columns := append(influx.CommoditySoldFieldKeys, commoditySoldTagKeys...)
	row, err := t.db.Query(influx.NewDBQuery(t.context).
		WithColumns(columns...).
		WithName("commodity_sold"))
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
		for i := index+3; i <= len(columns); i++ {
			valObj := value[i]
			if valObj != nil && valObj.(string) != "" {
				groupNames = append(groupNames, valObj.(string))
			}
		}
		// Get or create the entity
		entity := tp.createEntityIfAbsent(displayName, oid, entityType, groupNames...)
		// Parse commodity values
		for i, key := range influx.CommoditySoldFieldKeys {
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
			entity.createCommoditySoldIfAbsent(key, val)
		}
	}
	return nil
}

func (t *TopologyBuilder) processCommodityBought(tp *Topology) error {
	commodityBoughtTagKeys := []string{"oid", "provider_id", "entity_type", "display_name", "VM_CLUSTER", "HOST_CLUSTER"}
	columns := append(influx.CommodityBoughtFieldKeys, commodityBoughtTagKeys...)
	row, err := t.db.Query(influx.NewDBQuery(t.context).
		WithColumns(columns...).
		WithName("commodity_bought"))
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
		for i := index+4; i <= len(columns); i++ {
			valObj := value[i]
			if valObj != nil && valObj.(string) != "" {
				groupNames = append(groupNames, valObj.(string))
			}
		}
		// Get or create the entity
		entity := tp.createEntityIfAbsent(displayName, oid, entityType, groupNames...)
		// Parse commodity values
		for i, key := range influx.CommodityBoughtFieldKeys {
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
			entity.createCommodityBoughtIfAbsent(key, val, providerId)
		}
	}
	// Calculate average commodity bought values
	// Todo: Optimize and move up to the above loop
	for _, entity := range tp.Entities {
		entity.computeAvgBoughtValues()
	}
	return nil
}