package tests

import (
	"fmt"
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"strconv"
	"testing"
)

func ipv4Prefix(address string, prefix int, subnet bool) map[string]interface{} {
	ipv4Properties := make(map[string]interface{})
	ipv4PrefixMap := make(map[string]interface{})
	ipv4PrefixMap["address"] = address
	ipv4PrefixMap["prefix"] = prefix
	ipv4PrefixMap["subnet"] = subnet
	ipv4Properties["Properties"] = ipv4PrefixMap
	return ipv4Properties
}

func TestSingleAllocationPool(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	for i := 1; i <= 8; i++ {
		desiredSize := math.Pow(2, float64(i))
		var userInput = map[string]interface{}{"desiredSize": desiredSize}
		ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
		output, err := ipv4PrefixStruct.Invoke()
		expectedOutput := map[string]interface{}{"address": "192.168.1.0", "prefix": 32 - i, "subnet": false}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
	}
}

func TestSingleAllocationSubnet(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": true}
	for i := 1; i <= 7; i++ {
		desiredSize := math.Pow(2, float64(i))
		var userInput = map[string]interface{}{"desiredSize": desiredSize, "subnet": true}
		ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
		output, err := ipv4PrefixStruct.Invoke()
		expectedOutput := map[string]interface{}{"address": "192.168.1.0", "prefix": 32 - i, "subnet": true}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
	}
}

func TestAllocateRangeAtStartWithExistingResources(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4Prefix("192.168.1.16", 28, false)}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput = map[string]interface{}{"desiredSize": 10}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Invoke()
	expectedOutput := map[string]interface{}{"address": "192.168.1.0", "prefix": 28, "subnet": false}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestIpv4prefixCapacity24Mask(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4Prefix("192.168.1.16", 28, false)}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput map[string]interface{}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": strconv.Itoa(242), "utilizedCapacity": strconv.Itoa(14)}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestIpv4PrefixAllocationSubnetVsPool(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": true}
	var userInput = map[string]interface{}{"desiredSize": 2, "subnet": true}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Invoke()
	expectedOutput := map[string]interface{}{"address": "192.168.1.0", "prefix": 31, "subnet": true}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	userInput = map[string]interface{}{"desiredSize": 2, "subnet": true}
	ipv4PrefixStruct = src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err = ipv4PrefixStruct.Invoke()
	expectedOutput = map[string]interface{}{"address": "192.168.1.0", "prefix": 31, "subnet": true}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	userInput = map[string]interface{}{"desiredSize": 256}
	ipv4PrefixStruct = src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err = ipv4PrefixStruct.Invoke()
	expectedOutput = map[string]interface{}{"address": "192.168.1.0", "prefix": 24, "subnet": true}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	allocated = []map[string]interface{}{ipv4Prefix("192.168.1.0", 24, false)}
	userInput = map[string]interface{}{"desiredSize": 256}
	ipv4PrefixStruct = src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err = ipv4PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestIpv4PrefixAllocation24(t *testing.T) {
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var subnet []map[string]interface{}
	var expectedSubnets = []map[string]interface{}{
		ipv4Prefix("192.168.1.0", 28, false),   // 10 ->   0 -  15
		ipv4Prefix("192.168.1.32", 27, false),  // 19 ->  32 -  63
		ipv4Prefix("192.168.1.64", 26, false),  // 39 ->  64 - 127
		ipv4Prefix("192.168.1.16", 31, false),  //  2 ->  16 -  17
		ipv4Prefix("192.168.1.128", 28, false), // 14 -> 128 - 143
		ipv4Prefix("192.168.1.24", 29, false),  //  8 ->  24 -  31
		ipv4Prefix("192.168.1.192", 26, false), // 64 -> 196 - 255
	}
	counter := 0
	for _, i := range []int{10, 19, 39, 2, 14, 8, 64} {
		userInput := map[string]interface{}{"desiredSize": i}
		ipv4PrefixStruct := src.NewIpv4Prefix(subnet, resourcePool, userInput)
		output, err := ipv4PrefixStruct.Invoke()
		if eq := reflect.DeepEqual(output, expectedSubnets[counter]["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedSubnets[counter], output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		subnet = append(subnet, ipv4Prefix(output["address"].(string), output["prefix"].(int), false))
		counter++
	}

	// utilised capacity should be 16+32+64+2+16+8+64=202
	// free capacity should be 256-202=54
	// utilisation should be 202/256=78.9%
	// free blocks should be: 19-23, 144-195

	// Round 2, try to squeeze in additional subnets
	var expectedSubnets2 = []map[string]interface{}{
		ipv4Prefix("192.168.1.20", 30, false),
		ipv4Prefix("192.168.1.160", 27, false),
		ipv4Prefix("192.168.1.144", 29, false),
		ipv4Prefix("192.168.1.152", 30, false),
		ipv4Prefix("192.168.1.18", 31, false),
		ipv4Prefix("192.168.1.156", 30, false),
	}
	counter = 0
	for _, i := range []int{4, 32, 8, 4, 2, 4} {
		userInput := map[string]interface{}{"desiredSize": i}
		ipv4PrefixStruct := src.NewIpv4Prefix(subnet, resourcePool, userInput)
		output, err := ipv4PrefixStruct.Invoke()
		if eq := reflect.DeepEqual(output, expectedSubnets2[counter]["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedSubnets2[counter], output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		subnet = append(subnet, ipv4Prefix(output["address"].(string), output["prefix"].(int), false))
		counter++
	}
	// utilised capacity should be 202+4+32+8+4+2+4=256
	// free capacity should be 0
	// utilisation should be 100%

	// Round 3, no more capacity at utilisation 100%
	userInput := map[string]interface{}{"desiredSize": 2}
	ipv4PrefixStruct := src.NewIpv4Prefix(subnet, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Invoke()
	expectedOutputError := errors.New("Unable to allocate Ipv4 prefix from: 192.168.1.0/24. Insufficient capacity to allocate a new prefix of size: 2\n" +
		"Currently allocated addresses: 192.168.1.0, 192.168.1.32, 192.168.1.64, 192.168.1.16, 192.168.1.128, 192.168.1.24, 192.168.1.192, 192.168.1.20, 192.168.1.160, 192.168.1.144, 192.168.1.152, 192.168.1.18, 192.168.1.156, ")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}

// This test is the same as "ipv4 prefix allocation 24" everything is just multiplied by 256*256 to simplify the assertions
func TestIpv4PrefixAllocation8(t *testing.T) {
	var resourcePool = map[string]interface{}{"prefix": 8, "address": "10.0.0.0", "subnet": false}
	var subnet []map[string]interface{}
	var expectedSubnets = []map[string]interface{}{
		ipv4Prefix("10.0.0.0", 12, false),
		ipv4Prefix("10.32.0.0", 11, false),
		ipv4Prefix("10.64.0.0", 10, false),
		ipv4Prefix("10.16.0.0", 15, false),
		ipv4Prefix("10.128.0.0", 12, false),
		ipv4Prefix("10.24.0.0", 13, false),
		ipv4Prefix("10.192.0.0", 10, false),
	}
	counter := 0
	for _, i := range []int{655360, 1245184, 2555904, 131072, 917504, 524288, 4194304} {
		userInput := map[string]interface{}{"desiredSize": i}
		ipv4PrefixStruct := src.NewIpv4Prefix(subnet, resourcePool, userInput)
		output, err := ipv4PrefixStruct.Invoke()
		if eq := reflect.DeepEqual(output, expectedSubnets[counter]["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedSubnets[counter], output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		subnet = append(subnet, ipv4Prefix(output["address"].(string), output["prefix"].(int), false))
		counter++
	}
	// utilisation should be 202/256=78.9%

	// Round 2, try to squeeze in additional subnets
	var expectedSubnets2 = []map[string]interface{}{
		ipv4Prefix("10.20.0.0", 14, false),
		ipv4Prefix("10.160.0.0", 11, false),
		ipv4Prefix("10.144.0.0", 13, false),
		ipv4Prefix("10.152.0.0", 14, false),
		ipv4Prefix("10.18.0.0", 15, false),
		ipv4Prefix("10.156.0.0", 14, false),
	}
	counter = 0
	for _, i := range []int{262144, 2097152, 524288, 262144, 131072, 262144} {
		userInput := map[string]interface{}{"desiredSize": i}
		ipv4PrefixStruct := src.NewIpv4Prefix(subnet, resourcePool, userInput)
		output, err := ipv4PrefixStruct.Invoke()
		if eq := reflect.DeepEqual(output, expectedSubnets2[counter]["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedSubnets2[counter], output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		subnet = append(subnet, ipv4Prefix(output["address"].(string), output["prefix"].(int), false))
		counter++
	}

	// utilisation should be 100%

	// Round 3, no more capacity at utilisation 100%
	userInput := map[string]interface{}{"desiredSize": 2}
	ipv4PrefixStruct := src.NewIpv4Prefix(subnet, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Invoke()
	expectedOutputError := errors.New("Unable to allocate Ipv4 prefix from: 10.0.0.0/8. Insufficient capacity to allocate a new prefix of size: 2\n" +
		"Currently allocated addresses: 10.0.0.0, 10.32.0.0, 10.64.0.0, 10.16.0.0, 10.128.0.0, 10.24.0.0, 10.192.0.0, 10.20.0.0, 10.160.0.0, 10.144.0.0, 10.152.0.0, 10.18.0.0, 10.156.0.0, ")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}

func TestDesiredSizeMoreThanRoot(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput = map[string]interface{}{"desiredSize": 300}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Invoke()
	expectedOutputError := errors.New("Unable to allocate Ipv4 prefix from: 192.168.1.0/24. Insufficient capacity to allocate a new prefix of size: 300\nCurrently allocated addresses: ")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}

func TestDesiredSizeEqualsRoot(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput = map[string]interface{}{"desiredSize": 256}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Invoke()
	expectedOutput := map[string]interface{}{"address": "192.168.1.0", "prefix": 24, "subnet": false}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestFreeIpv4PrefixCapacity(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": true}
	var userInput = map[string]interface{}{"desiredSize": 2, "subnet": true}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output := ipv4PrefixStruct.FreeCapacity(
		ipv4Prefix("192.168.1.0", 24, false)["Properties"].(map[string]interface{}),
		100)
	expectedOutput := 156
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
}

func TestFreeIpv4PrefixUtilisation(t *testing.T) {
	var allocated = []map[string]interface{}{
		ipv4Prefix("192.168.1.0", 28, false)["Properties"].(map[string]interface{}),
		ipv4Prefix("192.168.1.128", 27, false)["Properties"].(map[string]interface{})}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": true}
	var userInput = map[string]interface{}{"desiredSize": 2, "subnet": true}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output := ipv4PrefixStruct.UtilizedCapacity(
		allocated,
		32)
	expectedOutput := 16 + 32 + 32
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %d expected, got: %d", expectedOutput, output)
	}
}

func TestClaimResourceWithDesiredValue(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput = map[string]interface{}{"desiredSize": 8, "desiredValue": "192.168.1.32"}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, err := ipv4PrefixStruct.Invoke()

	fmt.Printf("error when claiming resource with desired value: %s\n", err)

	expectedOutput := map[string]interface{}{"address": "192.168.1.32", "prefix": 29, "subnet": false}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestClaimResourceWithoutDesiredValue(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4Prefix("192.168.1.0", 29, false)}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput = map[string]interface{}{"desiredSize": 8}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, _ := ipv4PrefixStruct.Invoke()

	expectedOutput := map[string]interface{}{"address": "192.168.1.8", "prefix": 29, "subnet": false}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestClaimResourceInPoolWithAllocatedResources(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4Prefix("192.168.1.40", 29, false), ipv4Prefix("192.168.1.0", 29, false)}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput = map[string]interface{}{"desiredSize": 8, "desiredValue": "192.168.1.16"}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, _ := ipv4PrefixStruct.Invoke()

	expectedOutput := map[string]interface{}{"address": "192.168.1.16", "prefix": 29, "subnet": false}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestOverlappingResourcesWithDesiredValue(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4Prefix("192.168.1.128", 25, false)}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0", "subnet": false}
	var userInput = map[string]interface{}{"desiredSize": 256, "desiredValue": "192.168.1.0"}
	ipv4PrefixStruct := src.NewIpv4Prefix(allocated, resourcePool, userInput)
	output, _ := ipv4PrefixStruct.Invoke()

	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}
