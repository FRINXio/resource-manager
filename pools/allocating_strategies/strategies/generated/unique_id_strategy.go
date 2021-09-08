package src

// STRATEGY_START

import (
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

func (uniqueId *UniqueId) getNextFreeCounter() int {
	var max int
	value, ok := uniqueId.resourcePoolProperties["from"]
	if ok {
		max = value.(int) - 1
	} else {
		max = 0
	}
	for _, element := range uniqueId.currentResources {
		var properties = element["Properties"].(map[string]interface{})
		for k, v := range properties {
			if k == "counter" && int(v.(float64)) > max {
				max = int(v.(float64))
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
		case int:
			v = strconv.Itoa(v.(int))
		}
		value = strings.Replace(value.(string), "{" + k + "}", v.(string), 1)
	}
	var result = make(map[string]interface{})
	result["text"] = value
	result["counter"] = nextFreeCounter
	return result, nil
}

func (uniqueId *UniqueId) Capacity() map[string]interface{} {
	var allocatedCapacity = uniqueId.getNextFreeCounter() - 1
	var freeCapacity = int(^uint(0) >> 1) - allocatedCapacity
	var result = make(map[string]interface{})
	result["freeCapacity"] = freeCapacity
	result["utilizedCapacity"] = allocatedCapacity
	return result
}

// STRATEGY_END