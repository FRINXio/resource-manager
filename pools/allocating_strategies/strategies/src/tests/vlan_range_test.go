package tests

import (
	"fmt"
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

func vlanFloatRange(from float64, to float64) map[string]interface{} {
	rangeProperties := make(map[string]interface{})
	rangeMap := make(map[string]interface{})
	rangeMap["from"] = from
	rangeMap["to"] = to
	rangeProperties["Properties"] = rangeMap
	return rangeProperties
}

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
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
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
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
	userInput := map[string]interface{}{"desiredSize": float64(4096)}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": float64(0), "to": float64(4095)}
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
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
	userInput := map[string]interface{}{"desiredSize": float64(4097)}
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
	allocated := []map[string]interface{}{vlanFloatRange(float64(0), float64(2000)),
		vlanFloatRange(float64(2001), float64(4090))}
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
	userInput := map[string]interface{}{"desiredSize": float64(100)}
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
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(33)}
	userInput := map[string]interface{}{"desiredSize": float64(1)}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": float64(0), "to": float64(0)}
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
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
	userInput := map[string]interface{}{"desiredSize": float64(784)}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": float64(0), "to": float64(783)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestVlanRangeCapacityMeasure(t *testing.T) {
	allocated := []map[string]interface{}{vlanFloatRange(float64(0), float64(2000)),
		vlanFloatRange(float64(2001), float64(4090))}
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
	userInput := map[string]interface{}{"desiredSize": float64(100)}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, _ := vlanRangeStruct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": float64(5), "utilizedCapacity": float64(4091)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestAllocateReleasedRange(t *testing.T) {
	allocated := []map[string]interface{}{vlanFloatRange(float64(0), float64(100)),
		vlanFloatRange(float64(200), float64(300))}
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
	userInput := map[string]interface{}{"desiredSize": float64(10)}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": float64(101), "to": float64(110)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	userInput = map[string]interface{}{"desiredSize": float64(1000)}
	vlanRangeStruct = src.NewVlanRange(allocated, resourcePool, userInput)
	output, err = vlanRangeStruct.Invoke()
	expectedOutput = map[string]interface{}{"from": float64(301), "to": float64(1300)}
	eq = reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	allocated = []map[string]interface{}{vlanFloatRange(float64(100), float64(200))}
	userInput = map[string]interface{}{"desiredSize": float64(10)}
	vlanRangeStruct = src.NewVlanRange(allocated, resourcePool, userInput)
	output, err = vlanRangeStruct.Invoke()
	expectedOutput = map[string]interface{}{"from": float64(0), "to": float64(9)}
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
	allocated := []map[string]interface{}{vlanFloatRange(float64(0), float64(1000)),
		vlanFloatRange(float64(1001), float64(3000)), vlanFloatRange(float64(3001), float64(4090))}
	resourcePool := map[string]interface{}{"from": float64(0), "to": float64(4095)}
	userInput := map[string]interface{}{"desiredSize": float64(1)}
	vlanRangeStruct := src.NewVlanRange(allocated, resourcePool, userInput)
	output, err := vlanRangeStruct.Invoke()
	expectedOutput := map[string]interface{}{"from": float64(4091), "to": float64(4091)}
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	eq = reflect.DeepEqual(err, nil)
	if !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	// allocate range of 4 at the end of parent range, totally exhausting the range
	allocated = []map[string]interface{}{
		vlanFloatRange(float64(0), float64(1000)), vlanFloatRange(float64(1001), float64(3000)),
		vlanFloatRange(float64(3001), float64(4090)), vlanFloatRange(float64(4091), float64(4091))}
	userInput = map[string]interface{}{"desiredSize": float64(4)}
	vlanRangeStruct = src.NewVlanRange(allocated, resourcePool, userInput)
	output, err = vlanRangeStruct.Invoke()
	expectedOutput = map[string]interface{}{"from": float64(4092), "to": float64(4095)}
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
	allocated, _ := vlanFloatRange(float64(100), float64(4055))["Properties"]
	fmt.Println(allocated)
	output := src.FreeCapacityFloat(allocated.(map[string]interface{}), float64(100))
	expectedOutput := float64(3856)
	eq := reflect.DeepEqual(output, expectedOutput)
	if !eq {
		t.Fatalf("different output of %f expected, got: %f", expectedOutput, output)
	}
}
