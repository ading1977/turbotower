package utils

import (
	"github.com/turbonomic/turbotower/pkg/topology"
)

func GetSoldValues(commSold map[string]*topology.Commodity) map[string]float64 {
	soldValues := make(map[string]float64)
	for name, commSold := range commSold {
		soldValues[name] = commSold.Value

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
