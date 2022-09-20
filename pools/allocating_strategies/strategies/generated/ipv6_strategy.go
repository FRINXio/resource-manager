package src

import (
	"encoding/json"
	"github.com/pkg/errors"
	"math/big"
	"strconv"
)

type Ipv6 struct {
	currentResources       []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

func NewIpv6(currentResources []map[string]interface{},
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) Ipv6 {
	return Ipv6{currentResources, resourcePoolProperties, userInput}
}

func (ipv6 *Ipv6) UtilizedCapacity(allocatedRanges []map[string]interface{}, newlyAllocatedRangeCapacity float64) float64 {
	return float64(len(allocatedRanges)) + newlyAllocatedRangeCapacity
}

// FreeCapacity calculate free capacity based on previously allocated prefixes
func (ipv6 *Ipv6) FreeCapacity(parentPrefix string, utilisedCapacity float64) float64 {
	capacityString := strconv.FormatFloat(utilisedCapacity, 'f', -1, 64)
	capacityInt, _ := new(big.Int).SetString(capacityString, 10)
	parentPrefixInt, _ := strconv.Atoi(parentPrefix)
	addressesCount := ipv6SubnetAddresses(parentPrefixInt)
	addressesCount.Sub(addressesCount, capacityInt)
	return float64(addressesCount.Int64())
}

func (ipv6 *Ipv6) Capacity() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	rootAddressStr, ok := ipv6.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv6.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	isSubnet, ok := ipv6.resourcePoolProperties["subnet"]
	if !ok {
		return nil, errors.New("Unable to extract subnet property")
	}
	subnetItself := new(big.Int)
	if isSubnet.(bool) == true {
		subnetItself = big.NewInt(-2)
	} else {
		subnetItself = big.NewInt(0)
	}
	rootMask, err := NumberToInt(rootMask)
	if err != nil {
		return nil, err
	}
	freeInTotal := ipv6HostsInMask(rootAddressStr.(string), rootMask.(int))
	freeInTotal.Add(freeInTotal, subnetItself)

	result["freeCapacity"] = freeInTotal.Sub(freeInTotal, big.NewInt(int64(len(ipv6.currentResources)))).String()
	result["utilizedCapacity"] = strconv.Itoa(len(ipv6.currentResources))
	return result, nil
}

func (ipv6 *Ipv6) Invoke() (map[string]interface{}, error) {
	if ipv6.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	rootAddressStr, ok := ipv6.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv6.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	isSubnet, ok := ipv6.resourcePoolProperties["subnet"]
	if !ok {
		return nil, errors.New("Unable to extract subnet property")
	}
	rootMask, err := NumberToInt(rootMask)
	if err != nil {
		return nil, err
	}
	rootPrefixStr := prefixToStr(ipv6.resourcePoolProperties)
	rootCapacity := ipv6SubnetAddresses(rootMask.(int))
	rootAddressNum, err := Ipv6InetAton(rootAddressStr.(string))
	if err != nil {
		return nil, err
	}

	// unwrap and create currentResourcesSet
	currentResourcesSet := make(map[string]bool)
	for _, element := range ipv6.currentResources {
		value, ok := element["Properties"]
		if !ok {
			return nil, errors.New("Wrong properties in current resources")
		}
		if address, ok := value.(map[string]interface{})["address"]; ok {
			currentResourcesSet[address.(string)] = true
		}
	}

	var firstPossibleAddr = big.NewInt(0)
	var lastPossibleAddr = big.NewInt(0)

	if isSubnet.(bool) == true {
		firstPossibleAddr.Add(rootAddressNum, big.NewInt(1))
		lastPossibleAddr.Add(rootAddressNum, rootCapacity)
		lastPossibleAddr.Sub(lastPossibleAddr, big.NewInt(1))
	} else {
		firstPossibleAddr = rootAddressNum
		lastPossibleAddr.Add(rootAddressNum, rootCapacity)
	}

	var result = make(map[string]interface{})

	if value, ok := ipv6.userInput["desiredValue"]; ok {
		desiredValueNum, err := Ipv6InetAton(value.(string))
		if err != nil {
			return nil, err
		}
		if desiredValueNum.Cmp(firstPossibleAddr) >= 0 && desiredValueNum.Cmp(lastPossibleAddr) < 0 {
			desiredIpv6Address := Ipv6InetNtoa(desiredValueNum)
			if currentResourcesSet[desiredIpv6Address] {
				return nil, errors.New("Ipv6 address " + value.(string) + " was already claimed.")
			}
			result["address"] = value.(string)
			return result, nil
		} else {
			return nil, errors.New("Ipv6 address " + value.(string) + " is out of " + rootPrefixStr)
		}
	}

	for i := 0; new(big.Int).Sub(lastPossibleAddr, firstPossibleAddr).Cmp(big.NewInt(int64(i))) > 0; i++ {
		ipv6Address := Ipv6InetNtoa(new(big.Int).Add(firstPossibleAddr, big.NewInt(int64(i))))
		if !currentResourcesSet[ipv6Address] {
			result["address"] = ipv6Address
			return result, nil
		}
	}
	return nil, errors.New("Unable to allocate Ipv6 address from: " + rootPrefixStr + "." +
		"Insufficient capacity to allocate a new address.\n" +
		"Currently allocated addresses: " + addressesToStr(ipv6.currentResources))
}

func NumberToBigInt(number interface{}) (*big.Int, error) {
	switch number.(type) {
	case json.Number:
		intVal64, err := number.(json.Number).Float64()
		if err != nil {
			return nil, errors.New("Unable to convert a json number")
		}
		return big.NewInt(int64(int(intVal64))), nil
	case float64:
		return big.NewInt(int64(number.(float64))), nil
	case int:
		return big.NewInt(int64(number.(int))), nil
	case int64:
		return big.NewInt(number.(int64)), nil
	case *big.Int:
		return number.(*big.Int), nil
	case string:
		capacityInt, _ := new(big.Int).SetString(number.(string), 10)
		return capacityInt, nil
	}
	return big.NewInt(1), errors.New("Unable to convert number: " + number.(string) + " to a known type")
}
