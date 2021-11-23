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
	uniqueIdStruct := src.NewUniqueId(allocated, resourcePool)

	var output, err = uniqueIdStruct.Invoke()
	var expectedOutput = make(map[string]interface{})
	expectedOutput["counter"] = float64(4)
	expectedOutput["text"] = "VPN-4-Network19-VPN85-local"

	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	uniqueIdStruct = src.NewUniqueId(allocated, resourcePool)
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
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
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
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)
	}
}

func TestParamsWithoutResourcePool(t *testing.T) {
	var allocated []map[string]interface{}

	uniqueIdStruct := src.NewUniqueId(allocated, nil)
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
	uniqueIdStruct := src.NewUniqueId(allocated, resourcePool)
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
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
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
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)

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
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
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
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)
	}
}

func TestCapacityTo(t *testing.T) {
	to := 25
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"to": to, "idFormat": "{counter}"}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
	for i := 0; i <= 10; i++ {
		output, err := uniqueIdStruct.Invoke()
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)
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
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
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
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)
	}
	output, err := uniqueIdStruct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Unable to allocate Unique-id from idFormat: \"{counter}\"." +
		"Insufficient capacity to allocate a unique-id.")
	if eq := reflect.DeepEqual(err.Error(), expectedOutput.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestCapacityFromTo(t *testing.T) {
	from := 500
	to := 510
	var outputs []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": from, "to": to, "idFormat": "{counter}"}
	uniqueIdStruct := src.NewUniqueId(outputs, resourcePool)
	for i := 0; i <= 10; i++ {
		output, err := uniqueIdStruct.Invoke()
		outputs = append(outputs, uniqueId(output["counter"].(float64), output["text"].(string)))
		uniqueIdStruct = src.NewUniqueId(outputs, resourcePool)
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