package tests

import (
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func TestMissingParentRangeForVlanRange(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{}
	var userInput = map[string]interface{}{}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
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

func TestMissingUserInputForVlanRange(t *testing.T) {
	var allocated []map[string]interface{}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedOutput := errors.New("Missing desiredSize in userInput")
	eq = reflect.DeepEqual(err.Error(), expectedOutput.Error())
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, err)
	}
}

func TestAllocate4096(t *testing.T) {
	var allocated []map[string]interface{}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{"desiredSize": 4096}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": 0, "to": 4095}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestAllocate4097(t *testing.T) {
	var allocated []map[string]interface{}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{"desiredSize": 4097}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedErrorOutput := errors.New("Unable to allocate VLAN range: 0 - 4095." +
		" Insufficient capacity to allocate a new range of size: 4097.")
	eq = reflect.DeepEqual(err.Error(), expectedErrorOutput.Error())
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedErrorOutput.Error(), err)
	}
}

func TestAllocateNoCapacity(t *testing.T) {
	allocated := []map[string]interface{}{vlanRange(0, 2000), vlanRange(2001, 4090)}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{"desiredSize": 100}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	eq := reflect.DeepEqual(output, (map[string]interface{})(nil))
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	expectedErrorOutput := errors.New("Unable to allocate VLAN range: 0 - 4095." +
		" Insufficient capacity to allocate a new range of size: 100.")
	eq = reflect.DeepEqual(err.Error(), expectedErrorOutput.Error())
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedErrorOutput.Error(), err)
	}
}

func TestAllocateRange1(t *testing.T) {
	var allocated []map[string]interface{}
	resourcePool := map[string]interface{}{"from": 0, "to": 33}
	userInput := map[string]interface{}{"desiredSize": 1}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": 0, "to": 0}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestAllocateRange784(t *testing.T) {
	var allocated []map[string]interface{}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{"desiredSize": 784}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": 0, "to": 783}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestAllocateReleasedRange(t *testing.T) {
	allocated := []map[string]interface{}{vlanRange(0, 100), vlanRange(200, 300)}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{"desiredSize": 10}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": 101, "to": 110}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	userInput = map[string]interface{}{"desiredSize": 1000}
	vlanRangeStruct = src.NewVlanRange(allocated, resourcePool, userInput)
	output, err = vlanRangeStruct.Invoke()
	expectedOutput = map[string]interface{}{"from": 301, "to": 1300}
	eq = reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	allocated = []map[string]interface{}{vlanRange(100, 200)}
	userInput = map[string]interface{}{"desiredSize": 10}
	vlanRangeStruct = src.NewVlanRange(allocated, resourcePool, userInput)
	output, err = vlanRangeStruct.Invoke()
	expectedOutput = map[string]interface{}{"from": 0, "to": 9}
	eq = reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestAllocateRangeAtTheEnd(t *testing.T) {
	// allocate range of 1 at the end of parent range
	allocated := []map[string]interface{}{
		vlanRange(0, 1000), vlanRange(1001, 3000), vlanRange(3001, 4090)}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{"desiredSize": 1}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": 4091, "to": 4091}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	// allocate range of 4 at the end of parent range, totally exhausting the range
	allocated = []map[string]interface{}{vlanRange(0, 1000), vlanRange(1001, 3000),
		vlanRange(3001, 4090), vlanRange(4091, 4091)}
	userInput = map[string]interface{}{"desiredSize": 4}
	vlanRangeStruct = src.NewVlanRange(allocated, resourcePool, userInput)
	output, err = vlanRangeStruct.Invoke()
	expectedOutput = map[string]interface{}{"from": 4092, "to": 4095}
	eq = reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestVlanRangeFreeCapacity(t *testing.T) {
	allocated, _ := vlanRange(100, 4055)["Properties"]
	output, err := src.FreeRangeCapacity(allocated.(map[string]interface{}), 100)
	expectedOutput := 3856
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestVlanRangeCapacityMeasure(t *testing.T) {
	allocated := []map[string]interface{}{vlanRange(0, 2000), vlanRange(2001, 4090)}
	resourcePool := map[string]interface{}{"from": 0, "to": 4095}
	userInput := map[string]interface{}{"desiredSize": 100}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, _ := vlanRangeStruct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": float64(5), "utilizedCapacity": float64(4091)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestRangeUtilisation(t *testing.T) {
	allocated := []map[string]interface{}{vlanRange(0, 1000), vlanRange(1001, 3000), vlanRange(3001, 4090),
		vlanRange(4091, 4091)}
	output := src.UtilizedRangeCapacity(allocated, 1)
	expectedOutput := 4093
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
}

func TestRangeFreeCapacity(t *testing.T) {
	allocated := vlanRange(100, 4055)["Properties"]
	output, err := src.FreeRangeCapacity(allocated.(map[string]interface{}), 100)
	expectedOutput := 3856
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestFromCurrentResourcesToVlanRangeValues(t *testing.T) {
	allocated := []map[string]interface{}{vlanRange(0, 1000), vlanRange(1001, 3000), vlanRange(3001, 4090),
		vlanRange(4091, 4091)}
	output := src.FromCurrentResourcesToVlanRangeValues(allocated)
	expectedOutput := src.NewVlanRangeValues(0, 1000)
	eq := reflect.DeepEqual(expectedOutput, output[0])
	if !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
	expectedOutput = src.NewVlanRangeValues(1001, 3000)
	eq = reflect.DeepEqual(expectedOutput, output[1])
	if !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
	expectedOutput = src.NewVlanRangeValues(3001, 4090)
	eq = reflect.DeepEqual(expectedOutput, output[2])
	if !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
	expectedOutput = src.NewVlanRangeValues(4091, 4091)
	eq = reflect.DeepEqual(expectedOutput, output[3])
	if !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
}
