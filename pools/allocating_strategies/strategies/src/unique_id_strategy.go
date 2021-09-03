package src

import (
	"strconv"
	"strings"
)

// framework managed constants
type UniqueId struct {
	currentResources []map[string]interface{}
	resourcePoolProperties map[string]interface{}
}
// framework managed constants

func NewUniqueId(currentResources []map[string]interface{}, resourcePoolProperties map[string]interface{}) UniqueId {
	return UniqueId{currentResources, resourcePoolProperties}
}

// STRATEGY_START

func (uniqueId UniqueId) getNextFreeCounter() int {
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
			if k == "counter" && v.(int) > max {
				max = v.(int)
			}
		}
	}
	return max + 1
}

func (uniqueId UniqueId) invoke() map[string]interface{} {
	if uniqueId.resourcePoolProperties == nil {
		// console.error("Unable to extract resources")
		return nil
	}
	var nextFreeCounter = uniqueId.getNextFreeCounter()
	value, ok := uniqueId.resourcePoolProperties["idFormat"]
	if !ok {
		// console.error("Missing idFormat in resources")
		return nil
	}
	if !strings.Contains(value.(string), "{counter}") {
		//console.error("Missing {counter} in idFormat")
		return nil
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
	return result
}

func (uniqueId UniqueId) capacity() map[string]interface{} {
	var allocatedCapacity = uniqueId.getNextFreeCounter() - 1
	var freeCapacity = int(^uint(0) >> 1) - allocatedCapacity
	var result = make(map[string]interface{})
	result["freeCapacity"] = freeCapacity
	result["utilizedCapacity"] = allocatedCapacity
	return result
}

// STRATEGY_END