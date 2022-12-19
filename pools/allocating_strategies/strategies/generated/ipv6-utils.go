package src

import (
	"github.com/pkg/errors"
	"math/big"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Ipv6InetNtoa returns (string) ipv6 address from (big.Int) amount of free addresses
func Ipv6InetNtoa(addrBigInt *big.Int) string {
	step := big.NewInt(112)
	var remain = addrBigInt
	var parts []string
	for (step.Cmp(big.NewInt(0))) > 0 {
		divisor := new(big.Int).Exp(big.NewInt(2), step, nil)
		parts = append(parts, new(big.Int).Quo(remain, divisor).String())
		remain = remain.Mod(addrBigInt, divisor)
		step = step.Sub(step, big.NewInt(16))
	}

	parts = append(parts, remain.String())
	for index, part := range parts {
		number, _ := strconv.Atoi(part)
		parts[index] = strconv.FormatInt(int64(number), 16)
	}
	ip := strings.Join(parts, ":")
	ip = regexp.MustCompile(`\b:?(?:0+:?){2,}`).ReplaceAllString(ip, "::")
	return ip
}

// Ipv6InetAton returns number of assignable addresses based on address and mask
func Ipv6InetAton(addrstr string) (*big.Int, error) {
	validation := net.ParseIP(addrstr)
	a := validation.To16()
	if a == nil {
		return big.NewInt(-1), errors.New("Address: " + addrstr + " is invalid ipv6 address.")
	}

	addressCount := big.NewInt(0)
	exp := big.NewInt(0)
	parts := strings.Split(addrstr, ":")
	index := indexOf(parts, "")
	if index != -1 {
		for _, element := range parts {
			if len(element) > 4 || element > "ffff" {
				return big.NewInt(-1),
					errors.New("Address is invalid, cannot be parsed: " + addrstr + ". Contains invalid part: " + element)
			}
		}

		if index != -1 {
			for len(parts) < 8 {
				parts = append(parts[:index+1], parts[index:]...)
				parts[index] = ""
			}
		}
	}
	for i := len(parts) - 1; i >= 0; i-- {
		// The following code does the same thing as the commented but looks more complicated because of library big
		//number += BigInt(decimalFromHexNum) * (BigInt(2) ** BigInt(exp))
		//exp += BigInt(16)
		decimalFromHexNum, _ := strconv.ParseInt(parts[i], 16, 64)
		localAddressCount := new(big.Int).Exp(big.NewInt(2), exp, nil)
		localAddressCount.Mul(localAddressCount, big.NewInt(decimalFromHexNum))
		addressCount.Add(addressCount, localAddressCount)
		exp.Add(exp, big.NewInt(16))
	}
	return addressCount, nil
}

func indexOf(data []string, element string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	//if element not found
	return -1
}

// number of addresses in a subnet based on its mask
func ipv6SubnetAddresses(mask int) *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), uint(128-mask))
}

func ipv6HostsInMask(addressStr string, mask int) *big.Int {
	if mask == 128 {
		return big.NewInt(1)
	}
	if mask == 127 {
		return big.NewInt(2)
	}
	address, _ := Ipv6InetAton(addressStr)
	count := ipv6SubnetLastAddress(address, mask)
	count.Sub(count, address).Add(count, big.NewInt(1))
	return count
}

func ipv6SubnetLastAddress(subnet *big.Int, mask int) *big.Int {
	count := new(big.Int).Add(subnet, ipv6SubnetAddresses(mask))
	return count.Sub(count, big.NewInt(1))
}

var prefixRegex = "([0-9a-f:]+)/([0-9]{1,3})"

// Ipv6ParsePrefix parse prefix from a string e.g. beef::/64 into an object
func Ipv6ParsePrefix(ipv6PrefixAddressStr string) (map[string]interface{}, error) {

	regularExp := regexp.MustCompile(prefixRegex)
	match := regularExp.MatchString(ipv6PrefixAddressStr)
	if match == false {
		return nil, errors.New("Prefix: " + ipv6PrefixAddressStr + " is invalid, doesn't match regex: " + prefixRegex)
	}

	addrBigInt, err := Ipv6InetAton(regularExp.FindAllStringSubmatch(ipv6PrefixAddressStr, -1)[0][1])
	if err != nil {
		return nil, err
	}

	mask, err := strconv.Atoi(regularExp.FindAllStringSubmatch(ipv6PrefixAddressStr, -1)[0][2])
	if err != nil {
		return nil, err
	}

	if mask < 0 || mask > 128 {
		return nil, errors.New("Mask is invalid outside of ipv6 range: " + strconv.Itoa(mask))
	}

	if mask == 0 {
		addrBigInt = big.NewInt(0)
	} else {
		// making sure to nullify any bits set to 1 in subnet addr outside of mask
		// we can extract the result of this expression (addrBigInt >> big.Int(128 - mask)) << big.Int(128 - mask)
		// only in 2 steps:
		firstStep := new(big.Int).Rsh(addrBigInt, uint(128-mask))
		addrBigInt = new(big.Int).Lsh(firstStep, uint(128-mask))
	}

	var result = make(map[string]interface{})
	result["address"] = Ipv6InetNtoa(addrBigInt)
	result["prefix"] = mask
	return result, nil
}

func prefixesToStr(currentResourcesAddresses []Ipv6PrefixStruct) string {
	var addressesToStr = ""
	for _, allocatedAddr := range currentResourcesAddresses {
		value := make(map[string]interface{})
		value["address"] = allocatedAddr.address
		value["prefix"] = allocatedAddr.prefix
		addressesToStr += prefixToStr(value)
		addressesToStr += prefixToRangeStr(value)
		addressesToStr += ", "
	}
	return addressesToStr
}

func prefixToRangeStr(prefix map[string]interface{}) string {
	addressAton, _ := Ipv6InetAton(prefix["address"].(string))
	addressAton.Add(addressAton, ipv6SubnetAddresses(prefix["prefix"].(int)-1))
	addressNtoa := Ipv6InetNtoa(addressAton)
	return "[" + prefix["address"].(string) + "-" + addressNtoa + "]"
}
