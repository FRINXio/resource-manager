// framework managed constants
var currentResources = []
var resourcePoolProperties = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

/*
IPv6 address allocation strategy
- Expects Ipv6 prefix resource type to have 2 properties of type int ["address:string", "mask:int"]
- userInput.subnet is an optional parameter specifying whether root prefix will be used as a real subnet or just
  as an IP pool. Essentially whether to consider subnet address and broadcast when allocating addresses.
- Logs utilisation stats
- Allocates previously freed prefixes
- All addresses from parent prefix are used, including the first and last one
 */

// ipv6 int to str
function inet_ntoa(addrBigInt) {
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
}

// ipv6 str to BigInt
function inet_aton(addrstr) {
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
function subnetAddresses(mask) {
    return BigInt(1) << BigInt(128 - mask)
}

const prefixRegex = /([0-9a-f:]+)\/([0-9]{1,3})/

// parse prefix from a string e.g. beef::/64 into an object
function parsePrefix(str) {
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

function addressesToStr(currentResourcesUnwrapped) {
    let addressesToStr = ""
    for (let allocatedAddr of currentResourcesUnwrapped) {
        addressesToStr += allocatedAddr.address
        addressesToStr += ", "
    }
    return addressesToStr
}

function prefixToStr(prefix) {
    return `${prefix.address}/${prefix.prefix}`
}

// calculate utilized capacity based on previously allocated prefixes + a newly allocated prefix
function utilizedCapacity(allocatedAddresses, newlyAllocatedRangeCapacity) {
    return BigInt(allocatedAddresses.length) + BigInt(newlyAllocatedRangeCapacity)
}

// number of assignable addresses based on address and mask
function hostsInMask(addressStr, mask) {
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

// calculate free capacity based on previously allocated prefixes
function freeCapacity(parentPrefix, utilisedCapacity) {
    return subnetAddresses(parentPrefix.prefix) - BigInt(utilisedCapacity)
}

function capacity() {
    let subnetItself = userInput.subnet ? BigInt(1) : BigInt(0);
    let freeInTotal = hostsInMask(resourcePoolProperties.address, resourcePoolProperties.prefix) + subnetItself;
    return { freeCapacity: Number(freeInTotal - BigInt(currentResources.length)), utilizedCapacity: currentResources.length };
}

// log utilisation stats
function logStats(newlyAllocatedAddr, parentRange, isSubnet = false, allocatedAddresses = [], level = "log") {
    let newlyAllocatedPrefixCapacity
    if (newlyAllocatedAddr) {
        newlyAllocatedPrefixCapacity = BigInt(1)
    } else {
        newlyAllocatedPrefixCapacity = BigInt(0)
    }

    let utilisedCapacity = utilizedCapacity(allocatedAddresses, newlyAllocatedPrefixCapacity)
    if(isSubnet) {
        utilisedCapacity += BigInt(2)
    }
    let remainingCapacity = freeCapacity(parentRange, utilisedCapacity)
    let utilPercentage
    if (remainingCapacity === BigInt(0)) {
        utilPercentage = 100.0
    } else {
        let utilFloat = Number(utilisedCapacity * BigInt(1000) / subnetAddresses(parentRange.prefix)) / 1000
        utilPercentage = (utilFloat * 100)
    }

    console[level]("Remaining capacity: " + remainingCapacity)
    console[level]("Utilised capacity: " + utilisedCapacity)
    console[level](`Utilisation: ${utilPercentage.toFixed(1)}%`)
}

// main
function invoke() {
    let rootPrefixParsed = resourcePoolProperties
    if (rootPrefixParsed == null) {
        console.error("Unable to extract root prefix from pool name: " + rootPrefix)
        return null
    }
    let rootAddressStr = rootPrefixParsed.address
    let rootMask = rootPrefixParsed.prefix
    let rootPrefixStr = prefixToStr(rootPrefixParsed)
    let rootCapacity = subnetAddresses(rootMask)
    let rootAddressNum = inet_aton(rootAddressStr)

    // unwrap and sort currentResources
    let currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
    let currentResourcesSet = new Set(currentResourcesUnwrapped.map(ip => ip.address))

    let firstPossibleAddr
    let lastPossibleAddr
    if (userInput.subnet === true) {
        firstPossibleAddr = rootAddressNum + BigInt(1)
        lastPossibleAddr = rootAddressNum + rootCapacity - BigInt(1)
    } else {
        firstPossibleAddr = rootAddressNum
        lastPossibleAddr = rootAddressNum + rootCapacity
    }

    for (let i = firstPossibleAddr; i < lastPossibleAddr; i++) {
        if (!currentResourcesSet.has(inet_ntoa(i))) {
            // FIXME How to pass these stats ?
            // logStats(inet_ntoa(i), rootPrefixParsed, userInput.subnet === true, currentResourcesUnwrapped)
            return {
                "address": inet_ntoa(i)
            }
        }
    }

    // no suitable range found
    console.error("Unable to allocate Ipv6 address from: " + rootPrefixStr +
        ". Insufficient capacity to allocate a new address")
    console.error("Currently allocated addresses: " + addressesToStr(currentResourcesUnwrapped))
    logStats(null, rootPrefixParsed, userInput.subnet === true, currentResourcesUnwrapped, "error")
    return null
}

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg
    resourcePoolProperties = resourcePoolArg
    userInput = userInputArg
    return invoke()
}

function invokeWithParamsCapacity(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg
    resourcePoolProperties = resourcePoolArg
    userInput = userInputArg
    return capacity()
}

exports.invoke = invoke
exports.invokeWithParams = invokeWithParams
exports.invokeWithParamsCapacity = invokeWithParamsCapacity
exports.parsePrefix = parsePrefix
exports.utilizedCapacity = utilizedCapacity
exports.freeCapacity = freeCapacity
// For testing purposes
