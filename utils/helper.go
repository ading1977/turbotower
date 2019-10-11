package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbotower/pkg/topology"
)

func GetAvgBoughtUsed(commBoughtMap map[string][]*topology.CommodityBought) map[string]float64 {
	avgUsed := make(map[string]float64)
	l := len(commBoughtMap)
	for _, commBoughtList := range commBoughtMap {
		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("commBoughtList: %+v", commBoughtList)
		}
		for _, commBought := range commBoughtList {
			avgUsed[commBought.Name] += commBought.Used
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("avgUsed[%s]: %+v", commBought.Name, avgUsed[commBought.Name])
			}
		}
	}
	for k, v := range avgUsed {
		avgUsed[k] = v/float64(l)
	}
	if log.GetLevel() >= log.DebugLevel {
		log.Debugf("Entity created: %+v", avgUsed)
	}
	return avgUsed
}

func GetSoldUsedUtil(commSoldList []*topology.CommoditySold) map[string]float64 {
	soldUsedUtil := make(map[string]float64)
	for _, commSold := range commSoldList {
		soldUsedUtil[commSold.Name + "_USED"] = commSold.Used
		soldUsedUtil[commSold.Name + "_UTIL"] = commSold.Used/commSold.Capacity * 100
	}
	return soldUsedUtil
}

func Truncate(s string, max_len int) string {
	l := len(s)
	if l < max_len {
		return s
	}
	return "*" + s[(l-max_len+2):]
}
