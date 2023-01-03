package src

import (
	"github.com/pkg/errors"
	"math"
	"sort"
	"strconv"
)

type Ipv4Struct struct {
	address string
	prefix  int
}

type Ipv4Prefix struct {
	currentResources       []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

func NewIpv4Prefix(currentResources []map[string]interface{},
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) Ipv4Prefix {
	return Ipv4Prefix{currentResources, resourcePoolProperties, userInput}
}

func (ipv4prefix *Ipv4Prefix) UtilizedCapacity(
	allocatedRanges []map[string]interface{},
	newlyAllocatedRangeCapacity int) int {
	return prefixesCapacity(allocatedRanges) + newlyAllocatedRangeCapacity
}

// FreeCapacity calculate free capacity based on previously allocated prefixes
func (ipv4prefix *Ipv4Prefix) FreeCapacity(
	parentPrefix map[string]interface{},
	utilisedCapacity int) int {
	prefix, _ := parentPrefix["prefix"]
	return subnetAddresses(prefix.(int)) - utilisedCapacity
}

func getAddressAndPrefixFromCurrentResource(currentResource map[string]interface{}) (string, int, error) {
	value, ok := currentResource["Properties"]
	if !ok {
		return "", 0, errors.New("Unable to extract properties from resource")
	}
	address, ok := value.(map[string]interface{})["address"]
	if !ok {
		return "", 0, errors.New("Unable to extract address resource from properties")
	}
	prefix, ok := value.(map[string]interface{})["prefix"]
	if !ok {
		return "", 0, errors.New("Unable to extract prefix resource from properties")
	}
	prefix, err := NumberToInt(prefix)
	if err != nil {
		return "", 0, err
	}
	return address.(string), prefix.(int), nil
}

func (ipv4prefix *Ipv4Prefix) Capacity() (map[string]interface{}, error) {
	if ipv4prefix.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	rootAddressStr, ok := ipv4prefix.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv4prefix.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	totalCapacity := hostsInMask(rootAddressStr.(string), rootMask.(int)) + 2
	allocatedCapacity := 0
	var subnetItself = 0
	if subnet, ok := ipv4prefix.resourcePoolProperties["subnet"].(bool); ok && subnet == true {
		subnetItself = 2
	}

	for _, resource := range ipv4prefix.currentResources {
		address, prefix, err := getAddressAndPrefixFromCurrentResource(resource)
		if err != nil {
			return nil, err
		}
		allocatedCapacity += hostsInMask(address, prefix) + subnetItself
	}
	var result = make(map[string]interface{})
	result["freeCapacity"] = strconv.FormatFloat(float64(totalCapacity-allocatedCapacity), 'g', 30, 64)
	result["utilizedCapacity"] = strconv.Itoa(allocatedCapacity)

	return result, nil
}

func prefixesCapacity(currentResources []map[string]interface{}) int {
	width := 0
	for _, allocatedPrefix := range currentResources {
		value, _ := allocatedPrefix["prefix"]
		width += subnetAddresses(value.(int))
	}
	return width
}

// calculate the nearest possible address for a subnet where mask === newSubnetMask
//
//	that is outside of allocatedSubnet
func findNextFreeSubnetAddress(allocatedSubnet Ipv4Struct, newSubnetMask int) int {
	// find the first address after currently iterated allocated subnet
	addressNumber, _ := inetAton(allocatedSubnet.address)
	nextAvailableAddressNum := addressNumber + subnetAddresses(allocatedSubnet.prefix)
	// remove any bites from the address above after newSubnetMask
	newSubnetMaskNegative := 32 - newSubnetMask
	possibleSubnetNum := int(uint(nextAvailableAddressNum)>>newSubnetMaskNegative) << newSubnetMaskNegative
	// keep going until we find an address outside of currently iterated allocated subnet
	for nextAvailableAddressNum > int(possibleSubnetNum) {
		possibleSubnetNum = (int(uint(possibleSubnetNum)>>newSubnetMaskNegative) + 1) << newSubnetMaskNegative
	}
	return possibleSubnetNum
}

func (ipv4prefix *Ipv4Prefix) calculateDesiredSubnetMask() (int, int) {
	desiredSize, _ := ipv4prefix.userInput["desiredSize"]
	desiredSize, _ = NumberToInt(desiredSize)
	newSubnetBits := math.Ceil(math.Log(float64(desiredSize.(int))) / math.Log(2))
	newSubnetMask := 32 - newSubnetBits
	newSubnetCapacity := subnetAddresses(int(newSubnetMask))
	return int(newSubnetMask), newSubnetCapacity
}

func (ipv4prefix *Ipv4Prefix) Invoke() (map[string]interface{}, error) {
	if ipv4prefix.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	rootAddressStr, ok := ipv4prefix.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv4prefix.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	isSubnet, ok := ipv4prefix.resourcePoolProperties["subnet"]
	if !ok {
		return nil, errors.New("Unable to extract subnet resources")
	}
	rootMask, err := NumberToInt(rootMask)
	if err != nil {
		return nil, err
	}
	rootPrefixStr := prefixToStr(ipv4prefix.resourcePoolProperties)
	rootCapacity := subnetAddresses(rootMask.(int))
	rootAddressNum, err := inetAton(rootAddressStr.(string))
	if err != nil {
		return nil, err
	}

	value, ok := ipv4prefix.userInput["desiredSize"]
	if !ok {
		return nil, errors.New("Unable to allocate subnet from root prefix: " + rootPrefixStr +
			". Desired size of a new subnet size not provided as userInput.desiredSize")
	}
	desiredSize, err := NumberToInt(value)
	desiredValue, isDesiredValueOk := ipv4prefix.userInput["desiredValue"]

	if !isDesiredValueOk {
		desiredValue = nil
	}

	if err != nil {
		return nil, err
	}
	if desiredSize.(int) < 0 {
		return nil, errors.New("Unable to allocate subnet from root prefix: " + rootPrefixStr +
			". Desired size is invalid: " + strconv.Itoa(desiredSize.(int)) + ". Use values >= 2")
	}
	newSubnetMask, newSubnetCapacity := ipv4prefix.calculateDesiredSubnetMask()

	networkAddresses, err := networkAddressesInSubnet(rootAddressStr.(string), rootMask.(int), newSubnetCapacity, newSubnetMask)

	if err != nil && desiredValue != nil {
		return nil, errors.New("We weren't able to handle formatting of provided inputs, that were in incorrect shape")
	}

	if desiredValue != nil {
		if isHostAddressValid(networkAddresses, desiredValue.(string)) == false {
			return nil, errors.New("You provided invalid host address.")
		}
	}

	var currentResourcesStruct []Ipv4Struct
	for _, resource := range ipv4prefix.currentResources {
		address, prefix, err := getAddressAndPrefixFromCurrentResource(resource)
		if err != nil {
			return nil, err
		}
		currentResourcesStruct = append(currentResourcesStruct, Ipv4Struct{address: address, prefix: prefix})
	}

	// compare prefixes based on their broadcast address
	sort.Slice(currentResourcesStruct, func(i, j int) bool {
		address1Num, _ := inetAton(currentResourcesStruct[i].address)
		address2Num, _ := inetAton(currentResourcesStruct[j].address)
		endOfP1 := address1Num + subnetAddresses(currentResourcesStruct[i].prefix)
		endOfP2 := address2Num + subnetAddresses(currentResourcesStruct[j].prefix)
		return endOfP1 < endOfP2
	})

	possibleSubnetNum := rootAddressNum
	var result = make(map[string]interface{})

	// iterate over allocated subnets and see if a desired new subnet can be squeezed in
	for _, allocatedSubnet := range currentResourcesStruct {
		allocatedSubnetNum, _ := inetAton(allocatedSubnet.address)
		chunkCapacity := allocatedSubnetNum - possibleSubnetNum

		if desiredValue != nil {
			desiredValueNum, _ := inetAton(desiredValue.(string))
			for possibleSubnetNum <= desiredValueNum {
				if chunkCapacity >= desiredSize.(int) && desiredValueNum == possibleSubnetNum {
					// there is chunk with sufficient capacity between possibleSubnetNum and allocatedSubnet.address
					result["address"] = inetNtoa(possibleSubnetNum)
					result["prefix"] = newSubnetMask
					result["subnet"] = isSubnet

					return result, nil
				}

				chunkCapacity -= newSubnetCapacity
				possibleSubnetNum += newSubnetCapacity
			}
		} else {
			if chunkCapacity >= desiredSize.(int) {
				// there is chunk with sufficient capacity between possibleSubnetNum and allocatedSubnet.address
				result["address"] = inetNtoa(possibleSubnetNum)
				result["prefix"] = newSubnetMask
				result["subnet"] = isSubnet

				return result, nil
			}
		}

		// move possible subnet start to a valid address outside allocatedSubnet's addresses and continue the search
		possibleSubnetNum = findNextFreeSubnetAddress(allocatedSubnet, newSubnetMask)
	}

	if desiredValue != nil {
		desiredValueNum, _ := inetAton(desiredValue.(string))
		for possibleSubnetNum+newSubnetCapacity <= rootAddressNum+rootCapacity {
			var hasFreeSpaceInRange = possibleSubnetNum+newSubnetCapacity <= rootAddressNum+rootCapacity

			if hasFreeSpaceInRange && possibleSubnetNum == desiredValueNum {
				result["address"] = inetNtoa(possibleSubnetNum)
				result["prefix"] = newSubnetMask
				result["subnet"] = isSubnet
				return result, nil
			}

			possibleSubnetNum += newSubnetCapacity
		}
	} else {
		// check if there is any space left at the end of parent range
		if possibleSubnetNum+newSubnetCapacity <= rootAddressNum+rootCapacity {
			// there sure is some space, use it !
			result["address"] = inetNtoa(possibleSubnetNum)
			result["prefix"] = newSubnetMask
			result["subnet"] = isSubnet
			return result, nil
		}
	}

	return nil, errors.New("Unable to allocate Ipv4 prefix from: " + rootPrefixStr + ". " +
		"Insufficient capacity to allocate a new prefix of size: " + strconv.Itoa(desiredSize.(int)) + "\n" +
		"Currently allocated addresses: " + addressesToStr(ipv4prefix.currentResources))
}
