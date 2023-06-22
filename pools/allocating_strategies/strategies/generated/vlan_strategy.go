package src

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"regexp"
	"strconv"
	"time"
)

type Vlan struct {
	currentResources       []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

func NewVlan(currentResources []map[string]interface{},
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) Vlan {
	return Vlan{currentResources, resourcePoolProperties, userInput}
}

func UtilizedCapacity(allocatedRanges []map[string]interface{}, newlyAllocatedVlan float64) float64 {
	return float64(len(allocatedRanges)) + newlyAllocatedVlan
}

func (vlan *Vlan) FreeCapacity(vlanRange map[string]interface{}, utilisedCapacity float64) float64 {
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
	switch number.(type) {
	case json.Number:
		intVal64, err := number.(json.Number).Float64()
		if err != nil {
			return nil, errors.New("Unable to convert a json number")
		}
		return int(intVal64), nil
	case float64:
		return int(number.(float64)), nil
	case int:
		return number.(int), nil
	case int64:
		return int(number.(int64)), nil
	case int32:
		return int(number.(int32)), nil
	case big.Int:
		b := number.(big.Int)
		val := int(b.Int64())
		return val, nil
	}
	return nil, errors.New("Unable to convert number: " + number.(string) + " to a known type")
}

func (vlan *Vlan) Invoke() (map[string]interface{}, error) {
	time.Sleep(5 * time.Second)
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

	if value, ok := vlan.userInput["desiredValue"]; ok {
		re := regexp.MustCompile(`^\d+$`)
		if re.MatchString(value.(string)) {
			desiredValueNum, err := strconv.Atoi(value.(string))
			if err != nil {
				return nil, err
			}
			if desiredValueNum >= from.(int) && desiredValueNum < to.(int) {
				if contains(currentResourcesSet, float64(desiredValueNum)) {
					return nil, errors.New("VLAN " + value.(string) + " was already claimed.")
				}
				vlanProperties := make(map[string]interface{})
				vlanProperties["vlan"] = float64(desiredValueNum)
				return vlanProperties, nil
			} else {
				return nil, errors.New("VLAN " + value.(string) + " is out of range: " +
					strconv.Itoa(from.(int)) + " - " + strconv.Itoa(to.(int)))
			}
		} else {
			return nil, errors.New("VLAN must be a number. Received: " + value.(string) + ".")
		}
	}

	for i := from.(int); i <= to.(int); i++ {
		if !contains(currentResourcesSet, float64(i)) {
			// FIXME How to pass these stats
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
	freeCapacity := vlan.FreeCapacity(vlan.resourcePoolProperties, float64(len(vlan.currentResources)))
	result["freeCapacity"] = fmt.Sprintf("%v", freeCapacity)
	result["utilizedCapacity"] = strconv.Itoa(len(vlan.currentResources))
	return result, nil
}
