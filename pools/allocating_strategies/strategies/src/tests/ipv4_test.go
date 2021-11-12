package tests

import (
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"testing"
)

func ipv4(address string) map[string]interface{} {
	ipv4Properties := make(map[string]interface{})
	ipv4Map := make(map[string]interface{})
	ipv4Map["address"] = address
	ipv4Properties["Properties"] = ipv4Map
	return ipv4Properties
}

func TestAllocateAllAddresses24(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv4Struct := src.NewIpv4(allocated, resourcePool, userInput)

	var generated = ""
	for i := 1; i < 255; i++ {
		output, err := ipv4Struct.Invoke()
		expectedOutput := map[string]interface{}{"address": "192.168.1." + strconv.Itoa(i)}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		allocated = append(allocated, ipv4(output["address"].(string)))
		ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)

		generated += output["address"].(string) + ", "
	}

	// If treated as subnet, prefix is exhausted
	output, err := ipv4Struct.Invoke()
	expectedOutputError := errors.New("Unable to allocate Ipv4 address from: 192.168.1.0/24." +
		"Insufficient capacity to allocate a new address.\n" +
		"Currently allocated addresses: " + generated)
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// If treated as a pool, there are still 2 more addresses left
	resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0"}
	userInput = map[string]interface{}{}
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
	output, err = ipv4Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "192.168.1.0"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv4(output["address"].(string)))
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
	generated += output["address"].(string) + ", "

	output, err = ipv4Struct.Invoke()
	expectedOutput = map[string]interface{}{"address": "192.168.1.255"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv4(output["address"].(string)))
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
	generated += output["address"].(string) + ", "

	output, err = ipv4Struct.Invoke()
	expectedOutputError = errors.New("Unable to allocate Ipv4 address from: 192.168.1.0/24." +
		"Insufficient capacity to allocate a new address.\n" +
		"Currently allocated addresses: " + generated)
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}

func TestAllocateAllAddresses19(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 19, "address": "192.168.0.0"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv4Struct := src.NewIpv4(allocated, resourcePool, userInput)
	for  i := 0; i < 32; i++ {
		for j := 0; j < 256; j++ {
			if i == 0 && j == 0 {
				// First subnet addr: reserved
				continue
			}
			if i == 31 && j == 255 {
				// Broadcast: reserved
				continue
			}
			output, err := ipv4Struct.Invoke()
			expectedOutput := map[string]interface{}{"address": "192.168." + strconv.Itoa(i) + "." + strconv.Itoa(j)}
			if eq := reflect.DeepEqual(output, expectedOutput); !eq {
				t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
			}
			if eq := reflect.DeepEqual(err, nil); !eq {
				t.Fatalf("different output of nil expected, got: %s", err)
			}
			allocated = append(allocated, ipv4(output["address"].(string)))
			ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
		}
	}

	// If treated as subnet, prefix is exhausted
	output, _ := ipv4Struct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}

	// If treated as a pool, there are still 2 more addresses left
	resourcePool = map[string]interface{}{"prefix": 19, "address": "192.168.0.0"}
	userInput = map[string]interface{}{}
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
	output, err := ipv4Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "192.168.0.0"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv4(output["address"].(string)))
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)

	output, err = ipv4Struct.Invoke()
	expectedOutput = map[string]interface{}{"address": "192.168.31.255"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv4(output["address"].(string)))
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)

	output, _ = ipv4Struct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestAllocateIpv4AtStartWithExistingResources(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4("192.168.1.2")}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv4Struct := src.NewIpv4(allocated, resourcePool, userInput)
	output, err := ipv4Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "192.168.1.1"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestIpv4Capacity24Mask(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4("192.168.1.2")}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv4Struct := src.NewIpv4(allocated, resourcePool, userInput)
	output, err := ipv4Struct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": float64(254), "utilizedCapacity": float64(1)}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestIpv4Capacity16Mask(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4("192.168.1.2")}
	var resourcePool = map[string]interface{}{"prefix": 16, "address": "192.168.1.0"}
	var userInput map[string]interface{}
	ipv4Struct := src.NewIpv4(allocated, resourcePool, userInput)
	output, err := ipv4Struct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": float64(65533), "utilizedCapacity": float64(1)}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestIpv4Utilisation(t *testing.T) {
	var allocated = []map[string]interface{}{ipv4("192.168.1.2")}
	var resourcePool = map[string]interface{}{"prefix": 16, "address": "192.168.1.0"}
	var userInput map[string]interface{}
	ipv4Struct := src.NewIpv4(allocated, resourcePool, userInput)
	output := ipv4Struct.UtilizedCapacity(allocated, float64(1))
	expectedOutput := float64(2)
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %f expected, got: %f", expectedOutput, output)
	}
}

func TestIpv4DesireValue(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "192.168.1.0"}
	var userInput = map[string]interface{}{"desiredValue": "192.168.1.5"}
	ipv4Struct := src.NewIpv4(allocated, resourcePool, userInput)
	output, err := ipv4Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "192.168.1.5"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	// DesiredValue out of address/prefix
	userInput = map[string]interface{}{"desiredValue": "192.168.2.5"}
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
	output, err = ipv4Struct.Invoke()
	expectedOutputError := errors.New("Ipv4 address 192.168.2.5 is out of 192.168.1.0/24")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// Already claimed desiredValue
	allocated = []map[string]interface{}{ipv4("192.168.1.5")}
	userInput = map[string]interface{}{"desiredValue": "192.168.1.5"}
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
	output, err = ipv4Struct.Invoke()
	expectedOutputError = errors.New("Ipv4 address 192.168.1.5 was already claimed.")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// Wrong desiredValue input format
	allocated = []map[string]interface{}{}
	userInput = map[string]interface{}{"desiredValue": "Hello World"}
	ipv4Struct = src.NewIpv4(allocated, resourcePool, userInput)
	output, err = ipv4Struct.Invoke()
	expectedOutputError = errors.New("Address: Hello World is invalid, doesn't match regex: ^([0-9]{1,3})\\.([0-9]{1,3})\\.([0-9]{1,3})\\.([0-9]{1,3})$")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}
