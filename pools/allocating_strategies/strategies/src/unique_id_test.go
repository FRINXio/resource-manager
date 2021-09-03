package src

import (
	"reflect"
	"strconv"
	"testing"
)

func uniqueId(counter int, text string) map[string]interface{}{
	uniqueIdProperties := make(map[string]interface{})
	uniqueIdMap := make(map[string]interface{})
	uniqueIdMap["counter"] = counter
	uniqueIdMap["text"] = text
	uniqueIdProperties["Properties"] = uniqueIdMap
	return uniqueIdProperties
}

func TestUniqueIdOutputAndCapacity(t *testing.T) {
	var allocated = []map[string]interface{}{
		uniqueId(0, "first"),
		uniqueId(1, "second"),
		uniqueId(3, "third"),
	}

	var resourcePool = map[string]interface{}{
		"vpn": "VPN85", "network": "Network19",
		"idFormat": "VPN-{counter}-{network}-{vpn}-local",
	}
	uniqueIdStruct := NewUniqueId(allocated, resourcePool)

	var output = uniqueIdStruct.invoke()
	var expectedOutput = make(map[string]interface{})
	expectedOutput["counter"] = 4
	expectedOutput["text"] = "VPN-4-Network19-VPN85-local"

	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}

	uniqueIdStruct = NewUniqueId(allocated, resourcePool)
	var capacity = uniqueIdStruct.capacity()
	var expectedCapacity = make(map[string]interface{})
	expectedCapacity["freeCapacity"] = int(^uint(0) >> 1) - 3  // Number.MAX_SAFE_INTEGER - 3
	expectedCapacity["utilizedCapacity"] = 3

	eq = reflect.DeepEqual(capacity, expectedCapacity)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedCapacity, capacity)
	}
}

func TestSimpleRangeCounter(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": 1000, "idFormat": "{counter}"}
	uniqueIdStruct := NewUniqueId(outputs, resourcePool)
	for i := 0; i <= 10; i++ {
		var output = uniqueIdStruct.invoke()
		var expectedOutput = map[string]interface{}{"counter": 1000 + i, "text": strconv.Itoa(1000 + i)}
		eq := reflect.DeepEqual(output, expectedOutput)
		if !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		outputs = append(outputs, uniqueId(output["counter"].(int), output["text"].(string)))
		uniqueIdStruct = NewUniqueId(outputs, resourcePool)
	}

}

func TestParamsWithoutResourcePool(t *testing.T) {
	var allocated []map[string]interface{}

	uniqueIdStruct := NewUniqueId(allocated, nil)
	var output = uniqueIdStruct.invoke()

	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestResourcePoolWithoutIdFormat(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{
		"vpn": "VPN85", "network": "Network19",
		"idFormat": "VPN-{network}-{vpn}-local",
	}
	uniqueIdStruct := NewUniqueId(allocated, resourcePool)
	var output = uniqueIdStruct.invoke()

	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestMultipleL3vpnCounters(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"someProperty": "L3VPN", "idFormat": "{someProperty}{counter}"}
	uniqueIdStruct := NewUniqueId(outputs, resourcePool)
	for i := 1; i <= 10; i++ {
		var output = uniqueIdStruct.invoke()
		var expectedOutput = map[string]interface{}{"counter": i, "text": "L3VPN" + strconv.Itoa(i)}
		eq := reflect.DeepEqual(output, expectedOutput)
		if !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		outputs = append(outputs, uniqueId(output["counter"].(int), output["text"].(string)))
		uniqueIdStruct = NewUniqueId(outputs, resourcePool)

		var capacity = uniqueIdStruct.capacity()
		var expectedCapacity = make(map[string]interface{})
		expectedCapacity["freeCapacity"] = int(^uint(0) >> 1) - i  // Number.MAX_SAFE_INTEGER - 3
		expectedCapacity["utilizedCapacity"] = i
		eq = reflect.DeepEqual(capacity, expectedCapacity)
		if !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedCapacity, capacity)
		}
	}
}