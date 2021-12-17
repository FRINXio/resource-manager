package src

import (
	"encoding/json"
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

type vlanRangeValues struct {
	from float64
	to   float64
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

	from, err := NumberToFloat(from)
	if err != nil {
		return nil, err
	}
	to, err = NumberToFloat(to)
	if err != nil {
		return nil, err
	}
	desiredSize, err = NumberToFloat(desiredSize)
	if err != nil {
		return nil, err
	}

	if desiredSize.(float64) < float64(1) {
		return nil, errors.New("Unable to allocate VLAN range from: " + from.(string) + " - " + to.(string) +
			". Desired size is invalid: " + desiredSize.(string) + ". Use values >= 1")
	}

	currentResourcesVlanRangesValues := FromCurrentResourcesToVlanRangeValues(vlanRange.currentResources)

	sort.SliceStable(currentResourcesVlanRangesValues, func(i, j int) bool {
		return currentResourcesVlanRangesValues[i].to < currentResourcesVlanRangesValues[j].to
	})
	findingAvailableRange := vlanRangeValues{}
	findingAvailableRange.from = from.(float64)
	// iterate over allocated ranges and see if a desired new range can be squeezed in
	for _, allocatedRange := range currentResourcesVlanRangesValues {
		// set to bound to from bound of next range
		findingAvailableRange.to = allocatedRange.from - float64(1)
		// if there is enough space, allocate a chunk of that range
		if RangeCapacity(findingAvailableRange) >= desiredSize.(float64) {
			findingAvailableRange.to = findingAvailableRange.from + desiredSize.(float64) - float64(1)
			// FIXME How to pass these stats ?
			// logStats(findingAvailableRange, parentRange, currentResourcesUnwrapped)
			vlanProperties := make(map[string]interface{})
			vlanProperties["from"] = findingAvailableRange.from
			vlanProperties["to"] = findingAvailableRange.to
			return vlanProperties, nil
		}

		findingAvailableRange.from = allocatedRange.to + float64(1)
		findingAvailableRange.to = allocatedRange.to + float64(1)
	}

	// check if there is some space left at the end of parent range
	findingAvailableRange.to = to.(float64)
	if RangeCapacity(findingAvailableRange) >= desiredSize.(float64) {
		findingAvailableRange.to = findingAvailableRange.from + desiredSize.(float64) - float64(1)
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

func NumberToFloat(number interface{}) (interface{}, error) {
	var newNumber float64
	switch number.(type) {
	case json.Number:
		floatVal64, err := number.(json.Number).Float64()
		if err != nil {
			return nil, errors.New("Unable to convert a json number")
		}
		newNumber = floatVal64
	case float64:
		newNumber = number.(float64)
	case int:
		newNumber = float64(number.(int))
	}
	return newNumber, nil
}

func RangeCapacity(vlanRange vlanRangeValues) float64 {
	return vlanRange.to - vlanRange.from + float64(1)
}

func FromCurrentResourcesToVlanRangeValues(currentResources []map[string]interface{}) []vlanRangeValues {
	var currentResourcesUnwrapped []map[string]interface{}
	for _, element := range currentResources {
		value, ok := element["Properties"]
		if ok {
			currentResourcesUnwrapped = append(currentResourcesUnwrapped, value.(map[string]interface{}))
		}
	}

	var currentResourcesVlanRangesValues []vlanRangeValues
	for _, element := range currentResourcesUnwrapped {
		from, ok := element["from"]
		if ok {
			to, ok := element["to"]
			if ok {
				currentResourcesVlanRangesValues = append(currentResourcesVlanRangesValues,
					vlanRangeValues{from: from.(float64), to: to.(float64)})
			}
		}
	}
	return currentResourcesVlanRangesValues
}

func FreeCapacityFloat(vlanRange map[string]interface{}, utilisedCapacity float64) float64 {
	return vlanRange["to"].(float64) - vlanRange["from"].(float64) + float64(1) - utilisedCapacity
}

func (vlanRange *VlanRange) Capacity() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	var allocatedCapacity = float64(0)
	currentResourcesVlanRangesValues := FromCurrentResourcesToVlanRangeValues(vlanRange.currentResources)
	for _, allocatedRange := range currentResourcesVlanRangesValues {
		allocatedCapacity += RangeCapacity(allocatedRange)
	}
	result["freeCapacity"] = FreeCapacityFloat(vlanRange.resourcePoolProperties, allocatedCapacity)
	result["utilizedCapacity"] = allocatedCapacity
	return result, nil
}
