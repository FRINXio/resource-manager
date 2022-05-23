package tests

import (
	"fmt"
	"github.com/net-auto/resourceManager/pools/allocating_strategies/strategies/src"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"testing"
)

func ipv6(address string) map[string]interface{} {
	ipv6Properties := make(map[string]interface{})
	ipv6Map := make(map[string]interface{})
	ipv6Map["address"] = address
	ipv6Properties["Properties"] = ipv6Map
	return ipv6Properties
}

func TestAllocateAllAddresses120Ipv6(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 120, "address": "dddd::"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv6Struct := src.NewIpv6(allocated, resourcePool, userInput)

	var generated = ""
	for i := 1; i < 255; i++ {
		output, err := ipv6Struct.Invoke()
		expectedOutput := map[string]interface{}{"address": "dddd::" + fmt.Sprintf("%x", i)}
		if eq := reflect.DeepEqual(output, expectedOutput); !eq {
			t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
		}
		if eq := reflect.DeepEqual(err, nil); !eq {
			t.Fatalf("different output of nil expected, got: %s", err)
		}
		allocated = append(allocated, ipv6(output["address"].(string)))
		ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)

		generated += output["address"].(string) + ", "
	}

	// If treated as subnet, prefix is exhausted
	output, err := ipv6Struct.Invoke()
	expectedOutputError := errors.New("Unable to allocate Ipv6 address from: dddd::/120." +
		"Insufficient capacity to allocate a new address.\n" +
		"Currently allocated addresses: " + generated)
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// If treated as a pool, there are still 2 more addresses left
	resourcePool = map[string]interface{}{"prefix": 120, "address": "dddd::"}
	userInput = map[string]interface{}{}
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	output, err = ipv6Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "dddd::"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv6(output["address"].(string)))
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	generated += output["address"].(string) + ", "

	output, err = ipv6Struct.Invoke()
	expectedOutput = map[string]interface{}{"address": "dddd::ff"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv6(output["address"].(string)))
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	generated += output["address"].(string) + ", "

	output, err = ipv6Struct.Invoke()
	expectedOutputError = errors.New("Unable to allocate Ipv6 address from: dddd::/120." +
		"Insufficient capacity to allocate a new address.\n" +
		"Currently allocated addresses: " + generated)
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}

func TestAllocateAllAddresses117ForIpv6(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool = map[string]interface{}{"prefix": 115, "address": "dddd::"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv6Struct := src.NewIpv6(allocated, resourcePool, userInput)
	for  i := 0; i < 8; i++ {
		for j := 0; j < 256; j++ {
			if i == 0 && j == 0 {
				// First subnet addr: reserved
				continue
			}
			if i == 7 && j == 255 {
				// Broadcast: reserved
				continue
			}
			output, err := ipv6Struct.Invoke()

			// some formatting adjustments
			byteI := fmt.Sprintf("%x", i)
			byteJ := fmt.Sprintf("%x", j)
			if byteI == "0" {
				byteI = ""
				if byteJ == "0" {
					byteJ = ""
				}
			} else if len(byteJ) == 1 {
				byteJ = "0" + byteJ
			}

			expectedOutput := map[string]interface{}{"address": "dddd::" + byteI + "" + byteJ}
			if eq := reflect.DeepEqual(output, expectedOutput); !eq {
				t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
			}
			if eq := reflect.DeepEqual(err, nil); !eq {
				t.Fatalf("different output of nil expected, got: %s", err)
			}
			allocated = append(allocated, ipv6(output["address"].(string)))
			ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
		}
	}

	// If treated as subnet, prefix is exhausted
	resourcePool = map[string]interface{}{"prefix": 117, "address": "dddd::"}
	userInput = map[string]interface{}{"subnet": true}
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	output, _ := ipv6Struct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}

	// If treated as a pool, there are still 2 more addresses left
	userInput = map[string]interface{}{}
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	output, err := ipv6Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "dddd::"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv6(output["address"].(string)))
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)

	output, err = ipv6Struct.Invoke()
	expectedOutput = map[string]interface{}{"address": "dddd::7ff"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
	allocated = append(allocated, ipv6(output["address"].(string)))
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)

	output, _ = ipv6Struct.Invoke()
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
}

func TestAllocateIpv6AtStartWithExistingResources(t *testing.T) {
	var allocated = []map[string]interface{}{ipv6("dead::2")}
	var resourcePool = map[string]interface{}{"prefix": 24, "address": "dead::"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv6Struct := src.NewIpv6(allocated, resourcePool, userInput)
	output, err := ipv6Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "dead::1"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestIpv6Capacity24Mask(t *testing.T) {
	var allocated = []map[string]interface{}{ipv6("dead::2")}
	var resourcePool = map[string]interface{}{"prefix": 110, "address": "dead::"}
	var userInput = map[string]interface{}{"subnet": true}
	ipv6Struct := src.NewIpv6(allocated, resourcePool, userInput)
	output, err := ipv6Struct.Capacity()
	expectedOutput := map[string]interface{}{"freeCapacity": "262144", "utilizedCapacity": "1"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}
}

func TestIpv6Capacity16Mask(t *testing.T) {
	var allocated []map[string]interface{}
	var resourcePool map[string]interface{}
	var userInput map[string]interface{}
	ipv6Struct := src.NewIpv6(allocated, resourcePool, userInput)
	output := ipv6Struct.FreeCapacity("120", float64(100))
	expectedOutput := float64(156)
	log.Println(output)
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %f expected, got: %f", expectedOutput, output)
	}
}

func TestIpv6Utilisation(t *testing.T) {
	var allocated = []map[string]interface{}{ipv6("dead::beed")}
	var resourcePool map[string]interface{}
	var userInput map[string]interface{}
	ipv6Struct := src.NewIpv6(allocated, resourcePool, userInput)
	output := ipv6Struct.UtilizedCapacity(allocated, float64(1))
	expectedOutput := float64(2)
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %f expected, got: %f", expectedOutput, output)
	}
}

func TestIpv6DesireValue(t *testing.T) {
	var allocated = []map[string]interface{}{ipv6("dead::2")}
	var resourcePool = map[string]interface{}{"prefix": 110, "address": "dead::"}
	var userInput = map[string]interface{}{"desiredValue": "dead::3"}
	ipv6Struct := src.NewIpv6(allocated, resourcePool, userInput)
	output, err := ipv6Struct.Invoke()
	expectedOutput := map[string]interface{}{"address": "dead::3"}
	if eq := reflect.DeepEqual(output, expectedOutput); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutput, output)
	}
	if eq := reflect.DeepEqual(err, nil); !eq {
		t.Fatalf("different output of nil expected, got: %s", err)
	}

	// DesiredValue out of address/prefix
	userInput = map[string]interface{}{"desiredValue": "deee::2"}
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	output, err = ipv6Struct.Invoke()
	expectedOutputError := errors.New("Ipv6 address deee::2 is out of dead::/110")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// Already claimed desiredValue
	userInput = map[string]interface{}{"desiredValue": "dead::2"}
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	output, err = ipv6Struct.Invoke()
	expectedOutputError = errors.New("Ipv6 address dead::2 was already claimed.")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}

	// Wrong desiredValue input format
	allocated = []map[string]interface{}{}
	userInput = map[string]interface{}{"desiredValue": "Hello World"}
	ipv6Struct = src.NewIpv6(allocated, resourcePool, userInput)
	output, err = ipv6Struct.Invoke()
	expectedOutputError = errors.New("Address: Hello World is invalid ipv6 address.")
	if eq := reflect.DeepEqual(output, (map[string]interface{})(nil)); !eq {
		t.Fatalf("different output of nil expected, got: %s", output)
	}
	if eq := reflect.DeepEqual(err.Error(), expectedOutputError.Error()); !eq {
		t.Fatalf("different output of %s expected, got: %s", expectedOutputError, err)
	}
}
