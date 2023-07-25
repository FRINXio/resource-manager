package src

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type UniqueId struct {
	currentResources       []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

func NewUniqueId(
	currentResources []map[string]interface{},
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) UniqueId {
	return UniqueId{currentResources, resourcePoolProperties, userInput}
}

func (uniqueId *UniqueId) getNextFreeCounter(fromValue int, toValue int, desiredValue int, allocatedResources []map[string]interface{}) (int, error) {
	if desiredValue != -1 {
		if desiredValue >= fromValue && desiredValue <= toValue {
			for _, allocResource := range allocatedResources {
				counter, ok := allocResource["counter"]

				if !ok {
					continue
				}

				if counter == desiredValue {
					return -1, errors.New("desiredValue is already allocated")
				}
			}

			return desiredValue, nil
		} else {
			return -1, errors.New("desiredValue is not in range")
		}
	}

	for i := fromValue; i <= toValue; i++ {
		isAllocated := false
		for _, allocResource := range allocatedResources {
			counter, ok := allocResource["counter"]

			if !ok {
				continue
			}

			if counter == i {
				isAllocated = true
				break
			}
		}

		if !isAllocated {
			return i, nil
		}
	}

	return -1, nil
}

func (uniqueId *UniqueId) Invoke() (map[string]interface{}, error) {
	if uniqueId.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	var fromValue = 0 // max int
	resourcePoolfromValue, ok := uniqueId.resourcePoolProperties["from"]
	if ok {
		resourcePoolfromValue, _ := NumberToInt(resourcePoolfromValue)
		fromValue = resourcePoolfromValue.(int)
	} else {
		return nil, errors.New("Missing property 'from' in resource pool.")
	}

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
	} else {
		return nil, errors.New("Missing property 'to' in resource pool.")
	}

	replacePoolProperties := make(map[string]interface{})
	for k, v := range uniqueId.resourcePoolProperties {
		if k != "idFormat" && k != "counterFormatWidth" && k != "from" && k != "to" {
			replacePoolProperties[k] = v
		}
	}
	desiredValue := -1
	if value, ok := uniqueId.userInput["desiredValue"]; ok {
		value, _ = NumberToInt(value)
		desiredValue = value.(int)
		if desiredValue > int(int64(toValue)) {
			return nil, errors.New("Unable to allocate Unique-id desiredValue: " + strconv.FormatInt(int64(value.(int)), 10) + "." +
				" Value is out of scope: " + strconv.FormatInt(int64(toValue), 10))
		}
		if desiredValue < fromValue {
			return nil, errors.New("Unable to allocate Unique-id desiredValue: " + strconv.FormatInt(int64(value.(int)), 10) + "." +
				" Value is out of scope: " + strconv.FormatInt(int64(fromValue), 10))
		}
	}

	nextFreeCounter, err := uniqueId.getNextFreeCounter(fromValue, int(toValue), desiredValue, uniqueId.currentResources)
	if err != nil {
		return nil, err
	}
	if prefixNumber, ok := uniqueId.resourcePoolProperties["counterFormatWidth"]; ok {
		replacePoolProperties["counter"] = fmt.Sprintf(
			"%0"+strconv.Itoa(prefixNumber.(int))+"d", int(nextFreeCounter))
	} else {
		replacePoolProperties["counter"] = nextFreeCounter
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

		idFormat = strings.Replace(idFormat.(string), "{"+k+"}", v.(string), 1)
	}

	var result = make(map[string]interface{})
	result["text"] = idFormat
	result["counter"] = nextFreeCounter
	return result, nil
}

func (uniqueId *UniqueId) Capacity() (map[string]interface{}, error) {
	var result = make(map[string]interface{})

	if uniqueId.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resource pool properties")
	}

	var fromValue = 0
	resourcePoolfromValue, ok := uniqueId.resourcePoolProperties["from"]

	if ok {
		resourcePoolFromValue, _ := NumberToInt(resourcePoolfromValue)
		fromValue = resourcePoolFromValue.(int)
	} else {
		return nil, errors.New("Missing property 'from' in resource pool.")
	}

	var toValue = int(^uint(0) >> 1) // max int
	resourcePoolToValue, ok := uniqueId.resourcePoolProperties["to"]

	if ok {
		resourcePoolToValueInt, err := NumberToInt(resourcePoolToValue)

		if err != nil {
			return nil, errors.New("Unable to convert toValue to int")
		}

		toValue = resourcePoolToValueInt.(int)

		result["utilizedCapacity"] = len(uniqueId.currentResources)
		result["freeCapacity"] = toValue - fromValue - len(uniqueId.currentResources)
	} else {
		result["utilizedCapacity"] = len(uniqueId.currentResources)
		result["freeCapacity"] = toValue - fromValue - len(uniqueId.currentResources)
	}

	return result, nil
}
