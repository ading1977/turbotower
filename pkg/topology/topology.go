package topology

type CommoditySold struct {
	Name     string
	Used     float64
	Capacity float64
}

type CommodityBought struct {
	Name        string
	Used        float64
}

type Entity struct {
	Name            string
	OID             string
	CommoditySold   []*CommoditySold
	CommodityBought map[string][]*CommodityBought
}

type Topology struct {
	Entities map[string]*Entity
}

func NewCommoditySold(name string, used, capacity float64) *CommoditySold {
	return &CommoditySold{name, used, capacity}
}

func NewCommodityBought(name string, used float64) *CommodityBought {
	return &CommodityBought{name, used}
}

func NewEntity(name, oid string) *Entity {
	return &Entity{
		Name: name,
		OID: oid,
		CommodityBought: make(map[string][]*CommodityBought),
	}
}

func (e *Entity) AddCommoditySold(c *CommoditySold) *Entity {
	e.CommoditySold = append(e.CommoditySold, c)
	return e
}

func (e *Entity) AddCommodityBought(c *CommodityBought, provider string) *Entity {
	e.CommodityBought[provider] = append(e.CommodityBought[provider], c)
	return e
}

func NewTopology() *Topology {
	return &Topology{
		Entities: make(map[string]*Entity),
	}
}

func (t *Topology) AddEntity(e *Entity) {
	if _, found := t.Entities[e.OID]; found {
		return
	}
	t.Entities[e.OID] = e
}
