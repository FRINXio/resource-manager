package src

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
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

func (uniqueId *UniqueId) getNextFreeCounterAndResourcesSet() (float64, map[float64]bool) {
	var max float64
	currentResourcesSet := make(map[float64]bool)
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
		var counter = properties["counter"].(float64)
		if counter > max {
			max = counter
		}
		currentResourcesSet[counter] = true
	}
	return max + 1, currentResourcesSet
}

func (uniqueId *UniqueId) Invoke() (map[string]interface{}, error) {
	if uniqueId.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	nextFreeCounter, currentResourcesSet := uniqueId.getNextFreeCounterAndResourcesSet()
	idFormat, ok := uniqueId.resourcePoolProperties["idFormat"]
	if !ok {
		return nil, errors.New("Missing idFormat in resources")
	}
	if !strings.Contains(idFormat.(string), "{counter}") {
		return nil, errors.New("Missing {counter} in idFormat")
	}

	if toValue, ok := uniqueId.resourcePoolProperties["to"]; ok {
		toValue, err := NumberToInt(toValue)
		if err != nil || nextFreeCounter > float64(toValue.(int)) {
			return nil, errors.New("Unable to allocate Unique-id from idFormat: \"" + idFormat.(string) + "\"." +
				" Insufficient capacity to allocate a unique-id.")
		}
	}

	replacePoolProperties := make(map[string]interface{})
	for k, v := range uniqueId.resourcePoolProperties {
		if k != "idFormat" && k != "counterFormatWidth" && k!= "from" && k!= "to" {
			replacePoolProperties[k] = v
		}
	}
	if prefixNumber, ok := uniqueId.resourcePoolProperties["counterFormatWidth"]; ok {
		replacePoolProperties["counter"] = fmt.Sprintf(
			"%0" + strconv.Itoa(prefixNumber.(int)) +"d", int(nextFreeCounter))
	} else {
		replacePoolProperties["counter"] = nextFreeCounter
	}

	if value, ok := uniqueId.userInput["desiredValue"]; ok {
		var result = make(map[string]interface{})
		pattern := idFormat.(string)
		for k, v := range replacePoolProperties {
			if k != "counter" {
				pattern = strings.Replace(pattern, "{" + k + "}", v.(string), 1)
			} else {
				pattern = strings.Replace(pattern, "{counter}", "(?P<counter>\\d+)", 1)
			}
		}
		rgx := regexp.MustCompile(pattern)
		stringSubmatchs := rgx.FindStringSubmatch(value.(string))

		if len(stringSubmatchs) > 0 {
			for index, subxpname := range rgx.SubexpNames() {
				if subxpname == "counter" {
					counter, _ := strconv.ParseFloat(stringSubmatchs[index], 0)
					if currentResourcesSet[counter] {
						return nil, errors.New("Unique-id " + value.(string) + " was already claimed." )
					}
					result["counter"] = counter
					result["text"] = value
					return result, nil
				}
			}
		}
		return nil, errors.New("Wrong desired value, does not equal idFormat: " + idFormat.(string))
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
