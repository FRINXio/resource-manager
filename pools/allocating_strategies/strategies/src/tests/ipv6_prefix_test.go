package tests

import (
	"fmt"
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"testing"
)

func ipv6PrefixWithSubnet(address string, prefix int, subnet bool) map[string]interface{} {
	ipv6PrefixProperties := make(map[string]interface{})
	ipv6PrefixMap := make(map[string]interface{})
	ipv6PrefixMap["address"] = address
	ipv6PrefixMap["prefix"] = prefix
	ipv6PrefixMap["subnet"] = subnet
	ipv6PrefixProperties["Properties"] = ipv6PrefixMap
	return ipv6PrefixProperties
}

func ipv6Prefix(address string, prefix int) map[string]interface{} {
	ipv6PrefixProperties := make(map[string]interface{})
	ipv6PrefixMap := make(map[string]interface{})
	ipv6PrefixMap["address"] = address
	ipv6PrefixMap["prefix"] = prefix
	ipv6PrefixProperties["Properties"] = ipv6PrefixMap
	return ipv6PrefixProperties
}

func TestParseIpv6Address(t *testing.T) {
	var ipv6Addresses = []string{"dead::beef",
		"::1",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
		"::",
		"a897:fedc:1111:9999:f999::abcd"}

	for _, ipv6Address := range ipv6Addresses {
		ipv6AddressAton, _ := src.Ipv6InetAton(ipv6Address)
		if eq := reflect.DeepEqual(ipv6Address, src.Ipv6InetNtoa(ipv6AddressAton)); !eq {
			t.Fatalf("different output of %s expected, got: %s", ipv6Address, src.Ipv6InetNtoa(ipv6AddressAton))
		}
	}
}

func TestParseIpv6AddressPrefix(t *testing.T) {
	var expectedIpv6PrefixAddresses = []map[string]interface{}{
		ipv6Prefix("dead::", 64),
		ipv6Prefix("::", 19),
		ipv6Prefix("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", 128),
		ipv6Prefix("ffff:ffff:ffff:ffff:ffff:ffff:ffff:0", 112),
		ipv6Prefix("ff00::", 8),
		ipv6Prefix("ffff::", 16),
		ipv6Prefix("a897:fedc:1111:9999:f999:abcc::", 95),
	}
	var ipv6PrefixAddresses = []string{
		"dead::beef/64",
		"::1/19",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/112",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/8",
		"ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/16",
		"a897:fedc:1111:9999:f999:abcd::/95"}

	for index, ipv6PrefixAddress := range ipv6PrefixAddresses {
		ipv6ParsePrefixAddress, err := src.Ipv6ParsePrefix(ipv6PrefixAddress)
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}

		if eq := reflect.DeepEqual(ipv6ParsePrefixAddress, expectedIpv6PrefixAddresses[index]["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", ipv6ParsePrefixAddress, expectedIpv6PrefixAddresses[index]["Properties"])
		}
	}
}

func TestParseInvalidIpv6Address(t *testing.T) {
	var ipv6Addresses = []string{
		"xxxx::yyyy",
		"z",
		"888878468945",
	}

	for _, ipv6Address := range ipv6Addresses {
		ipv6AddressAton, err := src.Ipv6InetAton(ipv6Address)
		expected := errors.New("Address: " + ipv6Address + " is invalid ipv6 address.")
		if eq := reflect.DeepEqual(expected.Error(), err.Error()); !eq {
			t.Fatalf("different output of %s expected, got: %s", expected, err)
		}
		if eq := reflect.DeepEqual(big.NewInt(-1), ipv6AddressAton); !eq {
			t.Fatalf("different output of %s expected, got: %s", big.NewInt(-1), src.Ipv6InetNtoa(ipv6AddressAton))
		}
	}
}

func TestSingleAllocationPoolIpv6(t *testing.T) {
	var desiredSize *big.Int
	for i := 1; i <= 128-8-1; i++ {
		desiredSize = new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
		resourcePool := map[string]interface{}{"prefix": 8, "address": "bb00::", "subnet": true}
		userInput := map[string]interface{}{"desiredSize": desiredSize}
		ipv6PrefixStruct := src.NewIpv6Prefix([]map[string]interface{}{}, resourcePool, userInput)
		output, err := ipv6PrefixStruct.Invoke()
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		if eq := reflect.DeepEqual(output, ipv6PrefixWithSubnet("bb00::", 128-i-1, true)["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", output,
				ipv6PrefixWithSubnet("bb00::", 128-i-1, true)["Properties"])
		}
	}
}

func TestAllocateRangeAtStartWithExistingResourcesIpv6(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 2}
	var allocated []map[string]interface{}
	allocated = append(allocated, ipv6Prefix("dead::be02", 127))
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, err := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	if eq := reflect.DeepEqual(output, ipv6PrefixWithSubnet("dead::be00", 127, false)["Properties"]); !eq {
		t.Fatalf("different output of %s expected, got: %s",
			ipv6PrefixWithSubnet("dead::be00", 127, false)["Properties"], output)
	}
}

func TestDesiredSizeIsBiggerThanRootIpv6(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 300}
	ipv6PrefixStruct := src.NewIpv6Prefix([]map[string]interface{}{}, resourcePool, userInput)
	output, err := ipv6PrefixStruct.Invoke()
	expected := errors.New("Unable to allocate Ipv6 prefix from: dead::be00/120. " +
		"Insufficient capacity to allocate a new prefix of size: 300. Currently allocated prefixes: ")
	if eq := reflect.DeepEqual(expected.Error(), err.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expected.Error(), err.Error())
	}
	if eq := reflect.DeepEqual(map[string]interface{}(nil), output); !eq {
		t.Fatalf("different output of %s expected, got: %s", map[string]interface{}(nil), output)
	}
}

func TestDesiredSizeIsTheSameThanRootIpv6(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 104, "address": "dead::", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 16777216}
	ipv6PrefixStruct := src.NewIpv6Prefix([]map[string]interface{}{}, resourcePool, userInput)
	output, err := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	if eq := reflect.DeepEqual(output, ipv6PrefixWithSubnet("dead::", 104, false)["Properties"]); !eq {
		t.Fatalf("different output of %s expected, got: %s",
			ipv6PrefixWithSubnet("dead::", 104, false)["Properties"], output)
	}
}

func TestCapacity(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 16777216}
	ipv6PrefixStruct := src.NewIpv6Prefix([]map[string]interface{}{}, resourcePool, userInput)
	output, err := ipv6PrefixStruct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": "256", "utilizedCapacity": "0"}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
}

func TestAllocateAllAddresses104Ipv6Range(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 104, "address": "abcd:ef01:2345:6789::", "subnet": false}
	var expectedIpv6PrefixAddresses = []map[string]interface{}{
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::", 108, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::20:0", 107, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::40:0", 106, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::10:0", 111, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::80:0", 108, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::18:0", 109, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::c0:0", 106, false),
	}
	var ipv6PrefixAddresses = []int64{655360, 1245184, 2555904, 131072, 917504, 524288, 4194304}
	for index, ipv6PrefixAddress := range expectedIpv6PrefixAddresses {
		ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool,
			map[string]interface{}{"desiredSize": ipv6PrefixAddresses[index]})
		output, err := ipv6PrefixStruct.Invoke()
		allocated = append(allocated, ipv6PrefixWithSubnet(output["address"].(string), output["prefix"].(int), false))
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}

		if eq := reflect.DeepEqual(output, ipv6PrefixAddress["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", ipv6PrefixAddress["Properties"], output)
		}
	}

	// utilisation should be 78.9%

	// Round 2, try to squeeze in additional subnets
	var expectedIpv6PrefixAddresses2 = []map[string]interface{}{
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::14:0", 110, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::a0:0", 107, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::90:0", 109, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::98:0", 110, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::12:0", 111, false),
		ipv6PrefixWithSubnet("abcd:ef01:2345:6789::9c:0", 110, false),
	}
	var ipv6PrefixAddresses2 = []int64{262144, 2097152, 524288, 262144, 131072, 262144}
	for index, ipv6PrefixAddress := range expectedIpv6PrefixAddresses2 {
		ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool,
			map[string]interface{}{"desiredSize": ipv6PrefixAddresses2[index]})
		output, err := ipv6PrefixStruct.Invoke()
		allocated = append(allocated, ipv6PrefixWithSubnet(output["address"].(string), output["prefix"].(int), false))
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}

		if eq := reflect.DeepEqual(output, ipv6PrefixAddress["Properties"]); !eq {
			t.Fatalf("different output of %s expected, got: %s", ipv6PrefixAddress["Properties"], output)
		}
	}

	// utilisation should be 100%

	// Round 3, no more capacity at utilisation 100%
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool,
		map[string]interface{}{"desiredSize": 2})
	output, err := ipv6PrefixStruct.Invoke()
	expectedErr := errors.New("Unable to allocate Ipv6 prefix from: abcd:ef01:2345:6789::/104. " +
		"Insufficient capacity to allocate a new prefix of size: 2. Currently allocated prefixes: " +
		"abcd:ef01:2345:6789::/108[abcd:ef01:2345:6789::-abcd:ef01:2345:6789::20:0], " +
		"abcd:ef01:2345:6789::10:0/111[abcd:ef01:2345:6789::10:0-abcd:ef01:2345:6789::14:0], " +
		"abcd:ef01:2345:6789::12:0/111[abcd:ef01:2345:6789::12:0-abcd:ef01:2345:6789::16:0], " +
		"abcd:ef01:2345:6789::14:0/110[abcd:ef01:2345:6789::14:0-abcd:ef01:2345:6789::1c:0], " +
		"abcd:ef01:2345:6789::18:0/109[abcd:ef01:2345:6789::18:0-abcd:ef01:2345:6789::28:0], " +
		"abcd:ef01:2345:6789::20:0/107[abcd:ef01:2345:6789::20:0-abcd:ef01:2345:6789::60:0], " +
		"abcd:ef01:2345:6789::40:0/106[abcd:ef01:2345:6789::40:0-abcd:ef01:2345:6789::c0:0], " +
		"abcd:ef01:2345:6789::80:0/108[abcd:ef01:2345:6789::80:0-abcd:ef01:2345:6789::a0:0], " +
		"abcd:ef01:2345:6789::90:0/109[abcd:ef01:2345:6789::90:0-abcd:ef01:2345:6789::a0:0], " +
		"abcd:ef01:2345:6789::98:0/110[abcd:ef01:2345:6789::98:0-abcd:ef01:2345:6789::a0:0], " +
		"abcd:ef01:2345:6789::9c:0/110[abcd:ef01:2345:6789::9c:0-abcd:ef01:2345:6789::a4:0], " +
		"abcd:ef01:2345:6789::a0:0/107[abcd:ef01:2345:6789::a0:0-abcd:ef01:2345:6789::e0:0], " +
		"abcd:ef01:2345:6789::c0:0/106[abcd:ef01:2345:6789::c0:0-abcd:ef01:2345:6789::140:0], ")

	if eq := reflect.DeepEqual(err.Error(), expectedErr.Error()); !eq {
		fmt.Println(err.Error())
		fmt.Println(expectedErr.Error())
		t.Fatalf("different output of %s expected, got: %s", expectedErr.Error(), err.Error())
	}

	if eq := reflect.DeepEqual(output, map[string]interface{}(nil)); !eq {
		t.Fatalf("different output of %s expected, got: %s", map[string]interface{}(nil), output)
	}

	// Capacity test
	output, err = ipv6PrefixStruct.Capacity()
	expected := map[string]interface{}{"freeCapacity": "0", "utilizedCapacity": "16777216"}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	if eq := reflect.DeepEqual(output, expected); !eq {
		t.Fatalf("different output of %s expected, got: %s", expected, output)
	}
}

func TestAllocationOfResourceWithDesiredValueInBetweenResources(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 2, "desiredValue": "dead::be04"}
	var allocated []map[string]interface{}
	allocated = append(allocated, ipv6Prefix("dead::be02", 127))
	allocated = append(allocated, ipv6Prefix("dead::be08", 127))
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, err := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	if eq := reflect.DeepEqual(output, ipv6PrefixWithSubnet("dead::be04", 127, false)["Properties"]); !eq {
		t.Fatalf("different output of %s expected, got: %s",
			ipv6PrefixWithSubnet("dead::be00", 127, false)["Properties"], output)
	}
}

func TestAllocationOfResourceWithDesiredValueAfterResources(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 2, "desiredValue": "dead::be10"}
	var allocated []map[string]interface{}
	allocated = append(allocated, ipv6Prefix("dead::be02", 127))
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, err := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	if eq := reflect.DeepEqual(output, ipv6PrefixWithSubnet("dead::be10", 127, false)["Properties"]); !eq {
		t.Fatalf("different output of %s expected, got: %s",
			ipv6PrefixWithSubnet("dead::be00", 127, false)["Properties"], output)
	}
}

func TestAllocationOfResourceWithDesiredValueWithoutResources(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 2, "desiredValue": "dead::be30"}
	var allocated []map[string]interface{}
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, err := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	if eq := reflect.DeepEqual(output, ipv6PrefixWithSubnet("dead::be30", 127, false)["Properties"]); !eq {
		t.Fatalf("different output of %s expected, got: %s",
			ipv6PrefixWithSubnet("dead::be00", 127, false)["Properties"], output)
	}
}

func TestAllocationOfResourceWithInvalidDesiredValue(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 2, "desiredValue": "halabala"}
	var allocated []map[string]interface{}
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, _ := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(output, map[string]interface{}(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}

	userInput = map[string]interface{}{"desiredSize": 2, "desiredValue": "dead::be03"}
	ipv6PrefixStruct = src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, _ = ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(output, map[string]interface{}(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}

	userInput = map[string]interface{}{"desiredSize": 2, "desiredValue": "dead::be00"}
	allocated = append(allocated, ipv6Prefix("dead::be00", 127))
	ipv6PrefixStruct = src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, _ = ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(output, map[string]interface{}(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}

	userInput = map[string]interface{}{"desiredSize": 2, "desiredValue": "dead::be02"}
	var newAllocated []map[string]interface{}
	newAllocated = append(allocated, ipv6Prefix("dead::be02", 124))
	ipv6PrefixStruct = src.NewIpv6Prefix(newAllocated, resourcePool, userInput)
	output, _ = ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(output, map[string]interface{}(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestAllocationOfResourceWithOverflowingCapacity(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 126, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 5, "desiredValue": "dead::be00"}
	var allocated []map[string]interface{}
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, _ := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(output, map[string]interface{}(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestAllocationOfResourceWithDesiredValueNotNetwork(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 125, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 4, "desiredValue": "dead::be02"}
	var allocated []map[string]interface{}
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	output, _ := ipv6PrefixStruct.Invoke()
	if eq := reflect.DeepEqual(output, map[string]interface{}(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestCorrectlyCalculatedIPv6NetworkAddressForErrorMsg(t *testing.T) {
	resourcePool := map[string]interface{}{"prefix": 120, "address": "dead::be00", "subnet": false}
	userInput := map[string]interface{}{"desiredSize": 4, "desiredValue": "dead::be02"}
	var allocated = []map[string]interface{}{ipv6PrefixWithSubnet("dead::be00", 127, false), ipv6PrefixWithSubnet("dead::be10", 127, false)}
	ipv6PrefixStruct := src.NewIpv6Prefix(allocated, resourcePool, userInput)
	_, err := ipv6PrefixStruct.Invoke()
	expectedErrorOutput := errors.New("You provided invalid network address. Network address should be dead::be04")

	if eq := reflect.DeepEqual(err.Error(), expectedErrorOutput.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedErrorOutput, err)
	}
}
