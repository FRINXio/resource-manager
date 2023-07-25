package tests

import (
	"fmt"
	strategies "github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"testing"
)

func TestUniqueIdCapacity(t *testing.T) {
	allocated := []map[string]interface{}{
		{"counter": 0, "text": "id-0"},
		{"counter": 1, "text": "id-1"},
		{"counter": 2, "text": "id-2"},
	}
	userInput := map[string]interface{}{
		"desiredValue": 10,
	}
	allocated = append(allocated, map[string]interface{}{"counter": 0, "text": "id-0"})
	resourcePoolProperties := map[string]interface{}{
		"from":     0,
		"to":       10,
		"idFormat": "id-{counter}",
	}

	uniqueId := strategies.NewUniqueId(allocated, resourcePoolProperties, userInput)
	uniqueId.Invoke()

	capacity, err := uniqueId.Capacity()

	if err != nil {
		t.Errorf("UniqueId capacity test failed")
	}

	fmt.Println("UniqueId capacity: ", capacity)
}
