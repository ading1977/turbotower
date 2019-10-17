package topology

import (
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
)

type Commodity struct {
	Name  string
	Value float64
}

type Entity struct {
	Name            string
	EntityType      int32
	OID             int64
	CommoditySold   []*Commodity
	CommodityBought map[int64][]*Commodity
	Providers       []*Entity
	Consumers       []*Entity
}

func NewCommodity(name string, value float64) *Commodity {
	return &Commodity{name, value}
}

func NewEntity(name string, oid int64, entityType int32) *Entity {
	return &Entity{
		Name:            name,
		OID:             oid,
		EntityType:      entityType,
		CommodityBought: make(map[int64][]*Commodity),
	}
}

func (e *Entity) CreateCommoditySoldIfAbsent(name string, value float64) {
	for _, commSold := range e.CommoditySold {
		if commSold.Name == name {
			// There is already a commodity with the same name
			return
		}
	}
	e.CommoditySold = append(e.CommoditySold, NewCommodity(name, value))
}

func (e *Entity) CreateCommodityBoughtIfAbsent(name string, value float64, providerId int64) {
	if commBoughtList, found := e.CommodityBought[providerId]; found {
		for _, commBought := range commBoughtList {
			if commBought.Name == name {
				// There is already a commodity with the same name from the same provider
				return
			}
		}
	}
	// There is no such provider or there is no such commodity from this provider
	e.CommodityBought[providerId] = append(e.CommodityBought[providerId],
		NewCommodity(name, value))
}

func (e *Entity) PrintEntity() {
	entityType, _ := proto.EntityDTO_EntityType_name[e.EntityType]
	log.Infof("OID: %d Type: %s Name: %s", e.OID, entityType, e.Name)
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

func (e *Entity) GetProviderIds() []int64 {
	p := make([]int64, len(e.CommodityBought))
	i := 0
	for k := range e.CommodityBought {
		p[i] = k
		i++
	}
	return p
}

func (e *Entity) AddProvider(provider *Entity) {
	e.Providers = append(e.Providers, provider)
}

func (e *Entity) AddConsumer(consumer *Entity) {
	e.Consumers = append(e.Consumers, consumer)
}

type Topology struct {
	Entities        map[int64]*Entity
	EntityTypeIndex map[int32][]*Entity
}

func NewTopology() *Topology {
	return &Topology{
		Entities:        make(map[int64]*Entity),
		EntityTypeIndex: make(map[int32][]*Entity),
	}
}

func (t *Topology) CreateEntityIfAbsent(name string, oid int64, entityType int32) *Entity {
	e, found := t.Entities[oid]
	if found {
		return e
	}
	e = NewEntity(name, oid, entityType)
	t.Entities[oid] = e
	return e
}

func (t *Topology) GetEntity(oid int64) *Entity {
	e, found := t.Entities[oid]
	if found {
		return e
	}
	return nil
}

func (t *Topology) AddEntity(e *Entity) {
	if _, found := t.Entities[e.OID]; found {
		return
	}
	t.Entities[e.OID] = e
}

func (t *Topology) PrintEntities() {
	for _, entity := range t.Entities {
		entity.PrintEntity()
		log.Println()
	}
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

func (t *Topology) BuildGraph() {
	for _, entity := range t.Entities {
		t.EntityTypeIndex[entity.EntityType] = append(t.EntityTypeIndex[entity.EntityType], entity)
		for _, providerId := range entity.GetProviderIds() {
			if provider, found := t.Entities[providerId]; found {
				entity.AddProvider(provider)
				provider.AddConsumer(entity)
			} else {
				if log.GetLevel() >= log.DebugLevel {
					log.Debugf("Cannot locate provider entity with provider ID %s",
						providerId)
				}
			}
		}
	}
}
