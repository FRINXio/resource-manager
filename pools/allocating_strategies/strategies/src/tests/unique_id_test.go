package tests

import (
	"fmt"
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
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(allocated, resourcePool, userInput)

	var output, err = uniqueIdStruct.Invoke()
	var expectedOutput = make(map[string]interface{})
	expectedOutput["counter"] = float64(2)
	expectedOutput["text"] = "VPN-2-Network19-VPN85-local"

	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	uniqueIdStruct = src.NewUniqueId(allocated, resourcePool, userInput)
	capacity, err := uniqueIdStruct.Capacity()
	var expectedCapacity = make(map[string]interface{})
	expectedCapacity["freeCapacity"] = float64(^uint(0) >> 1) - 3  // Number.MAX_SAFE_INTEGER - 3
	expectedCapacity["utilizedCapacity"] = float64(3)

	if eq := reflect.DeepEqual(capacity, expectedCapacity); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedCapacity, capacity)
	}
}

func TestSimpleRangeCounter(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": 1000, "idFormat": "{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	for i := 0; i <= 10; i++ {
		var output, err = uniqueIdStruct.Invoke()
		var expectedOutput = map[string]interface{}{"counter": float64(1000 + i), "text": fmt.Sprint(1000 + i)}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
	}
}

func TestParamsWithoutResourcePool(t *testing.T) {
	var allocated []map[string]interface{}
	var userInput = map[string]interface{}{}

	uniqueIdStruct := src.NewUniqueId(allocated, nil, userInput)
	var output, err = uniqueIdStruct.Invoke()

	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Unable to extract resources")
	if eq := reflect.DeepEqual(err.Error(), expectedOutput.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestResourcePoolWithoutIdFormat(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{
		"vpn": "VPN85", "network": "Network19",
		"idFormat": "VPN-{network}-{vpn}-local",
	}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(allocated, resourcePool, userInput)
	output, err := uniqueIdStruct.Invoke()

	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Missing {counter} in idFormat")
	if eq := reflect.DeepEqual(err.Error(), expectedOutput.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestMultipleL3vpnCounters(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"someProperty": "L3VPN", "idFormat": "{someProperty}{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	for i := 1; i <= 10; i++ {
		output, err := uniqueIdStruct.Invoke()
		expectedOutput := map[string]interface{}{"counter": float64(i), "text": "L3VPN" + strconv.Itoa(i)}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)

		capacity, err := uniqueIdStruct.Capacity()
		var expectedCapacity = make(map[string]interface{})
		expectedCapacity["freeCapacity"] = float64(^uint(0) >> 1) - float64(i)  // Number.MAX_SAFE_INTEGER - 3
		expectedCapacity["utilizedCapacity"] = float64(i)
		if eq := reflect.DeepEqual(capacity, expectedCapacity); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedCapacity, capacity)
		}
	}
}

func TestSimplePrefixNumber(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"counterFormatWidth": 5, "idFormat": "{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	for i := 1; i <= 10; i++ {
		var output, err = uniqueIdStruct.Invoke()
		var expectedOutput = map[string]interface{}{"counter": float64(i), "text": fmt.Sprintf("%05d", i)}
		eq := reflect.DeepEqual(output, expectedOutput)
		if !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq = reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
	}
}

func TestCapacityTo(t *testing.T) {
	to := 25
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"to": to, "idFormat": "{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	for i := 0; i <= 10; i++ {
		output, err := uniqueIdStruct.Invoke()
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
		capacityOutput, err := uniqueIdStruct.Capacity()
		var expectedOutput = map[string]interface{}{"freeCapacity": float64(to - i), "utilizedCapacity": float64(i + 1) }
		if eq := reflect.DeepEqual(capacityOutput, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, capacityOutput)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
	}
}

func TestInvokeFromTo(t *testing.T) {
	from := 500
	to := 510
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": from, "to": to, "idFormat": "{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	for i := 0; i <= 10; i++ {
		var output, err = uniqueIdStruct.Invoke()
		var expectedOutput = map[string]interface{}{"counter": float64(from + i), "text": strconv.Itoa(from + i)}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
	}
	output, err := uniqueIdStruct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Unable to allocate Unique-id from idFormat: \"{counter}\"." +
		" Insufficient capacity to allocate a unique-id.")
	if eq := reflect.DeepEqual(err.Error(), expectedOutput.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestCapacityFromTo(t *testing.T) {
	from := 500
	to := 510
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": from, "to": to, "idFormat": "{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	for i := 0; i <= 10; i++ {
		output, err := uniqueIdStruct.Invoke()
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
		capacityOutput, err := uniqueIdStruct.Capacity()
		var expectedOutput = map[string]interface{}{"freeCapacity": float64(to - from - i), "utilizedCapacity": float64(i + 1) }
		if eq := reflect.DeepEqual(capacityOutput, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, capacityOutput)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
	}
}

func TestDesiredValueSimple(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"idFormat": "{counter}"}
	var userInput = map[string]interface{}{"desiredValue": 5}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	output, err := uniqueIdStruct.Invoke()
	var expectedOutput = map[string]interface{}{"counter": float64(5), "text": "5"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestDesiredValueSimple2(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"idFormat": "L3VPN{counter}"}
	var userInput = map[string]interface{}{"desiredValue": 5}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	output, err := uniqueIdStruct.Invoke()
	var expectedOutput = map[string]interface{}{"counter": float64(5), "text": "L3VPN" + strconv.Itoa(5)}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestDesiredValueSimple3(t *testing.T) {
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"idFormat": "L3VPN{counter}{pattern}", "pattern": "AAA"}
	var userInput = map[string]interface{}{"desiredValue": 5}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	output, err := uniqueIdStruct.Invoke()
	var expectedOutput = map[string]interface{}{"counter": float64(5), "text": "L3VPN" + strconv.Itoa(5) + "AAA"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestDesiredValueAlreadyClaimedValue(t *testing.T) {
	var outputs = []map[string]interface{}{uniqueId(float64(5), "L3VPN5AAA")}
	var resourcePool = map[string]interface{}{"idFormat": "L3VPN{counter}{pattern}", "pattern": "AAA"}
	var userInput = map[string]interface{}{"desiredValue": 5}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	output, err := uniqueIdStruct.Invoke()

	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Unique-id L3VPN5AAA was already claimed.")
	if eq := reflect.DeepEqual(err.Error(), expectedOutput.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestDesiredValueComplicated(t *testing.T) {
	from := 2
	to := 11
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"counterFormatWidth": 2, "from": from, "to": to, "someProperty": "L3VPN", "idFormat": "{someProperty}{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)
	for i := 0; i <= 8; i++ {
		output, err := uniqueIdStruct.Invoke()
		expectedOutput := map[string]interface{}{"counter": float64(i + from), "text": "L3VPN" + fmt.Sprintf("%02d", i + from)}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
	}
	userInput = map[string]interface{}{"desiredValue": 7}
	uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)

	output, err := uniqueIdStruct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	var expectedOutputError = errors.New("Unique-id L3VPN07 was already claimed.")
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	userInput = map[string]interface{}{}
	uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
	output, err = uniqueIdStruct.Invoke()
	expectedOutput := map[string]interface{}{"counter": float64(11), "text": "L3VPN11"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))

	uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
	output, err = uniqueIdStruct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutputError = errors.New("Unable to allocate Unique-id from idFormat: \"{someProperty}{counter}\". Insufficient capacity to allocate a unique-id.")
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	userInput = map[string]interface{}{"desiredValue": 25}
	uniqueIdStruct = src.NewUniqueId(outputs, resourcePool, userInput)
	output, err = uniqueIdStruct.Invoke()

	expectedOutputError = errors.New("Unable to allocate Unique-id desiredValue: 25. Value is out of scope: 11")
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}

func TestFirstMissing(t *testing.T) {
	from := 0
	to := 8
	var outputs = []map[string]interface{}{uniqueId(float64(7), "L3VPN-7"), uniqueId(float64(1), "L3VPN-1"),
		uniqueId(float64(3), "L3VPN-3"), uniqueId(float64(4), "L3VPN-4")}
	var resourcePool = map[string]interface{}{ "from": from, "to": to, "idFormat": "L3VPN-{counter}"}
	var userInput = map[string]interface{}{}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool, userInput)

	output, err := uniqueIdStruct.Invoke()
	expectedOutput := map[string]interface{}{"counter": float64(0), "text": "L3VPN-0"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}