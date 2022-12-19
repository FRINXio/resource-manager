package src

import (
	"github.com/pkg/errors"
	"strconv"
)

type Ipv4 struct {
	currentResources       []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

func NewIpv4(currentResources []map[string]interface{},
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) Ipv4 {
	return Ipv4{currentResources, resourcePoolProperties, userInput}
}

func (ipv4 *Ipv4) UtilizedCapacity(allocatedRanges []map[string]interface{}, newlyAllocatedRangeCapacity float64) float64 {
	return float64(len(allocatedRanges)) + newlyAllocatedRangeCapacity
}

// FreeCapacity calculate free capacity based on previously allocated prefixes
func (ipv4 *Ipv4) FreeCapacity(address string, mask int, utilisedCapacity float64, subnetItself int) float64 {
	return float64(hostsInMask(address, mask)) - utilisedCapacity + float64(subnetItself)
}

func (ipv4 *Ipv4) Capacity() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	rootAddressStr, ok := ipv4.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv4.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	subnet, ok := ipv4.resourcePoolProperties["subnet"]
	if !ok {
		return nil, errors.New("Unable to extract subnet property")
	}
	subnetItself := 2
	if subnet.(bool) == true {
		subnetItself = 0
	}
	freeCapacity := ipv4.FreeCapacity(rootAddressStr.(string), rootMask.(int), float64(len(ipv4.currentResources)), subnetItself)
	result["freeCapacity"] = strconv.FormatFloat(freeCapacity, 'g', 30, 64)
	result["utilizedCapacity"] = strconv.Itoa(len(ipv4.currentResources))
	return result, nil
}

func (ipv4 *Ipv4) Invoke() (map[string]interface{}, error) {
	if ipv4.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	rootAddressStr, ok := ipv4.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv4.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	isSubnet, ok := ipv4.resourcePoolProperties["subnet"]
	if !ok {
		return nil, errors.New("Unable to extract subnet property")
	}
	rootMask, err := NumberToInt(rootMask)
	if err != nil {
		return nil, err
	}
	rootPrefixStr := prefixToStr(ipv4.resourcePoolProperties)
	rootCapacity := subnetAddresses(rootMask.(int))
	rootAddressNum, err := inetAton(rootAddressStr.(string))
	if err != nil {
		return nil, err
	}

	// unwrap and create currentResourcesSet
	currentResourcesSet := make(map[string]bool)
	for _, element := range ipv4.currentResources {
		value, ok := element["Properties"]
		if !ok {
			return nil, errors.New("Wrong properties in current resources")
		}
		if address, ok := value.(map[string]interface{})["address"]; ok {
			currentResourcesSet[address.(string)] = true
		}
	}

	var firstPossibleAddr = 0
	var lastPossibleAddr = 0

	if isSubnet.(bool) == true {
		firstPossibleAddr = rootAddressNum + 1
		lastPossibleAddr = rootAddressNum + rootCapacity - 1
	} else {
		firstPossibleAddr = rootAddressNum
		lastPossibleAddr = rootAddressNum + rootCapacity
	}

	var result = make(map[string]interface{})

	if value, ok := ipv4.userInput["desiredValue"]; ok {
		desiredValueNum, err := inetAton(value.(string))
		if err != nil {
			return nil, err
		}
		if desiredValueNum >= firstPossibleAddr && desiredValueNum < lastPossibleAddr {
			desiredIpv4Address := inetNtoa(desiredValueNum)
			if currentResourcesSet[desiredIpv4Address] {
				return nil, errors.New("Ipv4 address " + value.(string) + " was already claimed.")
			}
			result["address"] = value.(string)
			return result, nil
		} else {
			return nil, errors.New("Ipv4 address " + value.(string) + " is out of " + rootPrefixStr)
		}
	}

	for i := firstPossibleAddr; i < lastPossibleAddr; i++ {
		if !currentResourcesSet[inetNtoa(i)] {
			result["address"] = inetNtoa(i)
			return result, nil
		}
	}
	return nil, errors.New("Unable to allocate Ipv4 address from: " + rootPrefixStr + "." +
		"Insufficient capacity to allocate a new address.\n" +
		"Currently allocated addresses: " + addressesToStr(ipv4.currentResources))
}
