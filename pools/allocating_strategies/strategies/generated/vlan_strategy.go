package src

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type Vlan struct {
	currentResources []map[string]interface{}
	resourcePoolProperties map[string]interface{}
}

func NewVlan(currentResources []map[string]interface{}, resourcePoolProperties map[string]interface{}) Vlan {
	return Vlan{currentResources, resourcePoolProperties}
}

func UtilizedCapacity(allocatedRanges []map[string]interface{}, newlyAllocatedVlan float64) float64 {
	return float64(len(allocatedRanges)) + newlyAllocatedVlan
}

func FreeCapacity(vlanRange map[string]interface{}, utilisedCapacity float64) float64 {
	return float64(vlanRange["to"].(int) - vlanRange["from"].(int) + 1 - int(utilisedCapacity))
}

func contains(slice []float64, val float64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func NumberToInt(number interface{}) (interface{}, error) {
	var newNumber int
	switch number.(type) {
	case json.Number:
		intVal64, err := number.(json.Number).Float64()
		if err != nil {
			return nil, errors.New("Unable to convert a json number")
		}
		newNumber = int(intVal64)
	case float64:
		newNumber = int(number.(float64))
	case int:
		newNumber = number.(int)
	}
	return newNumber, nil
}

func (vlan *Vlan) Invoke() (map[string]interface{}, error) {
	if vlan.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract parent vlan range from pool name")
	}
	parentRange := make(map[string]interface{})
	for k, v := range vlan.resourcePoolProperties {
		parentRange[k] = v
	}

	var currentResourcesUnwrapped []map[string]interface{}
	for _, element := range vlan.currentResources {
		value, ok := element["Properties"]
		if ok {
			currentResourcesUnwrapped = append(currentResourcesUnwrapped, value.(map[string]interface{}))
		}
	}
	var currentResourcesSet []float64
	for _, element := range currentResourcesUnwrapped {
		value, ok := element["vlan"]
		if ok {
			currentResourcesSet = append(currentResourcesSet, value.(float64))
		}
	}
	from, ok := parentRange["from"]
	if !ok {
		return nil, errors.New("Missing from in parentRange")
	}
	to, ok := parentRange["to"]
	if !ok {
		return nil, errors.New("Missing to in parentRange")
	}

	from, err := NumberToInt(from)
	if err != nil {
		return nil, err
	}
	to, err = NumberToInt(to)
	if err != nil {
		return nil, err
	}

	for i := from.(int); i <= to.(int); i++ {
		if !contains(currentResourcesSet, float64(i)) {
			// FIXME How to pass these stats ?
			// logStats(i, parentRange, currentResourcesUnwrapped)
			vlanProperties := make(map[string]interface{})
			vlanProperties["vlan"] = float64(i)
			return vlanProperties, nil
		}
	}
	return nil, errors.New("Unable to allocate VLAN. Insufficient capacity to allocate a new vlan")
}

func (vlan *Vlan) Capacity() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	result["freeCapacity"] = FreeCapacity(vlan.resourcePoolProperties, float64(len(vlan.currentResources)))
	result["utilizedCapacity"] = float64(len(vlan.currentResources))
	return result, nil
}
