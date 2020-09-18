// framework managed constants
var currentResources = []
var resourcePoolProperties = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

/*
IPv6 prefix allocation strategy

- Expects Ipv6 prefix resource type to have 2 properties of type int ["address:string", "mask:int"]
- userInput.desiredSize is a required parameter
- userInput.subnet is an optional parameter specifying whether allocated prefix will be used as a real subnet or just
  as an IP pool. Essentially whether to count subnet address and broadcast into the range or not. If set to true, the resulting
  prefix will have a capacity of desiredSize + 2. Defaults to false
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

// compare prefixes based on their broadcast address
function comparePrefix(prefix1, prefix2) {
    let endOfP1 = inet_aton(prefix1.address) + subnetAddresses(prefix1.prefix)
    let endOfP2 = inet_aton(prefix2.address) + subnetAddresses(prefix2.prefix)
    if (endOfP1 - endOfP2 < BigInt(0)) {
        return -1
    } else if (endOfP1 - endOfP2 > BigInt(0)) {
        return 1
    }

    return 0
}

// sum up capacity of an array of addresses
function prefixesCapacity(currentResources) {
    let width = BigInt(0)
    for (let allocatedPrefix of currentResources) {
        width += subnetAddresses(allocatedPrefix.prefix)
    }
    return width
}

// calculate utilized capacity based on previously allocated prefixes + a newly allocated prefix
function utilizedCapacity(allocatedRanges, newlyAllocatedRangeCapacity) {
    return prefixesCapacity(allocatedRanges) + newlyAllocatedRangeCapacity
}

// calculate free capacity based on previously allocated prefixes
function freeCapacity(parentPrefix, utilisedCapacity) {
    return subnetAddresses(parentPrefix.prefix) - utilisedCapacity
}

// log utilisation stats
function logStats(newlyAllocatedPrefix, parentRange, allocatedPrefixes = [], level = "log") {
    let newlyAllocatedPrefixCapacity = 0
    if (newlyAllocatedPrefix) {
        newlyAllocatedPrefixCapacity = subnetAddresses(newlyAllocatedPrefix.prefix)
    } else {
        newlyAllocatedPrefixCapacity = 0
    }

    let utilisedCapacity = utilizedCapacity(allocatedPrefixes, BigInt(newlyAllocatedPrefixCapacity))
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

function prefixesToString(currentResourcesUnwrapped) {
    let prefixesToStr = ""
    for (let allocatedPrefix of currentResourcesUnwrapped) {
        prefixesToStr += prefixToStr(allocatedPrefix)
        prefixesToStr += prefixToRangeStr(allocatedPrefix)
        prefixesToStr += ", "
    }
    return prefixesToStr
}

function prefixToStr(prefix) {
    return `${prefix.address}/${prefix.prefix}`
}

function prefixToRangeStr(prefix) {
    return `[${prefix.address}-${inet_ntoa(inet_aton(prefix.address) + subnetAddresses(prefix.prefix) - BigInt(1))}]`
}

function calculateDesiredSubnetMask() {
    let desiredSizeBigInt = BigInt(userInput.desiredSize)
    let newSubnetBits = BigInt(1)
    for (let i = 1; i <= 128; i++) {
        newSubnetBits *= BigInt(2)
        if (newSubnetBits >= desiredSizeBigInt) {
            newSubnetBits = i
            break;
        }
    }
    let newSubnetMask = 128 - newSubnetBits
    let newSubnetCapacity = subnetAddresses(newSubnetMask)
    return {newSubnetMask, newSubnetCapacity}
}

// calculate the nearest possible address for a subnet where mask === newSubnetMask
//  that is outside of allocatedSubnet
function findNextFreeSubnetAddress(allocatedSubnet, newSubnetMask) {
    // find the first address after currently iterated allocated subnet
    let nextAvailableAddressNum = inet_aton(allocatedSubnet.address) + subnetAddresses(allocatedSubnet.prefix)
    // remove any bites from the address above after newSubnetMask
    let newSubnetMaskNegative = BigInt(128) - BigInt(newSubnetMask)
    let possibleSubnetNum = (nextAvailableAddressNum >> newSubnetMaskNegative) << newSubnetMaskNegative
    // keep going until we find an address outside of currently iterated allocated subnet
    while (nextAvailableAddressNum > possibleSubnetNum) {
        possibleSubnetNum = ((possibleSubnetNum >> newSubnetMaskNegative) + BigInt(1)) << newSubnetMaskNegative
    }
    return possibleSubnetNum
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

    if (!userInput.desiredSize) {
        console.error("Unable to allocate subnet from root prefix: " + rootPrefixStr +
            ". Desired size of a new subnet size not provided as userInput.desiredSize")
        return null
    }

    // Convert desiredSize to bigint
    userInput.desiredSize = BigInt(userInput.desiredSize)

    if (userInput.desiredSize < BigInt(2)) {
        console.error("Unable to allocate subnet from root prefix: " + rootPrefixStr +
            ". Desired size is invalid: " + userInput.desiredSize + ". Use values >= 2")
        return null
    }

    if (userInput.subnet === true) {
        // reserve subnet address and broadcast
        userInput.desiredSize += BigInt(2)
    }

    // Calculate smallest possible subnet mask to fit desiredSize
    let {newSubnetMask, newSubnetCapacity} = calculateDesiredSubnetMask()

    // unwrap and sort currentResources
    let currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
    currentResourcesUnwrapped.sort(comparePrefix)

    let possibleSubnetNum = rootAddressNum
    // iterate over allocated subnets and see if a desired new subnet can be squeezed in
    for (let allocatedSubnet of currentResourcesUnwrapped) {

        let allocatedSubnetNum = inet_aton(allocatedSubnet.address)
        let chunkCapacity = allocatedSubnetNum - possibleSubnetNum
        if (chunkCapacity >= userInput.desiredSize) {
            // there is chunk with sufficient capacity between possibleSubnetNum and allocatedSubnet.address
            let newlyAllocatedPrefix = {
                "address": inet_ntoa(possibleSubnetNum),
                "prefix": newSubnetMask
            }
            // FIXME How to pass these stats ?
            // logStats(newlyAllocatedPrefix, rootPrefixParsed, currentResourcesUnwrapped)
            return newlyAllocatedPrefix
        }

        // move possible subnet start to a valid address outside of allocatedSubnet's addresses and continue the search
        possibleSubnetNum = findNextFreeSubnetAddress(allocatedSubnet, newSubnetMask)
    }

    // check if there is any space left at the end of parent range
    if (possibleSubnetNum + newSubnetCapacity <= rootAddressNum + rootCapacity) {
        // there sure is some space, use it !
        let newlyAllocatedPrefix = {
            "address": inet_ntoa(possibleSubnetNum),
            "prefix": newSubnetMask
        }
        // FIXME How to pass these stats ?
        // logStats(newlyAllocatedPrefix, rootPrefixParsed, currentResourcesUnwrapped)
        return newlyAllocatedPrefix
    }

    // no suitable range found
    console.error("Unable to allocate Ipv6 prefix from: " + rootPrefixStr +
        ". Insufficient capacity to allocate a new prefix of size: " + userInput.desiredSize)
    console.error("Currently allocated prefixes: " + prefixesToString(currentResourcesUnwrapped))
    logStats(null, rootPrefixParsed, currentResourcesUnwrapped, "error")
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

exports.invoke = invoke
exports.invokeWithParams = invokeWithParams
exports.utilizedCapacity = utilizedCapacity
exports.freeCapacity = freeCapacity
exports.parsePrefix = parsePrefix
exports.inet_ntoa = inet_ntoa
exports.inet_aton = inet_aton
// For testing purposes
