package src

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

func inetNtoa(addrint int) string {
	return strconv.Itoa((addrint >> 24) & 0xff) + "." +
		strconv.Itoa((addrint >> 16) & 0xff) + "." +
		strconv.Itoa((addrint >> 8) & 0xff)+ "." +
		strconv.Itoa(addrint & 0xff)
}


func inetAton(addrstr string) (int, error) {
	re, _ := regexp.Compile("^([0-9]{1,3})\\.([0-9]{1,3})\\.([0-9]{1,3})\\.([0-9]{1,3})$")
	res := re.FindStringSubmatch(addrstr)
	if len(res) < 5 {
		return 0, errors.New("Address: " + addrstr + " is invalid, doesn't match regex: " + re.String())
	}
	for i := 1; i <= 4; i++ {
		n, err := strconv.Atoi(res[i])
		if err != nil && (n < 0 || n > 255){
			return 0, errors.New("Address: " + addrstr + " is invalid, outside of ipv4 range: " + addrstr)
		}
	}
	n1, _ := strconv.Atoi(res[1])
	n2, _ := strconv.Atoi(res[2])
	n3, _ := strconv.Atoi(res[3])
	n4, _ := strconv.Atoi(res[4])
	return (n1 << 24) | (n2 << 16) | (n3 << 8) | n4, nil
}
// parse prefix from a string e.g. 1.2.3.4/18 into an object
// number of addresses in a subnet based on its mask
func subnetAddresses(mask int) int {
	return 1 << (32 - mask)
}

func hostsInMask(addressStr string, mask int) int {
	if mask == 32 {
		return 1
	}
	if mask == 31 {
		return 2
	}
	address, _ := inetAton(addressStr)

	return subnetLastAddress(address, mask) - (address + 1)
}

func subnetLastAddress(subnet int, mask int) int {
	return subnet + subnetAddresses(mask) - 1
}

func addressesToStr(currentResourcesUnwrapped []map[string]interface{}) string{
	var addressesToStr= ""
	for _, allocatedAddr := range currentResourcesUnwrapped {
		addressStr, _ := allocatedAddr["address"]
		addressesToStr += addressStr.(string)
		addressesToStr += ", "
	}
	return addressesToStr
}

func prefixToStr(prefix map[string]interface{}) string {
	addressStr, _ := prefix["address"]
	prefixStr, _ := prefix["prefix"]
	return addressStr.(string) + "/" + strconv.Itoa(prefixStr.(int))
}