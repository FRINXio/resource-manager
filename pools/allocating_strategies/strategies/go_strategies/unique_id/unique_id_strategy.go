package unique_id

import (
	"strconv"
	"strings"
)

var currentResources []map[string]interface{}
var resourcePoolProperties map[string]interface{}
var userInput map[string]interface{}


func getNextFreeCounter(properties map[string]interface{}) int {
	var max int
	value, ok := properties["from"]
	if ok {
		max = value.(int) - 1
	} else {
		max = 0
	}
	for _, element := range currentResources {
		var properties = element["Properties"].(map[string]interface{})
		for k, v := range properties {
			if k == "counter" && v.(int) > max {
				max = v.(int)
			}
		}
	}
	return max + 1
}

func invoke() map[string]interface{} {
	if resourcePoolProperties == nil {
		// console.error("Unable to extract resources")
		return nil
	}
	var nextFreeCounter = getNextFreeCounter(resourcePoolProperties)
	value, ok := resourcePoolProperties["idFormat"]
	if !ok {
		// console.error("Missing idFormat in resources")
		return nil
	}
	if !strings.Contains(value.(string), "{counter}") {
		//console.error("Missing {counter} in idFormat")
		return nil
	}
	replacePoolProperties := make(map[string]interface{})
	for k, v := range resourcePoolProperties {
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

func capacity() map[string]interface{} {
	var allocatedCapacity = getNextFreeCounter(resourcePoolProperties) - 1
	var freeCapacity = int(^uint(0) >> 1) - allocatedCapacity
	var result = make(map[string]interface{})
	result["freeCapacity"] = freeCapacity
	result["utilizedCapacity"] = allocatedCapacity
	return result
}

// STRATEGY_END

// For testing purposes
func invokeWithParams(
	currentResourcesArg []map[string]interface{},
	resourcePoolArg map[string]interface{},
	userInputArg map[string]interface{}) map[string]interface{} {

	currentResources = currentResourcesArg
	resourcePoolProperties = resourcePoolArg
	userInput = userInputArg
	return invoke()
}

func invokeWithParamsCapacity(
	currentResourcesArg []map[string]interface{},
	resourcePoolArg map[string]interface{},
	userInputArg map[string]interface{}) map[string]interface{} {

	currentResources = currentResourcesArg
	resourcePoolProperties = resourcePoolArg
	userInput = userInputArg
	return capacity()
}