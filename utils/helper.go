package utils

import (
	log "github.com/sirupsen/logrus"
	"github.com/turbonomic/turbotower/pkg/topology"
)

func GetAvgBoughtValues(commBoughtMap map[int64][]*topology.Commodity) map[string]float64 {
	avg := make(map[string]float64)
	l := len(commBoughtMap)
	for _, commBoughtList := range commBoughtMap {
		if log.GetLevel() >= log.DebugLevel {
			log.Debugf("commBoughtList: %+v", commBoughtList)
		}
		for _, commBought := range commBoughtList {
			avg[commBought.Name] += commBought.Value
			if log.GetLevel() >= log.DebugLevel {
				log.Debugf("avg[%s]: %+v", commBought.Name, avg[commBought.Name])
			}
		}
	}
	for k, v := range avg {
		avg[k] = v / float64(l)
	}
	return avg
}


func GetSoldValues(commSoldList []*topology.Commodity) map[string]float64 {
	soldValues := make(map[string]float64)
	for _, commSold := range commSoldList {
		soldValues[commSold.Name] = commSold.Value

	}
	return soldValues
}

func Truncate(s string, max_len int) string {
	l := len(s)
	if l < max_len {
		return s
	}
	return "*" + s[(l-max_len+2):]
}
