package src

import (
	"github.com/pkg/errors"
)

type Ipv4 struct {
	currentResources []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput map[string]interface{}
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
func (ipv4 *Ipv4) FreeCapacity(address string, mask int, utilisedCapacity float64) float64 {
	value, ok := ipv4.userInput["subnet"]
	var subnetItself int
	if ok && value.(bool) {
		subnetItself = 1
	} else {
		subnetItself = 0
	}
	return float64(hostsInMask(address, mask)) - utilisedCapacity + float64(subnetItself)
}

func (ipv4 *Ipv4) Capacity() (map[string]interface{}, error){
	var result = make(map[string]interface{})
	rootAddressStr, _ := ipv4.resourcePoolProperties["address"]
	rootMask, _ := ipv4.resourcePoolProperties["prefix"]
	result["freeCapacity"] = ipv4.FreeCapacity(rootAddressStr.(string), rootMask.(int), float64(len(ipv4.currentResources)))
	result["utilizedCapacity"] = float64(len(ipv4.currentResources))
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
	rootPrefixStr := prefixToStr(ipv4.resourcePoolProperties)
	rootCapacity := subnetAddresses(rootMask.(int))
	rootAddressNum, err := inetAton(rootAddressStr.(string))
	if err != nil {
		return nil, err
	}

	// unwrap and sort currentResources
	var currentResourcesUnwrapped []map[string]interface{}
	for _, element := range ipv4.currentResources {
		value, ok := element["Properties"]
		if ok {
			currentResourcesUnwrapped = append(currentResourcesUnwrapped, value.(map[string]interface{}))
		}
	}
	var currentResourcesSet []string
	for _, element := range currentResourcesUnwrapped {
		value, ok := element["address"]
		if ok {
			currentResourcesSet = append(currentResourcesSet, value.(string))
		}
	}

	var firstPossibleAddr = 0
	var lastPossibleAddr = 0

	value, ok := ipv4.userInput["subnet"]
	if ok && value == true{
		firstPossibleAddr = rootAddressNum + 1
		lastPossibleAddr = rootAddressNum + rootCapacity - 1
	} else {
		firstPossibleAddr = rootAddressNum
		lastPossibleAddr = rootAddressNum + rootCapacity
	}

	var result = make(map[string]interface{})

	value, ok = ipv4.userInput["desiredValue"]
	if ok {
		desiredValueNum, err := inetAton(value.(string))
		if err != nil {
			return nil, err
		}
		if desiredValueNum >= firstPossibleAddr && desiredValueNum < lastPossibleAddr {
			desiredIpv4Address := inetNtoa(desiredValueNum)
			if stringInSlice(desiredIpv4Address, currentResourcesSet)  {
				return nil, errors.New("Ipv4 address " + value.(string) + " was already claimed." )
			}
			result["address"] = value.(string)
			return result, nil
		} else {
			return nil, errors.New("Ipv4 address " + value.(string) + " is out of " + rootPrefixStr )
		}
	}

	for  i := firstPossibleAddr; i < lastPossibleAddr; i++ {
		if !stringInSlice(inetNtoa(i), currentResourcesSet)  {
			// FIXME How to pass these stats ?
			// logStats(inet_ntoa(i), rootPrefixParsed, userInput.subnet === true, currentResourcesUnwrapped)
			result["address"] = inetNtoa(i)
			return result, nil
		}
	}
	return nil, errors.New("Unable to allocate Ipv4 address from: " + rootPrefixStr + "." +
		"Insufficient capacity to allocate a new address.\n" +
		"Currently allocated addresses: " + addressesToStr(currentResourcesUnwrapped))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}