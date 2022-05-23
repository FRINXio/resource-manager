package src

import (
	"fmt"
	"github.com/pkg/errors"
	"sort"
)

type VlanRange struct {
	currentResources       []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

func NewVlanRange(currentResources []map[string]interface{}, resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) VlanRange {
	return VlanRange{currentResources, resourcePoolProperties, userInput}
}

type VlanRangeValues struct {
	from int
	to   int
}

func NewVlanRangeValues(from int, to int) VlanRangeValues {
	return VlanRangeValues{from: from, to: to}
}

func (vlanRange *VlanRange) Invoke() (map[string]interface{}, error) {
	if vlanRange.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract parent vlan range from pool name")
	}
	parentRange := make(map[string]interface{})
	for k, v := range vlanRange.resourcePoolProperties {
		parentRange[k] = v
	}

	from, ok := parentRange["from"]
	if !ok {
		return nil, errors.New("Missing from in parentRange")
	}
	to, ok := parentRange["to"]
	if !ok {
		return nil, errors.New("Missing to in parentRange")
	}
	desiredSize, ok := vlanRange.userInput["desiredSize"]
	if !ok {
		return nil, errors.New("Missing desiredSize in userInput")
	}

	from, err := NumberToInt(from)
	if err != nil {
		return nil, err
	}
	to, err = NumberToInt(to)
	if err != nil {
		return nil, err
	}
	desiredSize, err = NumberToInt(desiredSize)
	if err != nil {
		return nil, err
	}

	if desiredSize.(int) < 1 {
		return nil, errors.New("Unable to allocate VLAN range from: " + from.(string) + " - " + to.(string) +
			". Desired size is invalid: " + desiredSize.(string) + ". Use values >= 1")
	}

	currentResourcesVlanRangesValues := FromCurrentResourcesToVlanRangeValues(vlanRange.currentResources)

	sort.SliceStable(currentResourcesVlanRangesValues, func(i, j int) bool {
		return currentResourcesVlanRangesValues[i].to < currentResourcesVlanRangesValues[j].to
	})
	findingAvailableRange := VlanRangeValues{}
	findingAvailableRange.from = from.(int)
	// iterate over allocated ranges and see if a desired new range can be squeezed in
	for _, allocatedRange := range currentResourcesVlanRangesValues {
		// set to bound to from bound of next range
		findingAvailableRange.to = allocatedRange.from - 1
		// if there is enough space, allocate a chunk of that range
		if RangeCapacity(findingAvailableRange) >= desiredSize.(int) {
			findingAvailableRange.to = findingAvailableRange.from + desiredSize.(int) - 1
			// FIXME How to pass these stats ?
			// logStats(findingAvailableRange, parentRange, currentResourcesUnwrapped)
			vlanProperties := make(map[string]interface{})
			vlanProperties["from"] = findingAvailableRange.from
			vlanProperties["to"] = findingAvailableRange.to
			return vlanProperties, nil
		}

		findingAvailableRange.from = allocatedRange.to + 1
		findingAvailableRange.to = allocatedRange.to + 1
	}

	// check if there is some space left at the end of parent range
	findingAvailableRange.to = to.(int)
	if RangeCapacity(findingAvailableRange) >= desiredSize.(int) {
		findingAvailableRange.to = findingAvailableRange.from + desiredSize.(int) - 1
		// FIXME How to pass these stats ?
		// logStats(findingAvailableRange, parentRange, currentResourcesUnwrapped)
		vlanProperties := make(map[string]interface{})
		vlanProperties["from"] = findingAvailableRange.from
		vlanProperties["to"] = findingAvailableRange.to
		return vlanProperties, nil
	}
	return nil, errors.New("Unable to allocate VLAN range: " + fmt.Sprint(from) + " - " + fmt.Sprint(to) +
		". Insufficient capacity to allocate a new range of size: " + fmt.Sprint(desiredSize) + ".")
}

func RangeCapacity(vlanRange VlanRangeValues) int {
	return vlanRange.to - vlanRange.from + 1
}

func FromCurrentResourcesToVlanRangeValues(currentResources []map[string]interface{}) []VlanRangeValues {
	var currentResourcesUnwrapped []map[string]interface{}
	for _, element := range currentResources {
		value, ok := element["Properties"]
		if ok {
			currentResourcesUnwrapped = append(currentResourcesUnwrapped, value.(map[string]interface{}))
		}
	}

	var currentResourcesVlanRangesValues []VlanRangeValues
	for _, element := range currentResourcesUnwrapped {
		from, ok := element["from"]
		if ok {
			from, err := NumberToInt(from)
			if err != nil {
				return nil
			}
			to, ok := element["to"]
			if ok {
				to, err := NumberToInt(to)
				if err != nil {
					return nil
				}
				currentResourcesVlanRangesValues = append(currentResourcesVlanRangesValues,
					NewVlanRangeValues(from.(int), to.(int)))
			}
		}
	}
	return currentResourcesVlanRangesValues
}

func FreeRangeCapacity(vlanRange map[string]interface{}, utilisedCapacity int) (int, error) {
	from, err := NumberToInt(vlanRange["from"])
	if err != nil {
		return 0, err
	}
	to, err := NumberToInt(vlanRange["to"])
	if err != nil {
		return 0, err
	}
	return to.(int) - from.(int) + 1 - utilisedCapacity, nil
}

func UtilizedRangeCapacity(currentResources []map[string]interface{}, newlyAllocatedRangeCapacity int) int {
	var allocatedCapacity = 0
	currentResourcesVlanRangesValues := FromCurrentResourcesToVlanRangeValues(currentResources)
	for _, allocatedRange := range currentResourcesVlanRangesValues {
		allocatedCapacity += RangeCapacity(allocatedRange)
	}
	return allocatedCapacity + newlyAllocatedRangeCapacity
}

func (vlanRange *VlanRange) Capacity() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	allocatedCapacity := UtilizedRangeCapacity(vlanRange.currentResources, 0)
	freeCapacity, err := FreeRangeCapacity(vlanRange.resourcePoolProperties, allocatedCapacity)
	if err != nil {
		return nil, err
	}
	result["freeCapacity"] = float64(freeCapacity)
	result["utilizedCapacity"] = float64(allocatedCapacity)
	return result, nil
}
