package src

import "github.com/pkg/errors"

// STRATEGY_START

type Vlan struct {
	currentResources []map[string]interface{}
	resourcePoolProperties map[string]interface{}
}

func NewVlan(currentResources []map[string]interface{}, resourcePoolProperties map[string]interface{}) Vlan {
	return Vlan{currentResources, resourcePoolProperties}
}

func UtilizedCapacity(allocatedRanges []map[string]interface{}, newlyAllocatedVlan int) int {
	return len(allocatedRanges) + newlyAllocatedVlan
}

func FreeCapacity(vlanRange map[string]interface{}, utilisedCapacity int) int {
	return vlanRange["to"].(int) - vlanRange["from"].(int) + 1 - utilisedCapacity
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
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
	var currentResourcesSet []int
	for _, element := range currentResourcesUnwrapped {
		value, ok := element["vlan"]
		if ok {
			currentResourcesSet = append(currentResourcesSet, value.(int))
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
	for i := from.(int); i <= to.(int); i++ {
		if !contains(currentResourcesSet, i) {
			// FIXME How to pass these stats ?
			// logStats(i, parentRange, currentResourcesUnwrapped)
			vlanProperties := make(map[string]interface{})
			vlanProperties["vlan"] = i
			return vlanProperties, nil
		}
	}
	return nil, errors.New("Unable to allocate VLAN. Insufficient capacity to allocate a new vlan")
}

func (vlan *Vlan) Capacity() map[string]interface{} {
	var result = make(map[string]interface{})
	result["freeCapacity"] = FreeCapacity(vlan.resourcePoolProperties, len(vlan.currentResources))
	result["utilizedCapacity"] = len(vlan.currentResources)
	return result
}

// STRATEGY_END
