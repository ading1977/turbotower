package topology

import (
	set "github.com/deckarep/golang-set"
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbo-go-sdk/pkg/proto"
	"strings"
)

type neighborFunc func(e *Entity) []*Entity

func GetProviders(e *Entity) []*Entity {
	return e.Providers
}

func GetConsumers(e *Entity) []*Entity {
	return e.Consumers
}

type SupplyChainNode struct {
	EntityType             int32
	Depth                  int
	Members                set.Set
	ConnectedProviderTypes set.Set
	ConnectedConsumerTypes set.Set
}

type SupplyChainResolver struct {
	VisitedEntityTypes set.Set
	VisitedEntities    set.Set
	NodeMap            map[int32]*SupplyChainNode
	Frontier           []*Entity
}

func NewSupplyChainNode(entityType int32, depth int) *SupplyChainNode {
	return &SupplyChainNode{
		EntityType:             entityType,
		Depth:                  depth,
		Members:                set.NewSet(),
		ConnectedProviderTypes: set.NewSet(),
		ConnectedConsumerTypes: set.NewSet(),
	}
}

func (n *SupplyChainNode) addMember(entity *Entity) {
	n.Members.Add(entity)
}

func (n *SupplyChainNode) printNode() {
	log.Infof("Depth: %d", n.Depth)
	var providerTypes, consumerTypes, members []string
	for providerType := range n.ConnectedProviderTypes.Iterator().C {
		entityType, _ := proto.EntityDTO_EntityType_name[providerType.(int32)]
		providerTypes = append(providerTypes, entityType)
	}
	for consumerType := range n.ConnectedConsumerTypes.Iterator().C {
		entityType, _ := proto.EntityDTO_EntityType_name[consumerType.(int32)]
		consumerTypes = append(consumerTypes, entityType)
	}
	log.Infof("Provider types: %s", strings.Join(providerTypes, " "))
	log.Infof("Consumer types: %s", strings.Join(consumerTypes, " "))
	for member := range n.Members.Iterator().C {
		members = append(members, member.(*Entity).Name)
	}
	log.Infof("Members: %s", strings.Join(members, " "))
	log.Infof("Member count: %d", len(members))
}

func NewSupplyChainResolver() *SupplyChainResolver {
	return &SupplyChainResolver{
		VisitedEntityTypes: set.NewSet(),
		VisitedEntities:    set.NewSet(),
		NodeMap:            make(map[int32]*SupplyChainNode),
	}
}

func (s *SupplyChainResolver) GetSupplyChainNodes(startingVertices []*Entity) {
	s.Frontier = startingVertices
	log.Infof("Collect supply chain providers")
	// Collect supply chain providers
	s.traverseSupplyChain(GetProviders, 1, 1)
	// Collect supply chain consumers
	log.Infof("Collect supply chain consumers")
	var frontier []*Entity
	for _, vertex := range startingVertices {
		for _, neighbor := range GetConsumers(vertex) {
			frontier = append(frontier, neighbor)
		}
	}
	s.Frontier = frontier
	s.traverseSupplyChain(GetConsumers, 0, -1)
	s.collectNodeProviderConsumerTypes()
	s.PrintNodeMap()
}

func (s *SupplyChainResolver) traverseSupplyChain(neighborFunc neighborFunc,
	currentDepth int, increment int) {
	var nextFrontier []*Entity
	var visitedEntityTypesInThisDepth = set.NewSet()
	// Process the current depth
	//log.Infof("Current depth %d", currentDepth)
	for len(s.Frontier) > 0 {
		// Dequeue
		vertex := s.Frontier[0]
		s.Frontier = s.Frontier[1:]
		if s.VisitedEntities.Contains(vertex) {
			continue
		}
		s.VisitedEntities.Add(vertex)
		//log.Infof("Visiting %s", vertex.Name)
		// Only add a node when we have not already visited an entity of the same type
		if !s.VisitedEntityTypes.Contains(vertex.EntityType) {
			neighbors := neighborFunc(vertex)
			for _, neighbor := range neighbors {
				if !s.VisitedEntities.Contains(neighbor) {
					nextFrontier = append(nextFrontier, neighbor)
				}
			}
			node, found := s.NodeMap[vertex.EntityType]
			if !found {
				entityType, _ := proto.EntityDTO_EntityType_name[vertex.EntityType]
				log.Infof("Create a new supply chain node for %s", entityType)
				node = NewSupplyChainNode(vertex.EntityType, currentDepth)
				s.NodeMap[vertex.EntityType] = node
			}
			//log.Infof("Adding a member %s to node type %v", vertex.Name, vertex.EntityType)
			node.addMember(vertex)
			visitedEntityTypesInThisDepth.Add(vertex.EntityType)
		}
	}
	s.VisitedEntityTypes = s.VisitedEntityTypes.Union(visitedEntityTypesInThisDepth)
	//log.Infof("Visited entity types %v", s.Context.VisitedEntityTypes)
	// Process the next depth
	if len(nextFrontier) > 0 {
		s.Frontier = nextFrontier
		s.traverseSupplyChain(neighborFunc, currentDepth+increment, increment)
	}
}

func (s *SupplyChainResolver) collectNodeProviderConsumerTypes() {
	for _, node := range s.NodeMap {
		for member := range node.Members.Iterator().C {
			for _, provider := range member.(*Entity).Providers {
				node.ConnectedProviderTypes.Add(provider.EntityType)
			}
			for _, consumer := range member.(*Entity).Consumers {
				node.ConnectedConsumerTypes.Add(consumer.EntityType)
			}
		}
	}
}

func (s *SupplyChainResolver) PrintNodeMap() {
	for eType, node := range s.NodeMap {
		entityType, _ := proto.EntityDTO_EntityType_name[eType]
		log.Infof("Entity Type: %s", entityType)
		node.printNode()
		log.Println()
	}
}
