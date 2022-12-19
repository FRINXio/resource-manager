// ipv4 int to str
export function inet_ntoa(addrint) {
    return ((addrint >> 24) & 0xff) + "." +
        ((addrint >> 16) & 0xff) + "." +
        ((addrint >> 8) & 0xff) + "." +
        (addrint & 0xff)
}

// ipv4 str to int
export function inet_aton(addrstr) {
    var re = /^([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})$/
    var res = re.exec(addrstr)

    if (res === null) {
        console.error("Address: " + addrstr + " is invalid, doesn't match regex: " + re)
        return null
    }

    for (var i = 1; i <= 4; i++) {
        if (res[i] < 0 || res[i] > 255) {
            console.error("Address: " + addrstr + " is invalid, outside of ipv4 range: " + addrstr)
            return null
        }
    }

    return (res[1] << 24) | (res[2] << 16) | (res[3] << 8) | res[4]
}// parse prefix from a string e.g. 1.2.3.4/18 into an object
// number of addresses in a subnet based on its mask
export function subnetAddresses(mask) {
    return 1 << (32 - mask)
}

// number of assignable addresses based on address and mask
export function hostsInMask(addressStr, mask) {
    if (mask == 32) {
        return 1;
    }
    if (mask == 31) {
        return 2;
    }
    let address = inet_aton(addressStr);

    return subnetLastAddress(address, mask) - (address + 1);
}

function subnetLastAddress(subnet, mask) {
    return subnet + subnetAddresses(mask) - 1;
}

const prefixRegex = /([0-9.]+)\/([0-9]{1,2})/

export function parsePrefix(str) {
    let res = prefixRegex.exec(str)
    if (res == null) {
        console.error("Prefix is invalid, doesn't match regex: " + prefixRegex)
        return null
    }

    let addrNum = inet_aton(res[1])
    if (addrNum == null) {
        return null
    }

    let mask = parseInt(res[2], 10)
    if (mask < 0 || mask > 32) {
        console.error("Mask is invalid outside of ipv4 range: " + mask)
        return null
    }

    if (mask === 0) {
        addrNum = 0
    } else {
        // making sure to nullify any bits set to 1 in subnet addr outside of mask
        addrNum = (addrNum >>> (32 - mask)) << (32 - mask)
    }

    return {
        "address": inet_ntoa(addrNum),
        "prefix": mask
    }
}

export function addressesToStr(currentResourcesUnwrapped) {
    let addressesToStr = ""
    for (let allocatedAddr of currentResourcesUnwrapped) {
        addressesToStr += allocatedAddr.address
        addressesToStr += ", "
    }
    return addressesToStr
}

export function prefixToStr(prefix) {
    return `${prefix.address}/${prefix.prefix}`
}
