package topology

import (
	set "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type Commodity struct {
	Name  string
	Value float64
}

type Entity struct {
	Name               string
	EntityType         int32
	OID                int64
	CommoditySold      map[string]*Commodity
	CommodityBought    map[int64]map[string]*Commodity
	AvgCommBoughtValue map[string]float64
	Providers          []*Entity
	Consumers          []*Entity
	Groups             set.Set
}

type Topology struct {
	Entities        map[int64]*Entity
	EntityTypeIndex map[int32][]*Entity
}

func newCommodity(name string, value float64) *Commodity {
	return &Commodity{name, value}
}

func newEntity(name string, oid int64, entityType int32) *Entity {
	return &Entity{
		Name:               name,
		OID:                oid,
		EntityType:         entityType,
		CommoditySold:      make(map[string]*Commodity),
		CommodityBought:    make(map[int64]map[string]*Commodity),
		AvgCommBoughtValue: make(map[string]float64),
		Groups:             set.NewSet(),
	}
}

func (e *Entity) createCommoditySoldIfAbsent(name string, value float64) {
	if _, found := e.CommoditySold[name]; !found {
		e.CommoditySold[name] = newCommodity(name, value)
	}
}

func (e *Entity) createCommodityBoughtIfAbsent(name string, value float64, providerId int64) {
	if commBought, found := e.CommodityBought[providerId]; found {
		if _, found := commBought[name]; !found {
			// There is no such commodity from this provider, add it to the map
			commBought[name] = newCommodity(name, value)
		}
		return
	}
	// There is no such provider
	e.CommodityBought[providerId] = map[string]*Commodity {
		name: newCommodity(name, value),
	}
}

func (e *Entity) printEntity() {
	entityType, _ := proto.EntityDTO_EntityType_name[e.EntityType]
	log.Infof("OID: %d Type: %s Name: %s", e.OID, entityType, e.Name)
	log.Infof("Belongs to %v", e.Groups)
	log.Infof("Commodity bought:")
	for providerId, commBoughtList := range e.CommodityBought {
		log.Printf("    Provider: %d", providerId)
		log.Printf("        %-40s%-15s", "Metric", "Value")
		for _, commBought := range commBoughtList {
			log.Printf("        %-40s%-15f", commBought.Name, commBought.Value)
		}
	}
	log.Infof("Commodity Sold:")
	log.Printf("        %-40s%-15s", "Metric", "Value")
	for _, commSold := range e.CommoditySold {
		log.Printf("        %-40s%-15f", commSold.Name, commSold.Value)
	}
}

func (e *Entity) getProviderIds() []int64 {
	p := make([]int64, len(e.CommodityBought))
	i := 0
	for k := range e.CommodityBought {
		p[i] = k
		i++
	}
	return p
}

func (e *Entity) computeAvgBoughtValues() {
	l := len(e.CommodityBought)
	for _, commBoughtMap := range e.CommodityBought {
		for name, commBought := range commBoughtMap {
			e.AvgCommBoughtValue[name] += commBought.Value
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("avg[%s]: %+v", commBought.Name, e.AvgCommBoughtValue[commBought.Name])
			}
		}
	}
	for k, v := range e.AvgCommBoughtValue {
		e.AvgCommBoughtValue[k] = v / float64(l)
	}
}

func newTopology() *Topology {
	return &Topology{
		Entities:        make(map[int64]*Entity),
		EntityTypeIndex: make(map[int32][]*Entity),
	}
}

func (t *Topology) createEntityIfAbsent(name string, oid int64, entityType int32, groups ...string) *Entity {
	e, found := t.Entities[oid]
	if !found {
		e = newEntity(name, oid, entityType)
		t.Entities[oid] = e
	}
	for _, group := range groups {
		e.Groups.Add(group)
	}
	return e
}

func (t *Topology) getEntitiesInCluster(clusterName string, entityType int32) []*Entity {
	//log.Infof("Get entities in cluster %v", clusterName)
	var entities []*Entity
	if entityList, found := t.EntityTypeIndex[entityType]; found {
		for _, entity := range entityList {
			//log.Infof("Checking %v %v", entity.Name, entity.Groups)
			if entity.Groups.Contains(clusterName) {
				//log.Infof("Adding entity %v", entity.Name)
				entities = append(entities, entity)
			}
		}
	}
	return entities
}

func (t *Topology) GetContainerPodsInCluster(clusterName string) []*Entity {
	return t.getEntitiesInCluster(clusterName, int32(proto.EntityDTO_CONTAINER_POD))
}

func (t *Topology) GetVirtualMachinesInCluster(clusterName string) []*Entity {
	return t.getEntitiesInCluster(clusterName, int32(proto.EntityDTO_VIRTUAL_MACHINE))
}

func (t *Topology) GetEntityByNameAndType(name string, entityType int32) *Entity {
	if entityList, found := t.EntityTypeIndex[entityType]; found {
		for _, entity := range entityList {
			if entity.Name == name {
				return entity
			}
		}
	}
	return nil
}

func (t *Topology) PrintEntityTypeIndex() {
	log.Infof("%-20s%-15s", "Type", "Count")
	for t, e := range t.EntityTypeIndex {
		entityType, _ := proto.EntityDTO_EntityType_name[t]
		log.Infof("%-20s%-15d", entityType, len(e))
	}
}

func (t *Topology) PrintGraph() {
	for _, e := range t.Entities {
		entityType, _ := proto.EntityDTO_EntityType_name[e.EntityType]
		log.Infof("Entity: %s [%s]", entityType, e.Name)
		log.Infof("    Consumers:")
		for _, consumer := range e.Consumers {
			entityType, _ := proto.EntityDTO_EntityType_name[consumer.EntityType]
			log.Infof("        %s [%s]", entityType, consumer.Name)
		}
		log.Infof("    Providers:")
		for _, provider := range e.Providers {
			entityType, _ := proto.EntityDTO_EntityType_name[provider.EntityType]
			log.Infof("        %s [%s]", entityType, provider.Name)
		}
	}
}

func (e *Entity) addProvider(provider *Entity) {
	e.Providers = append(e.Providers, provider)
}

func (e *Entity) addConsumer(consumer *Entity) {
	e.Consumers = append(e.Consumers, consumer)
}

func (t *Topology) buildGraph() {
	for _, entity := range t.Entities {
		t.EntityTypeIndex[entity.EntityType] = append(t.EntityTypeIndex[entity.EntityType], entity)
		for _, providerId := range entity.getProviderIds() {
			if provider, found := t.Entities[providerId]; found {
				entity.addProvider(provider)
				provider.addConsumer(entity)
			} else {
				if log.GetLevel() >= log.DebugLevel {
					log.Debugf("Cannot locate provider entity with provider ID %s",
						providerId)
				}
			}
		}
	}
}
