package tests

import (
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"testing"
)

func uniqueId(counter float64, text string) map[string]interface{}{
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
	uniqueIdStruct := src.NewUniqueId(allocated, resourcePool)

	var output, err = uniqueIdStruct.Invoke()
	var expectedOutput = make(map[string]interface{})
	expectedOutput["counter"] = 4
	expectedOutput["text"] = "VPN-4-Network19-VPN85-local"

	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	uniqueIdStruct = src.NewUniqueId(allocated, resourcePool)
	var capacity = uniqueIdStruct.Capacity()
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
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
	for i := 0; i <= 10; i++ {
		var output, err = uniqueIdStruct.Invoke()
		var expectedOutput = map[string]interface{}{"counter": 1000 + i, "text": strconv.Itoa(1000 + i)}
		eq := reflect.DeepEqual(output, expectedOutput)
		if !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		eq = reflect.DeepEqual(err, nil)
		if !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		outputs = append(outputs, uniqueId(float64(output["counter"].(int)), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)
	}

}

func TestParamsWithoutResourcePool(t *testing.T) {
	var allocated []map[string]interface{}

	uniqueIdStruct := src.NewUniqueId(allocated, nil)
	var output, err = uniqueIdStruct.Invoke()

	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Unable to extract resources")
	eq = reflect.DeepEqual(err.Error(), expectedOutput.Error())
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestResourcePoolWithoutIdFormat(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{
		"vpn": "VPN85", "network": "Network19",
		"idFormat": "VPN-{network}-{vpn}-local",
	}
	uniqueIdStruct := src.NewUniqueId(allocated, resourcePool)
	output, err := uniqueIdStruct.Invoke()

	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Missing {counter} in idFormat")
	eq = reflect.DeepEqual(err.Error(), expectedOutput.Error())
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestMultipleL3vpnCounters(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"someProperty": "L3VPN", "idFormat": "{someProperty}{counter}"}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
	for i := 1; i <= 10; i++ {
		output, err := uniqueIdStruct.Invoke()
		expectedOutput := map[string]interface{}{"counter": i, "text": "L3VPN" + strconv.Itoa(i)}
		eq := reflect.DeepEqual(output, expectedOutput)
		if !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		eq = reflect.DeepEqual(err, nil)
		if !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		outputs = append(outputs, uniqueId(float64(output["counter"].(int)), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)

		var capacity = uniqueIdStruct.Capacity()
		var expectedCapacity = make(map[string]interface{})
		expectedCapacity["freeCapacity"] = int(^uint(0) >> 1) - i  // Number.MAX_SAFE_INTEGER - 3
		expectedCapacity["utilizedCapacity"] = i
		eq = reflect.DeepEqual(capacity, expectedCapacity)
		if !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedCapacity, capacity)
		}
	}
}