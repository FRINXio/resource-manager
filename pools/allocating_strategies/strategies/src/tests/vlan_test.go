package tests

import (
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func vlan(vlan float64) map[string]interface{} {
	vlanProperties := make(map[string]interface{})
	vlanMap := make(map[string]interface{})
	vlanMap["vlan"] = vlan
	vlanProperties["Properties"] = vlanMap
	return vlanProperties
}

func vlans(from int, to int) []map[string]interface{} {
	var vlansArray []map[string]interface{}
	for i := from; i <= to; i++ {
		vlansArray = append(vlansArray, vlan(float64(i)))
	}
	return vlansArray
}

func vlanRange(from int, to int) map[string]interface{} {
	rangeProperties := make(map[string]interface{})
	rangeMap := make(map[string]interface{})
	rangeMap["from"] = from
	rangeMap["to"] = to
	rangeProperties["Properties"] = rangeMap
	return rangeProperties
}

func TestMissingParentRange(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{}
	vlanStruct := src.NewVlan(allocated, resourcePool, nil)

	output, err := vlanStruct.Invoke()

	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Missing from in parentRange")
	eq = reflect.DeepEqual(err.Error(), expectedOutput.Error())
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestAllocateVlan(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": 0, "to": 4095}
	vlanStruct := src.NewVlan(allocated, resourcePool, nil)
	output, err := vlanStruct.Invoke()
	expectedOutput := map[string]interface{}{"vlan": float64(0)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	allocated = []map[string]interface{}{vlan(1)}
	resourcePool = map[string]interface{}{"from": 0, "to": 4095}
	vlanStruct = src.NewVlan(allocated, resourcePool, nil)
	output, err = vlanStruct.Invoke()
	expectedOutput = map[string]interface{}{"vlan": float64(0)}
	eq = reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}

	allocated = []map[string]interface{}{vlan(278)}
	resourcePool = map[string]interface{}{"from": 278, "to": 333}
	vlanStruct = src.NewVlan(allocated, resourcePool, nil)
	output, err = vlanStruct.Invoke()
	expectedOutput = map[string]interface{}{"vlan": float64(279)}
	eq = reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}

	resourcePool = map[string]interface{}{"from": 0, "to": 4095}
	vlanStruct = src.NewVlan(vlans(0, 4094), resourcePool, nil)
	output, err = vlanStruct.Invoke()
	expectedOutput = map[string]interface{}{"vlan": float64(4095)}
	eq = reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestAllocateVlanFull(t *testing.T) {
	var resourcePool = map[string]interface{}{"from": 0, "to": 4095}
	vlanStruct := src.NewVlan(vlans(0, 4095), resourcePool, nil)
	output, err := vlanStruct.Invoke()
	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Unable to allocate VLAN. Insufficient capacity to allocate a new vlan")
	eq = reflect.DeepEqual(err.Error(), expectedOutput.Error())
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestVlanCapacityMeasureEmpty(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": 0, "to": 4095}
	vlanStruct := src.NewVlan(allocated, resourcePool, nil)
	output, _ := vlanStruct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": float64(4096), "utilizedCapacity": float64(0)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestVlanCapacityMeasureFull(t *testing.T) {
	var resourcePool = map[string]interface{}{"from": 0, "to": 4095}
	vlanStruct := src.NewVlan(vlans(0, 4095), resourcePool, nil)
	output, _ := vlanStruct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": float64(0), "utilizedCapacity": float64(4096)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestFreeCapacity(t *testing.T) {
	vlanRangeProperties, _ := vlanRange(100, 900)["Properties"]
	output := src.FreeCapacity(vlanRangeProperties.(map[string]interface{}), float64(100))
	expectedOutput := float64(701)
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %f expected, got: %f", expectedOutput, output)
	}
}

func TestUtilisation(t *testing.T) {
	allocated := []map[string]interface{}{
		vlan(0), vlan(1), vlan(1000)}
	output := src.UtilizedCapacity(allocated, 1)
	expectedOutput := float64(4)
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %f expected, got: %f", expectedOutput, output)
	}
}

func TestVlanDesireValue(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"from": 0, "to": 4095}
	var userInput = map[string]interface{}{"desiredValue": "1111"}
	vlanStruct := src.NewVlan(allocated, resourcePool, userInput)
	output, err := vlanStruct.Invoke()
	expectedOutput := map[string]interface{}{"vlan": float64(1111)}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	// DesiredValue out of address/prefix
	userInput = map[string]interface{}{"desiredValue": "4096"}
	vlanStruct = src.NewVlan(allocated, resourcePool, userInput)
	output, err = vlanStruct.Invoke()
	expectedOutputError := errors.New("VLAN 4096 is out of range: 0 - 4095")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// Already claimed desiredValue
	allocated = []map[string]interface{}{vlan(278)}
	userInput = map[string]interface{}{"desiredValue": "278"}
	vlanStruct = src.NewVlan(allocated, resourcePool, userInput)
	output, err = vlanStruct.Invoke()
	expectedOutputError = errors.New("VLAN 278 was already claimed.")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// Wrong desiredValue input format
	allocated = []map[string]interface{}{}
	userInput = map[string]interface{}{"desiredValue": "Hello World"}
	vlanStruct = src.NewVlan(allocated, resourcePool, userInput)
	output, err = vlanStruct.Invoke()
	expectedOutputError = errors.New("VLAN must be a number. Received: Hello World.")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}
