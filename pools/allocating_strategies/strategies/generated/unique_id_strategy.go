package src

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type UniqueId struct {
	currentResources []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput map[string]interface{}
}

func NewUniqueId(currentResources []map[string]interface{},
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) UniqueId {
	return UniqueId{currentResources, resourcePoolProperties, userInput}
}

func (uniqueId *UniqueId) getNextFreeCounterAndResourcesSet(fromValue int) (float64, map[float64]bool) {
	var start = float64(fromValue)
	currentResourcesSet := make(map[float64]bool)

	for _, element := range uniqueId.currentResources {
		var properties = element["Properties"].(map[string]interface{})
		var counter = properties["counter"].(float64)
		currentResourcesSet[counter] = true
	}
	for i := start; i < float64(len(uniqueId.currentResources)) + start; i++ {
		if !currentResourcesSet[i + 1] {
			return i + 1, currentResourcesSet
		}
	}
	return float64(len(uniqueId.currentResources)) + start + 1, currentResourcesSet
}

func (uniqueId *UniqueId) Invoke() (map[string]interface{}, error) {
	if uniqueId.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	var fromValue = 0 // max int
	resourcePoolfromValue, ok := uniqueId.resourcePoolProperties["from"]
	if ok {
		resourcePoolfromValue, _ := NumberToInt(resourcePoolfromValue)
		fromValue = resourcePoolfromValue.(int) - 1
	}

	nextFreeCounter, currentResourcesSet := uniqueId.getNextFreeCounterAndResourcesSet(fromValue)
	idFormat, ok := uniqueId.resourcePoolProperties["idFormat"]
	if !ok {
		return nil, errors.New("Missing idFormat in resources")
	}
	if !strings.Contains(idFormat.(string), "{counter}") {
		return nil, errors.New("Missing {counter} in idFormat")
	}

	var toValue = ^uint(0) >> 1 // max int
	resourcePooltoValue, ok := uniqueId.resourcePoolProperties["to"]
	if ok {
		resourcePooltoValue, _ := NumberToInt(resourcePooltoValue)
		toValue = uint(resourcePooltoValue.(int))
	}

	replacePoolProperties := make(map[string]interface{})
	for k, v := range uniqueId.resourcePoolProperties {
		if k != "idFormat" && k != "counterFormatWidth" && k != "from" && k != "to" {
			replacePoolProperties[k] = v
		}
	}

	if value, ok := uniqueId.userInput["desiredValue"]; ok {
		value, _ = NumberToInt(value)
		nextFreeCounter = float64(value.(int))
		if nextFreeCounter > float64(toValue) {
			return nil, errors.New("Unable to allocate Unique-id desiredValue: " + strconv.FormatInt(int64(value.(int)), 10) + "." +
				" Value is out of scope: " + strconv.FormatInt(int64(toValue), 10))
		}
		if nextFreeCounter < float64(fromValue) {
			return nil, errors.New("Unable to allocate Unique-id desiredValue: " + strconv.FormatInt(int64(value.(int)), 10) + "." +
				" Value is out of scope: " + strconv.FormatInt(int64(fromValue), 10))
		}
	}
	if prefixNumber, ok := uniqueId.resourcePoolProperties["counterFormatWidth"]; ok {
		replacePoolProperties["counter"] = fmt.Sprintf(
			"%0" + strconv.Itoa(prefixNumber.(int)) +"d", int(nextFreeCounter))
	} else {
		replacePoolProperties["counter"] = nextFreeCounter
	}

	if nextFreeCounter > float64(toValue) {
		return nil, errors.New("Unable to allocate Unique-id from idFormat: \"" + idFormat.(string) + "\"." +
			" Insufficient capacity to allocate a unique-id.")
	}

	for k, v := range replacePoolProperties {
		switch v.(type) {
		case float64:
			v = fmt.Sprint(v.(float64))
		case int:
			v = strconv.Itoa(v.(int))
		case json.Number:
			intVal64, err := v.(json.Number).Float64()
			if err != nil {
				return nil, errors.New("Unable to convert a json number")
			}
			v = fmt.Sprint(intVal64)
		}

		idFormat = strings.Replace(idFormat.(string), "{" + k + "}", v.(string), 1)
	}
	if currentResourcesSet[nextFreeCounter] {
		return nil, errors.New("Unique-id " + idFormat.(string) + " was already claimed." )
	}

	var result = make(map[string]interface{})
	result["text"] = idFormat
	result["counter"] = nextFreeCounter
	return result, nil
}

func (uniqueId *UniqueId) Capacity() (map[string]interface{}, error) {
	var allocatedCapacity = float64(len(uniqueId.currentResources))
	var result = make(map[string]interface{})
	var fromValue float64
	var toValue float64
	to, ok := uniqueId.resourcePoolProperties["to"]
	if ok {
		toValue = float64(to.(int))
	} else {
		toValue = float64(^uint(0) >> 1)
	}
	from, ok := uniqueId.resourcePoolProperties["from"]
	if ok {
		fromValue = float64(from.(int))
	} else {
		fromValue = float64(0)
	}
	result["freeCapacity"] = toValue - allocatedCapacity - fromValue + 1
	result["utilizedCapacity"] = allocatedCapacity
	return result, nil
}
