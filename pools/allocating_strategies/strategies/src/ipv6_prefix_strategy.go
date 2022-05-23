package src

import (
	"github.com/pkg/errors"
	"math/big"
	"sort"
	"strconv"
)

type Ipv6Prefix struct {
	currentResources       []map[string]interface{}
	resourcePoolProperties map[string]interface{}
	userInput              map[string]interface{}
}

type Ipv6PrefixStruct struct {
	address string
	prefix  int
}

func NewIpv6Prefix(currentResources []map[string]interface{},
	resourcePoolProperties map[string]interface{},
	userInput map[string]interface{}) Ipv6Prefix {
	return Ipv6Prefix{currentResources, resourcePoolProperties, userInput}
}

func (ipv6Prefix *Ipv6Prefix) UtilizedCapacity(allocatedRanges []map[string]interface{}, newlyAllocatedRangeCapacity float64) float64 {
	return float64(len(allocatedRanges)) + newlyAllocatedRangeCapacity
}

// FreeCapacity calculate free capacity based on previously allocated prefixes
func (ipv6Prefix *Ipv6Prefix) FreeCapacity(parentPrefix string, utilisedCapacity float64) float64 {
	capacityString := strconv.FormatFloat(utilisedCapacity, 'f', -1, 64)
	capacityInt, _ := new(big.Int).SetString(capacityString, 10)
	parentPrefixInt, _ := strconv.Atoi(parentPrefix)
	addressesCount := ipv6SubnetAddresses(parentPrefixInt)
	addressesCount.Sub(addressesCount, capacityInt)
	return float64(addressesCount.Int64())
}

func (ipv6Prefix *Ipv6Prefix) Capacity() (map[string]interface{}, error) {
	var result = make(map[string]interface{})
	rootAddressStr, ok := ipv6Prefix.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv6Prefix.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	subnetItself := new(big.Int)
	i, ok := ipv6Prefix.userInput["subnet"]
	if i != nil {
		subnetItself = big.NewInt(1)
	} else {
		subnetItself = big.NewInt(0)
	}

	var allocatedCapacity = big.NewInt(0)
	for _, resource := range ipv6Prefix.currentResources {
		prefix, err := NumberToInt(resource["Properties"].(map[string]interface{})["prefix"].(interface{}))
		if err != nil {
			return nil, err
		}
		allocatedCapacity.Add(allocatedCapacity, ipv6HostsInMask(
			resource["Properties"].(map[string]interface{})["address"].(string),
			prefix.(int)))
	}
	rootMask, err := NumberToInt(rootMask)
	if err != nil {
	    return nil, err
	}
	totalCapacity := ipv6HostsInMask(rootAddressStr.(string), rootMask.(int))
	totalCapacity.Sub(totalCapacity, allocatedCapacity)
	totalCapacity.Add(totalCapacity, subnetItself)

	result["freeCapacity"] = totalCapacity.String()
	result["utilizedCapacity"] = allocatedCapacity.String()
	return result, nil
}

func (ipv6Prefix *Ipv6Prefix) Invoke() (map[string]interface{}, error) {
	if ipv6Prefix.resourcePoolProperties == nil {
		return nil, errors.New("Unable to extract resources")
	}
	rootAddressStr, ok := ipv6Prefix.resourcePoolProperties["address"]
	if !ok {
		return nil, errors.New("Unable to extract address resource")
	}
	rootMask, ok := ipv6Prefix.resourcePoolProperties["prefix"]
	if !ok {
		return nil, errors.New("Unable to extract prefix resources")
	}
	rootMask, err := NumberToInt(rootMask)
	if err != nil {
		return nil, err
	}
	rootPrefixStr := prefixToStr(ipv6Prefix.resourcePoolProperties)
	rootCapacity := ipv6SubnetAddresses(rootMask.(int))
	rootAddressNum, err := Ipv6InetAton(rootAddressStr.(string))
	if err != nil {
		return nil, err
	}

	value, ok := ipv6Prefix.userInput["desiredSize"]
	if !ok {
		return nil, errors.New("Unable to allocate subnet from root prefix: " + rootPrefixStr +
			". Desired size of a new subnet size not provided as userInput.desiredSize")
	}
	desiredSize, err := NumberToBigInt(value)
	if err != nil {
		return nil, err
	}
	if desiredSize.Cmp(big.NewInt(2)) < 0 {
		return nil, errors.New("Unable to allocate subnet from root prefix: " + rootPrefixStr +
			". Desired size is invalid: " + desiredSize.String() + ". Use values >= 2")
	}

	if value, ok := ipv6Prefix.userInput["subnet"]; ok && value == true {
		// reserve subnet address and broadcast
		desiredSize.Add(desiredSize, big.NewInt(2))
	}

	// Calculate smallest possible subnet mask to fit desiredSize
	newSubnetMask, newSubnetCapacity := calculateDesiredSubnetMask(desiredSize)

	var currentResourcesStruct []Ipv6PrefixStruct
	for _, resource := range ipv6Prefix.currentResources {
		address, prefix, err := getAddressAndPrefixFromCurrentResource(resource)
		if err != nil {
			return nil, err
		}
		currentResourcesStruct = append(currentResourcesStruct, Ipv6PrefixStruct{address: address, prefix: prefix})
	}

	sort.Slice(currentResourcesStruct, func(i, j int) bool {
		address1, _ := Ipv6InetAton(currentResourcesStruct[i].address)
		address2, _ := Ipv6InetAton(currentResourcesStruct[j].address)
		endOfP1 := new(big.Int).Add(address1, ipv6SubnetAddresses(currentResourcesStruct[i].prefix))
		endOfP2 := new(big.Int).Add(address2, ipv6SubnetAddresses(currentResourcesStruct[j].prefix))
		result := endOfP1.Cmp(endOfP2)
		if result >= 0 {
			return false
		}
		return true
	})

	possibleSubnetNum := rootAddressNum
	// iterate over allocated subnets and see if a desired new subnet can be squeezed in
	for _, currentResource := range currentResourcesStruct {

		allocatedSubnetNum, err := Ipv6InetAton(currentResource.address)
		if err != nil {
			return nil, errors.New("Wrong property address: " + currentResource.address + " in current resources")
		}
		chunkCapacity := new(big.Int).Sub(allocatedSubnetNum, possibleSubnetNum)
		if chunkCapacity.Cmp(desiredSize) >= 0 {
			// there is chunk with sufficient capacity between possibleSubnetNum and allocatedSubnet.address
			var newlyAllocatedPrefix = make(map[string]interface{})
			newlyAllocatedPrefix["address"] = Ipv6InetNtoa(possibleSubnetNum)
			newlyAllocatedPrefix["prefix"] = newSubnetMask

			return newlyAllocatedPrefix, nil
		}

		// move possible subnet start to a valid address outside of allocatedSubnet's addresses and continue the search
		possibleSubnetNum, err = findNextFreeSubnetAddress(currentResource, newSubnetMask)
		if err != nil {
			return nil, err
		}
	}

	// check if there is any space left at the end of parent range
	currentAmount := new(big.Int).Add(possibleSubnetNum, newSubnetCapacity)
	rootAmount := new(big.Int).Add(rootAddressNum, rootCapacity)
	if currentAmount.Cmp(rootAmount) < 1 {
		// there sure is some space, use it !
		var newlyAllocatedPrefix = make(map[string]interface{})
		newlyAllocatedPrefix["address"] = Ipv6InetNtoa(possibleSubnetNum)
		newlyAllocatedPrefix["prefix"] = newSubnetMask

		return newlyAllocatedPrefix, nil
	}
	// no suitable range found
	return nil, errors.New("Unable to allocate Ipv6 prefix from: " + rootPrefixStr +
		". Insufficient capacity to allocate a new prefix of size: " + desiredSize.String() +
		". Currently allocated prefixes: " + prefixesToStr(currentResourcesStruct))
}

func calculateDesiredSubnetMask(desiredSize *big.Int) (int, *big.Int) {
	newSubnetBits := big.NewInt(1)
	newSubnet := 1
	for i := 1; i <= 128; i++ {
		newSubnetBits.Mul(newSubnetBits, big.NewInt(2))
		if newSubnetBits.Cmp(desiredSize) >= 0 {
			newSubnet = i
			break
		}
	}
	newSubnetMask := 128 - newSubnet
	newSubnetCapacity := ipv6SubnetAddresses(newSubnetMask)
	return newSubnetMask, newSubnetCapacity
}

// calculate the nearest possible address for a subnet where mask === newSubnetMask
//  that is outside allocatedSubnet
func findNextFreeSubnetAddress(allocatedSubnet Ipv6PrefixStruct, newSubnetMask int) (*big.Int, error) {
	address, err := Ipv6InetAton(allocatedSubnet.address)
	if err != nil {
		return nil, err
	}
	// find the first address after currently iterated allocated subnet
	nextAvailableAddressNum := new(big.Int).Add(address, ipv6SubnetAddresses(allocatedSubnet.prefix))
	// remove any bites from the address above after newSubnetMask
	newSubnetMaskNegative := new(big.Int).Sub(big.NewInt(128), big.NewInt(int64(newSubnetMask)))
	possibleSubnetNum := new(big.Int).Rsh(nextAvailableAddressNum, uint(newSubnetMaskNegative.Int64()))
	possibleSubnetNum.Lsh(possibleSubnetNum, uint(newSubnetMaskNegative.Int64()))
	// keep going until we find an address outside of currently iterated allocated subnet
	for nextAvailableAddressNum.Cmp(possibleSubnetNum) > 0 {
		possibleSubnetNum.Rsh(possibleSubnetNum, uint(newSubnetMaskNegative.Int64()))
		possibleSubnetNum.Add(possibleSubnetNum, big.NewInt(1))
		possibleSubnetNum.Lsh(possibleSubnetNum, uint(newSubnetMaskNegative.Int64()))
	}
	return possibleSubnetNum, nil
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
