
// ipv6 int to str
export function inet_ntoa(addrBigInt) {
    try {
        let step = BigInt(112)
        let remain = addrBigInt
        const parts = []

        while (step > BigInt(0)) {
            const divisor = BigInt(2) ** BigInt(step)
            parts.push(remain / divisor)
            remain = addrBigInt % divisor
            step -= BigInt(16)
        }
        parts.push(remain)

        let ip = parts.map(n => Number(n).toString(16)).join(":")
        return ip.replace(/\b:?(?:0+:?){2,}/, "::")
    } catch (error) {
        console.error("Address is invalid, cannot be serialized from: " + addrBigInt)
        return null
    }
}// number of assignable addresses based on address and mask
// ipv6 str to BigInt
export function inet_aton(addrstr) {
    let number = BigInt(0)
    let exp = BigInt(0)

    try {
        const parts = addrstr.split(":")
        const index = parts.indexOf("")

        for (const part of parts) {
            if (part.length > 4 || part > "ffff") {
                console.error("Address is invalid, cannot be parsed: " + addrstr + ". Contains invalid part: " + part)
                return null
            }
        }

        if (index !== -1) {
            while (parts.length < 8) {
                parts.splice(index, 0, "")
            }
        }

        for (const n of parts.map(part => part ? `0x${part}` : `0`).map(Number).reverse()) {
            number += BigInt(n) * (BigInt(2) ** BigInt(exp))
            exp += BigInt(16)
        }

        return number
    } catch (error) {
        console.error("Address is invalid, cannot be parsed: " + addrstr)
        return null
    }
}

// number of addresses in a subnet based on its mask
export function subnetAddresses(mask) {
    return BigInt(1) << BigInt(128 - mask)
}

const prefixRegex = /([0-9a-f:]+)\/([0-9]{1,3})/

// parse prefix from a string e.g. beef::/64 into an object
export function parsePrefix(str) {
    let res = prefixRegex.exec(str)
    if (res == null) {
        console.error("Prefix is invalid, doesn't match regex: " + prefixRegex)
        return null
    }

    let addrBigInt = inet_aton(res[1])
    if (addrBigInt == null) {
        return null
    }

    let mask = parseInt(res[2], 10)
    if (mask < 0 || mask > 128) {
        console.error("Mask is invalid outside of ipv6 range: " + mask)
        return null
    }

    if (mask === 0) {
        addrBigInt = 0
    } else {
        // making sure to nullify any bits set to 1 in subnet addr outside of mask
        addrBigInt = (addrBigInt >> BigInt(128 - mask)) << BigInt(128 - mask)
    }

    return {
        "address": inet_ntoa(addrBigInt),
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

export function hostsInMask(addressStr, mask) {
    if (mask == 128) {
        return 1;
    }
    if (mask == 127) {
        return 2;
    }

    let address = inet_aton(addressStr);

    return subnetLastAddress(address, mask) - BigInt(address) + BigInt(1);
}

function subnetLastAddress(subnet, mask) {
    return BigInt(subnet) + subnetAddresses(mask) - BigInt(1);
}
