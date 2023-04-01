package src

import (
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"math"
	"net"
	"reflect"
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

	isSubnet, isSubnetOk := ipv4prefix.resourcePoolProperties["subnet"].(bool)
	totalCapacity := hostsInMask(rootAddressStr.(string), rootMask.(int))
	allocatedCapacity := 0

	if isSubnetOk && isSubnet && (rootMask.(int) == 31 || rootMask.(int) == 32) {
		var result = make(map[string]interface{})
		result["freeCapacity"] = strconv.Itoa(0)
		result["utilizedCapacity"] = strconv.Itoa(0)

		return result, nil
	}

	if isSubnetOk && isSubnet == false && rootMask.(int) != 31 && rootMask.(int) != 32 {
		// we are adding 2 because total capacity is not only host addresses but also network and broadcast addresses
		// 31 and 32 prefixes return constant value, so there is not needed this calculation
		totalCapacity += 2
	}

	for _, resource := range ipv4prefix.currentResources {
		address, prefix, err := getAddressAndPrefixFromCurrentResource(resource)
		if err != nil {
			return nil, err
		}

		if prefix == 31 || prefix == 32 {
			allocatedCapacity += hostsInMask(address, prefix)
		} else {
			allocatedCapacity += hostsInMask(address, prefix) + 2
		}
	}

	var freeCapacity float64
	// we are handling the edge case when allocated capacity is same as total capacity,
	// then we need to add network and broadcast address (+2) of the main subnet
	if totalCapacity+2 == allocatedCapacity {
		freeCapacity = 0
	} else {
		freeCapacity = float64(totalCapacity - allocatedCapacity)
	}

	var result = make(map[string]interface{})
	result["freeCapacity"] = strconv.FormatFloat(freeCapacity, 'g', 30, 64)
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
	addressNumber, _ := InetAton(allocatedSubnet.address)
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

func calculateDesiredSubnetIpv4Mask(desiredSize int) (int, int) {
	newSubnetBits := math.Ceil(math.Log(float64(desiredSize)) / math.Log(2))
	newSubnetMask := 32 - newSubnetBits
	newSubnetCapacity := subnetAddresses(int(newSubnetMask))
	return int(newSubnetMask), newSubnetCapacity
}

func isIPv4AddrNetwork(addr string, prefix int) (bool, string, error) {
	_, ipNet, ipErr := net.ParseCIDR(fmt.Sprintf("%s/%d", addr, prefix))

	if ipErr != nil {
		return false, ipNet.IP.String(), ipErr
	}

	if addr == ipNet.IP.String() {
		return true, ipNet.IP.String(), nil
	}

	return false, ipNet.IP.String(), nil
}

func nextFreeNetworkAddressAfter(networkAddress string, prefix int, capacity int, allocatedResources []Ipv4Struct) string {
	networkAddressNum, _ := InetAton(networkAddress)
	var possibleSubnetNum = networkAddressNum

	for _, allocatedSubnet := range allocatedResources {
		allocatedSubnetNum, _ := InetAton(allocatedSubnet.address)
		chunkCapacity := allocatedSubnetNum - possibleSubnetNum

		fmt.Println(inetNtoa(possibleSubnetNum))

		if chunkCapacity >= capacity && allocatedSubnetNum > networkAddressNum {
			return inetNtoa(possibleSubnetNum)
		}

		// move possible subnet start to a valid address outside allocatedSubnet's addresses and continue the search
		possibleSubnetNum = findNextFreeSubnetAddress(allocatedSubnet, prefix)
	}

	if networkAddressNum > possibleSubnetNum {
		return networkAddress
	} else {
		_, netAddr, _ := isIPv4AddrNetwork(inetNtoa(possibleSubnetNum), prefix)
		return netAddr
	}
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
	rootAddressNum, err := InetAton(rootAddressStr.(string))
	if err != nil {
		return nil, err
	}

	value, ok := ipv4prefix.userInput["desiredSize"]
	if !ok {
		return nil, errors.New("Unable to allocate subnet from root prefix: " + rootPrefixStr +
			". Desired size of a new subnet size not provided as userInput.desiredSize")
	}

	var desiredSize interface{}
	var desiredSizeErr error
	if isSubnet.(bool) {
		desiredSize, desiredSizeErr = NumberToInt(value)
		if desiredSizeErr == nil {
			desSize, ok := desiredSize.(int)

			if !ok {
				return nil, gqlerror.Errorf("Unable to claim resource: usrInput.desiredSize was sent in bad format. Required format is int or string and was: %s", reflect.TypeOf(desiredSize))
			}

			desiredSize = desSize + 2
		} else {
			return nil, gqlerror.Errorf("Unable to claim resource: %s", desiredSizeErr)
		}
	} else {
		desiredSize, desiredSizeErr = NumberToInt(value)

		if desiredSizeErr != nil {
			return nil, gqlerror.Errorf("Unable to claim resource: %s", desiredSizeErr)
		}
	}
	desiredValue, isDesiredValueOk := ipv4prefix.userInput["desiredValue"]

	if !isDesiredValueOk {
		desiredValue = nil
	}

	if desiredSizeErr != nil || err != nil {
		return nil, err
	}
	if desiredSize.(int) < 0 {
		return nil, errors.New("Unable to allocate subnet from root prefix: " + rootPrefixStr +
			". Desired size is invalid: " + strconv.Itoa(desiredSize.(int)) + ". Use values >= 2")
	}
	newSubnetMask, newSubnetCapacity := calculateDesiredSubnetIpv4Mask(desiredSize.(int))

	if isSubnet.(bool) && (newSubnetMask == 31 || newSubnetMask == 32) {
		return nil, errors.Errorf("It is not possible to allocate resource with prefix %d, together with subnet set as true", newSubnetMask)
	}

	if isSubnet.(bool) && rootCapacity < desiredSize.(int) {
		return nil, errors.New("Unable to allocate Ipv4 prefix from: " + rootPrefixStr + ". " +
			"Insufficient capacity to allocate a new prefix of size: " + strconv.Itoa(desiredSize.(int)-2) + "\n" +
			"Currently allocated addresses: " + addressesToStr(ipv4prefix.currentResources))
	}

	if err != nil && desiredValue != nil {
		return nil, errors.New("We weren't able to handle formatting of provided inputs, that were of incorrect shape")
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
		address1Num, _ := InetAton(currentResourcesStruct[i].address)
		address2Num, _ := InetAton(currentResourcesStruct[j].address)
		endOfP1 := address1Num + subnetAddresses(currentResourcesStruct[i].prefix)
		endOfP2 := address2Num + subnetAddresses(currentResourcesStruct[j].prefix)
		return endOfP1 < endOfP2
	})

	possibleSubnetNum := rootAddressNum
	var result = make(map[string]interface{})

	if desiredValue != nil {
		isNetwork, networkAddr, ipErr := isIPv4AddrNetwork(desiredValue.(string), newSubnetMask)

		if ipErr != nil {
			return nil, ipErr
		}

		if !isNetwork {
			nextFreeNetworkAddress := nextFreeNetworkAddressAfter(networkAddr, newSubnetMask, newSubnetCapacity, currentResourcesStruct)
			return nil, errors.Errorf("You provided invalid network address. Network address should be %s", nextFreeNetworkAddress)
		}

		if len(currentResourcesStruct) > 0 {
			desiredValueNum, er := InetAton(desiredValue.(string))
			lastResource := currentResourcesStruct[len(currentResourcesStruct)-1]
			lastResourceSubnetNum, e := InetAton(lastResource.address)

			if er != nil || e != nil {
				return nil, errors.New("We weren't able to handle formatting of provided inputs, that were of incorrect format")
			}

			if lastResourceBroadcastNum := subnetLastAddress(lastResourceSubnetNum, lastResource.prefix); desiredValueNum > lastResourceBroadcastNum {
				if desiredValueNum+newSubnetCapacity <= rootAddressNum+rootCapacity {
					// there is chunk with sufficient capacity between possibleSubnetNum and allocatedSubnet.address
					result["address"] = desiredValue
					result["prefix"] = newSubnetMask
					result["subnet"] = isSubnet

					return result, nil
				}
			}
		} else {
			desiredValueNum, e := InetAton(desiredValue.(string))

			if e != nil {
				return nil, errors.New("We weren't able to handle formatting of provided inputs, that were of incorrect format")
			}

			if desiredValueNum+newSubnetCapacity <= rootAddressNum+rootCapacity {
				// there is chunk with sufficient capacity between possibleSubnetNum and allocatedSubnet.address
				result["address"] = desiredValue
				result["prefix"] = newSubnetMask
				result["subnet"] = isSubnet

				return result, nil
			}

			return nil, errors.New("Unable to allocate Ipv4 prefix from: " + rootPrefixStr + ". " +
				"Insufficient capacity to allocate a new prefix of size: " + strconv.Itoa(desiredSize.(int)) +
				" and with desiredValue of: " + desiredValue.(string) + "\n" +
				"Currently allocated addresses: " + addressesToStr(ipv4prefix.currentResources))
		}
	}

	// iterate over allocated subnets and see if a desired new subnet can be squeezed in
	for _, allocatedSubnet := range currentResourcesStruct {
		allocatedSubnetNum, _ := InetAton(allocatedSubnet.address)
		chunkCapacity := allocatedSubnetNum - possibleSubnetNum

		if desiredValue != nil {
			desiredValueNum, _ := InetAton(desiredValue.(string))
			desiredBroadcastNum, _ := InetAton(inetNtoa(subnetLastAddress(desiredValueNum, newSubnetMask)))

			if chunkCapacity >= newSubnetCapacity && desiredBroadcastNum < allocatedSubnetNum {
				// there is chunk with sufficient capacity between possibleSubnetNum and allocatedSubnet.address
				result["address"] = inetNtoa(desiredValueNum)
				result["prefix"] = newSubnetMask
				result["subnet"] = isSubnet

				return result, nil
			}
		} else {
			if chunkCapacity >= newSubnetCapacity {
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

	// check if there is any space left at the end of parent range
	if desiredValue == nil && possibleSubnetNum+newSubnetCapacity <= rootAddressNum+rootCapacity {
		// there sure is some space, use it !
		result["address"] = inetNtoa(possibleSubnetNum)
		result["prefix"] = newSubnetMask
		result["subnet"] = isSubnet
		return result, nil
	}

	var desiredSizeStr string
	if isSubnet.(bool) {
		desiredSizeStr = strconv.Itoa(desiredSize.(int) - 2)
	} else {
		desiredSizeStr = strconv.Itoa(desiredSize.(int))
	}

	return nil, errors.New("Unable to allocate Ipv4 prefix from: " + rootPrefixStr + ". " +
		"Insufficient capacity to allocate a new prefix of size: " + desiredSizeStr + "\n" +
		"Currently allocated addresses: " + addressesToStr(ipv4prefix.currentResources))
}
