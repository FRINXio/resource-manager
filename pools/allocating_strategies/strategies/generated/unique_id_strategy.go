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
}

func NewUniqueId(currentResources []map[string]interface{}, resourcePoolProperties map[string]interface{}) UniqueId {
	return UniqueId{currentResources, resourcePoolProperties}
}

func (uniqueId *UniqueId) getNextFreeCounter() float64 {
	var max float64
	value, ok := uniqueId.resourcePoolProperties["from"]
	switch value.(type) {
	case json.Number:
		value, _ = value.(json.Number).Float64()
	case float64:
		value = value.(float64)
	case int:
		value = float64(value.(int))
	}
	if ok {
		max = value.(float64) - 1
	} else {
		max = float64(0)
	}
	for _, element := range uniqueId.currentResources {
		var properties = element["Properties"].(map[string]interface{})
		for k, v := range properties {
			if k == "counter" && v.(float64) > max {
				max = v.(float64)
			}
		}
	}
	return max + 1
}

func (uniqueId *UniqueId) Invoke() (map[string]interface{}, error) {
	if uniqueId.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	var nextFreeCounter = uniqueId.getNextFreeCounter()
	value, ok := uniqueId.resourcePoolProperties["idFormat"]
	if !ok {
		return nil, errors.New("Missing idFormat in resources")
	}
	if !strings.Contains(value.(string), "{counter}") {
		return nil, errors.New("Missing {counter} in idFormat")
	}
	replacePoolProperties := make(map[string]interface{})
	for k, v := range uniqueId.resourcePoolProperties {
		if k != "idFormat" {
			replacePoolProperties[k] = v
		}
	}
	replacePoolProperties["counter"] = nextFreeCounter
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

		value = strings.Replace(value.(string), "{" + k + "}", v.(string), 1)
	}
	var result = make(map[string]interface{})
	result["text"] = value
	result["counter"] = nextFreeCounter
	return result, nil
}

func (uniqueId *UniqueId) Capacity() (map[string]interface{}, error) {
	var allocatedCapacity = float64(len(uniqueId.currentResources))
	var freeCapacity = float64(^uint(0) >> 1) - allocatedCapacity
	var result = make(map[string]interface{})
	result["freeCapacity"] = freeCapacity
	result["utilizedCapacity"] = allocatedCapacity
	return result, nil
}
