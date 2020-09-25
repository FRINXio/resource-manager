// Code generated, DO NOT EDIT.

package pools

const IPV4_PREFIX = `

/*
IPv4 prefix allocation strategy

- Expects IPv4 prefix resource type to have 2 properties of type int ["address:string", "mask:int"]
- userInput.desiredSize is a required parameter e.g. desiredSize == 10  ---produces-prefix-of--->  192.168.1.0/28
- userInput.subnet is an optional parameter specifying whether allocated prefix will be used as a real subnet or just
  as an IP pool. Essentially whether to count subnet address and broadcast into the range or not. If set to true, the resulting
  prefix will have a capacity of desiredSize + 2. Defaults to false
- Logs utilisation stats
- Allocates previously freed prefixes
- All addresses from parent prefix are used, including the first and last one
 */

// ipv4 int to str
function inet_ntoa(addrint) {
    return ((addrint >> 24) & 0xff) + "." +
        ((addrint >> 16) & 0xff) + "." +
        ((addrint >> 8) & 0xff) + "." +
        (addrint & 0xff)
}

// ipv4 str to int
function inet_aton(addrstr) {
    var re = /^([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})$/
    var res = re.exec(addrstr)

    if (res === null) {
        console.error("Address is invalid, doesn't match regex: " + re)
        return null
    }

    for (var i = 1; i <= 4; i++) {
        if (res[i] < 0 || res[i] > 255) {
            console.error("Address is invalid, outside of ipv4 range: " + addrstr)
            return null
        }
    }

    return (res[1] << 24) | (res[2] << 16) | (res[3] << 8) | res[4]
}

// number of addresses in a subnet based on its mask
function subnetAddresses(mask) {
    return 1<<(32-mask)
}

// TODO code reuse between strategies

const prefixRegex = /([0-9.]+)\/([0-9]{1,2})/

// parse prefix from a string e.g. 1.2.3.4/18 into an object
function parsePrefix(str) {
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
        addrNum = (addrNum >>> (32-mask)) << (32-mask)
    }

    return {
        "address": inet_ntoa(addrNum),
        "prefix": mask
    }
}

// compare prefixes based on their broadcast address
function comparePrefix(prefix1, prefix2) {
    let endOfP1 = inet_aton(prefix1.address) + subnetAddresses(prefix1.prefix)
    let endOfP2 = inet_aton(prefix2.address) + subnetAddresses(prefix2.prefix)
    return endOfP1 - endOfP2
}

// sum up capacity of an array of addresses
function prefixesCapacity(currentResources) {
    let width = 0
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

    let utilisedCapacity = utilizedCapacity(allocatedPrefixes, newlyAllocatedPrefixCapacity)
    let remainingCapacity = freeCapacity(parentRange, utilisedCapacity)
    let utilPercentage
    if (remainingCapacity === 0) {
        utilPercentage = 100.0
    } else {
        utilPercentage = (utilisedCapacity / subnetAddresses(parentRange.prefix)) * 100
    }
    console[level]("Remaining capacity: " + remainingCapacity)
    console[level]("Utilised capacity: " + utilisedCapacity)
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
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
    return ` + "`" + `${prefix.address}/${prefix.prefix}` + "`" + `
}

function prefixToRangeStr(prefix) {
    return ` + "`" + `[${prefix.address}-${inet_ntoa(inet_aton(prefix.address) + subnetAddresses(prefix.prefix) - 1)}]` + "`" + `
}

function calculateDesiredSubnetMask() {
    let newSubnetBits = Math.ceil(Math.log(userInput.desiredSize) / Math.log(2))
    let newSubnetMask = 32 - newSubnetBits
    let newSubnetCapacity = subnetAddresses(newSubnetMask)
    return {newSubnetMask, newSubnetCapacity};
}

// calculate the nearest possible address for a subnet where mask === newSubnetMask
//  that is outside of allocatedSubnet
function findNextFreeSubnetAddress(allocatedSubnet, newSubnetMask) {
    // find the first address after currently iterated allocated subnet
    let nextAvailableAddressNum = inet_aton(allocatedSubnet.address) + subnetAddresses(allocatedSubnet.prefix)
    // remove any bites from the address above after newSubnetMask
    let newSubnetMaskNegative = 32 - newSubnetMask
    let possibleSubnetNum = (nextAvailableAddressNum >>> newSubnetMaskNegative) << newSubnetMaskNegative
    // keep going until we find an address outside of currently iterated allocated subnet
    while (nextAvailableAddressNum > possibleSubnetNum) {
        possibleSubnetNum = ((possibleSubnetNum >>> newSubnetMaskNegative) + 1) << newSubnetMaskNegative
    }
    return possibleSubnetNum;
}

// main
function invoke() {
    let rootPrefix = resourcePool.ResourcePoolName

    let rootPrefixParsed = parsePrefix(rootPrefix)
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

    if (userInput.desiredSize < 2) {
        console.error("Unable to allocate subnet from root prefix: " + rootPrefixStr +
            ". Desired size is invalid: " + userInput.desiredSize + ". Use values >= 2")
        return null
    }

    if (userInput.subnet === true) {
        // reserve subnet address and broadcast
        userInput.desiredSize += 2
    }

    // Calculate smallest possible subnet mask to fit desiredSize
    let {newSubnetMask, newSubnetCapacity} = calculateDesiredSubnetMask();

    // unwrap and sort currentResources
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
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
        possibleSubnetNum = findNextFreeSubnetAddress(allocatedSubnet, newSubnetMask);
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
    console.error("Unable to allocate Ipv4 prefix from: " + rootPrefixStr +
        ". Insufficient capacity to allocate a new prefix of size: " + userInput.desiredSize)
    console.error("Currently allocated prefixes: " + prefixesToString(currentResourcesUnwrapped))
    logStats(null, rootPrefixParsed, currentResourcesUnwrapped, "error")
    return null
}


`

const IPV4 = `

/*
IPv4 address allocation strategy

- Expects IPv4 prefix resource type to have 2 properties of type int ["address:string", "mask:int"]
- userInput.subnet is an optional parameter specifying whether root prefix will be used as a real subnet or just
  as an IP pool. Essentially whether to consider subnet address and broadcast when allocating addresses.
- Logs utilisation stats
- Allocates previously freed prefixes
- All addresses from parent prefix are used, including the first and last one
 */

// ipv4 int to str
function inet_ntoa(addrint) {
    return ((addrint >> 24) & 0xff) + "." +
        ((addrint >> 16) & 0xff) + "." +
        ((addrint >> 8) & 0xff) + "." +
        (addrint & 0xff)
}

// ipv4 str to int
function inet_aton(addrstr) {
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
}

// number of addresses in a subnet based on its mask
function subnetAddresses(mask) {
    return 1<<(32-mask)
}

const prefixRegex = /([0-9.]+)\/([0-9]{1,2})/

// parse prefix from a string e.g. 1.2.3.4/18 into an object
function parsePrefix(str) {
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
        addrNum = (addrNum >>> (32-mask)) << (32-mask)
    }

    return {
        "address": inet_ntoa(addrNum),
        "prefix": mask
    }
}

// calculate utilized capacity based on previously allocated prefixes + a newly allocated prefix
function utilizedCapacity(allocatedAddresses, newlyAllocatedRangeCapacity) {
    return allocatedAddresses.length + newlyAllocatedRangeCapacity
}

// calculate free capacity based on previously allocated prefixes
function freeCapacity(parentPrefix, utilisedCapacity) {
    return subnetAddresses(parentPrefix.prefix) - utilisedCapacity
}

// log utilisation stats
function logStats(newlyAllocatedAddr, parentRange, isSubnet = false, allocatedAddresses = [], level = "log") {
    let newlyAllocatedPrefixCapacity = 0
    if (newlyAllocatedAddr) {
        newlyAllocatedPrefixCapacity = 1
    } else {
        newlyAllocatedPrefixCapacity = 0
    }

    let utilisedCapacity = utilizedCapacity(allocatedAddresses, newlyAllocatedPrefixCapacity)
    if(isSubnet) {
        utilisedCapacity += 2
    }
    let remainingCapacity = freeCapacity(parentRange, utilisedCapacity)
    let utilPercentage
    if (remainingCapacity === 0) {
        utilPercentage = 100.0
    } else {
        utilPercentage = (utilisedCapacity / subnetAddresses(parentRange.prefix)) * 100
    }
    console[level]("Remaining capacity: " + remainingCapacity)
    console[level]("Utilised capacity: " + utilisedCapacity)
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
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
    return ` + "`" + `${prefix.address}/${prefix.prefix}` + "`" + `
}

// main
function invoke() {
    let rootPrefix = resourcePool.ResourcePoolName

    let rootPrefixParsed = parsePrefix(rootPrefix)
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
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
    let currentResourcesSet = new Set(currentResourcesUnwrapped.map(ip => ip.address))

    let firstPossibleAddr = 0
    let lastPossibleAddr = 0
    if (userInput.subnet === true) {
        firstPossibleAddr = rootAddressNum + 1
        lastPossibleAddr = rootAddressNum + rootCapacity - 1
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
    console.error("Unable to allocate Ipv4 address from: " + rootPrefixStr +
        ". Insufficient capacity to allocate a new address")
    console.error("Currently allocated addresses: " + addressesToStr(currentResourcesUnwrapped))
    logStats(null, rootPrefixParsed, userInput.subnet === true, currentResourcesUnwrapped, "error")
    return null
}


`

const IPV6_PREFIX = `

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

        for (const n of parts.map(part => part ? ` + "`" + `0x${part}` + "`" + ` : ` + "`" + `0` + "`" + `).map(Number).reverse()) {
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
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
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
    return ` + "`" + `${prefix.address}/${prefix.prefix}` + "`" + `
}

function prefixToRangeStr(prefix) {
    return ` + "`" + `[${prefix.address}-${inet_ntoa(inet_aton(prefix.address) + subnetAddresses(prefix.prefix) - BigInt(1))}]` + "`" + `
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
    let rootPrefix = resourcePool.ResourcePoolName

    let rootPrefixParsed = parsePrefix(rootPrefix)
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
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
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


`

const IPV6 = `

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

        for (const n of parts.map(part => part ? ` + "`" + `0x${part}` + "`" + ` : ` + "`" + `0` + "`" + `).map(Number).reverse()) {
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
    return ` + "`" + `${prefix.address}/${prefix.prefix}` + "`" + `
}

// calculate utilized capacity based on previously allocated prefixes + a newly allocated prefix
function utilizedCapacity(allocatedAddresses, newlyAllocatedRangeCapacity) {
    return BigInt(allocatedAddresses.length) + BigInt(newlyAllocatedRangeCapacity)
}

// calculate free capacity based on previously allocated prefixes
function freeCapacity(parentPrefix, utilisedCapacity) {
    return subnetAddresses(parentPrefix.prefix) - BigInt(utilisedCapacity)
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
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
}

// main
function invoke() {
    let rootPrefix = resourcePool.ResourcePoolName

    let rootPrefixParsed = parsePrefix(rootPrefix)
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
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
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


`

const RANDOM_S_INT32 = `

/*
signed int32 random allocation strategy

- Expects random signed int32 resource type to have 1 properties of type int ["int"]
- Logs utilisation stats
- MIN value is  -2147483648, MAX value is 2147483647
- Allocates previously freed resources
 */

const rangeRegx = /\[(-?[0-9]+)-([0-9]+)\]/

const INT_MIN = -2147483648
const INT_MAX = 2147483647

// TODO this parse_range function is pretty common, how to share common code ? this is different to 3rd party libs

function parse_range(str) {
    let res = rangeRegx.exec(str)
    if (res == null) {
        console.error("Int range cannot be parsed from pool name: " + str + ". Not matching pattern: " + rangeRegx)
        return null
    }

    from = parseInt(res[1])
    to = parseInt(res[2])
    if (from < INT_MIN || from >= INT_MAX) {
        console.error("Int range invalid, from end is: " + from)
        return null
    }
    if (to <= INT_MIN || to > INT_MAX) {
        console.error("Int range invalid, to end is: " + from)
        return null
    }

    if (from >= to) {
        console.error("Int range invalid, from end: " + from + " and to end: " + to + " do not form a range")
        return null
    }

    return {
        "from": from,
        "to": to
    }
}

function rangeCapacity(s_int32Range) {
    return s_int32Range.to - s_int32Range.from + 1
}

function freeCapacity(parentRange, utilisedCapacity) {
    return rangeCapacity(parentRange) - utilisedCapacity
}

function utilizedCapacity(allocatedRanges, newlyAllocatedS_int32) {
    return allocatedRanges.length + (newlyAllocatedS_int32 != null)
}

function logStats(newlyAllocatedS_int32, parentRange, allocatedS_int32s = [], level = "log") {
    let utilisedCapacity = utilizedCapacity(allocatedS_int32s, newlyAllocatedS_int32)
    let remainingCapacity = freeCapacity(parentRange, utilisedCapacity)
    let utilPercentage
    if (remainingCapacity === 0) {
        utilPercentage = 100.0
    } else {
        utilPercentage = (utilisedCapacity / rangeCapacity(parentRange)) * 100
    }
    console[level]("Remaining capacity: " + remainingCapacity)
    console[level]("Utilised capacity: " + utilisedCapacity)
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
}

function invoke() {
    let parentRangeStr = resourcePool.ResourcePoolName
    let parentRange = parse_range(parentRangeStr)
    if (parentRange == null) {
        console.error("Unable to allocate random s_int32" +
            ". Unable to extract parent int range from pool name: " + parentRangeStr)
        return null
    }

    // unwrap currentResources
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
    let currentResourcesSet = new Set(currentResourcesUnwrapped.map(s_int32 => s_int32.int))

    for (let i = parentRange.from; i <= parentRange.to; i++) {
        newInt = getRandomInt(parentRange.from, parentRange.to)
        if (!currentResourcesSet.has(newInt)) {
            // FIXME How to pass these stats ?
            // logStats(i, parentRange, currentResourcesUnwrapped)
            return {
                "int": newInt
            }
        }
    }

    // no more numbers
    console.error("Unable to allocate random s_int32 from: " + rangeToStr(parentRange) +
        ". Insufficient capacity to allocate a new s_int32")
    logStats(null, parentRange, currentResourcesUnwrapped, "error")
    return null
}

// returns an int between min and max (inclusive)
function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 0.9999) + min);
}

function rangeToStr(range) {
    return ` + "`" + `[${range.from}-${range.to}]` + "`" + `
}


`

const ROUTE_DISTINGUISHER = `

/*
RD allocation strategy - this strategy expects all inputs to be provided and just validates and formats them

- Expects route_distuingusher resource type to have 1 properties of type string ["rd:string"]
- Input params
-  userInput.ipv4 - valid ipv4 addr
-  userInput.as - valid AS number (2bytes or 4 bytes)
-  userInput.assignedNumber
- Valid input combinations are: ipv4 + assignedNumber(2 bytes), as(4 byte) + assignedNumber(2 bytes), as(2 byte) + assignedNumber(4 bytes)
- Allocates previously freed RDs
 */


const BYTE_2_MAX = 65536
const BYTE_4_MAX = 4294967296
const BYTE_6_MAX = BigInt(256) ** BigInt(6)

function rangeCapacity() {
    // this is theoretical max capacity for 6 bytes
    return BYTE_6_MAX
}

function freeCapacity(utilisedCapacity) {
    return rangeCapacity() - utilisedCapacity
}

function utilizedCapacity(allocatedRanges, newlyAllocatedVlan) {
    // FIXME using BigInts but allocatedRanges only fits 2^32 items MAX
    return BigInt(allocatedRanges.length) + BigInt(newlyAllocatedVlan != null)
}

function logStats(newlyAllocatedRd, allocatedRds = [], level = "log") {
    let utilisedCapacity = utilizedCapacity(allocatedRds, newlyAllocatedRd)
    let remainingCapacity = freeCapacity( utilisedCapacity)
    let utilPercentage
    if (remainingCapacity === BigInt(0)) {
        utilPercentage = 100.0
    } else {
        let utilFloat = Number(utilisedCapacity * BigInt(1000) / rangeCapacity()) / 1000
        utilPercentage = (utilFloat * 100)
    }
    
    console[level]("Remaining capacity: " + remainingCapacity)
    console[level]("Utilised capacity: " + utilisedCapacity)
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
}

function invoke() {
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
    let currentResourcesSet = new Set(currentResourcesUnwrapped.map(ip => ip.rd))

    let is2ByteAssignedNumber = false
    let assignedNumber = -1
    if (userInput.hasOwnProperty("assignedNumber") && userInput.assignedNumber) {
        assignedNumber = userInput.assignedNumber
        if (assignedNumber > 0 && assignedNumber < BYTE_2_MAX) {
            is2ByteAssignedNumber = true
        } else if (assignedNumber > 0 && assignedNumber < BYTE_4_MAX) {
            is2ByteAssignedNumber = false
        } else {
            console.error("Unable to allocate RD for assigned number: " + userInput.assignedNumber + ". Number is invalid")
            logStats(null, currentResourcesUnwrapped, "error")
            return null
        }
    }

    let is2ByteAs = false
    let asNumber = -1
    if (userInput.hasOwnProperty("asNumber") && userInput.asNumber) {
        asNumber = userInput.asNumber
        if (asNumber > 0 && asNumber < BYTE_2_MAX) {
            is2ByteAs = true
        } else if (asNumber > 0 && asNumber < BYTE_4_MAX) {
            is2ByteAs = false
        } else {
            console.error("Unable to allocate RD for AS number: " + userInput.asNumber + ". AS is invalid")
            logStats(null, currentResourcesUnwrapped, "error")
            return null
        }
    }

    let ipv4 = "-1"
    if (userInput.hasOwnProperty("ipv4") && userInput.ipv4) {

        var re = /^([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})$/
        var res = re.exec(userInput.ipv4)

        if (res === null) {
            console.error("Unable to allocate RD, invalid IPv4: " + userInput.ipv4 + " provided")
            logStats(null, currentResourcesUnwrapped, "error")
            return null
        }

        if (userInput.ipv4 > "255.255.255.255") {
            console.error("Unable to allocate RD, invalid IPv4: " + userInput.ipv4 + " provided")
            logStats(null, currentResourcesUnwrapped, "error")
            return null
        }

        ipv4 = userInput.ipv4
    }

    if (ipv4 !== "-1" && asNumber !== -1) {
        console.error("Unable to allocate RD, both AS: " + asNumber + " number and IPv4: " + ipv4 + " provided")
        logStats(null, currentResourcesUnwrapped, "error")
        return null
    }

    if (asNumber !== -1 && !is2ByteAs && assignedNumber !== -1 && !is2ByteAssignedNumber) {
        console.error("Unable to allocate RD, 4 byte AS: " + asNumber + " and 4 byte assigned number: " + assignedNumber + " provided")
        logStats(null, currentResourcesUnwrapped, "error")
        return null
    }

    if (ipv4 !== "-1" && assignedNumber !== -1 && !is2ByteAssignedNumber) {
        console.error("Unable to allocate RD, 4 byte assigned number: " + assignedNumber + " provided with an IP address")
        logStats(null, currentResourcesUnwrapped, "error")
        return null
    }

    // TYPE 0
    if (asNumber !== -1 && is2ByteAs && assignedNumber !== -1) {
        let newRd = ` + "`" + `${asNumber}:${assignedNumber}` + "`" + `
        if (currentResourcesSet.has(newRd)) {
            console.error("Unable to allocate RD, duplicate RD created: " + newRd)
            logStats(null, currentResourcesUnwrapped, "error")
            return null
        }
        // logStats(newRd, currentResourcesUnwrapped)
        return {"rd": newRd}
    }

    // TYPE 1
    if (ipv4 !== "-1" && assignedNumber !== -1 && is2ByteAssignedNumber) {
        let newRd = ` + "`" + `${ipv4}:${assignedNumber}` + "`" + `
        if (currentResourcesSet.has(newRd)) {
            console.error("Unable to allocate RD, duplicate RD created: " + newRd)
            logStats(null, currentResourcesUnwrapped, "error")
            return null
        }
        // logStats(newRd, currentResourcesUnwrapped)
        return {"rd": newRd}
    }

    // TYPE 2
    if (asNumber !== -1 && !is2ByteAs && assignedNumber !== -1 && is2ByteAssignedNumber) {
        let newRd = ` + "`" + `${asNumber}:${assignedNumber}` + "`" + `
        if (currentResourcesSet.has(newRd)) {
            console.error("Unable to allocate RD, duplicate RD created: " + newRd)
            logStats(null, currentResourcesUnwrapped, "error")
            return null
        }
        // logStats(newRd, currentResourcesUnwrapped)
        return {"rd": newRd}
    }

    // didn't match known RD types
    console.error("Unable to allocate RD, check the input parameters. User provided input: " + JSON.stringify(userInput))
    logStats(null, currentResourcesUnwrapped, "error")
    return null
}


`

const VLAN_RANGE = `

/*
VLAN range allocation strategy

- Expects VLAN_range resource type to have 2 properties of type int ["from", "to"]
- Produced ranges are inclusive
- Produced ranges are non-overlapping
- Logs utilisation stats
- userInput.desiredSize is a required parameter e.g. desiredSize == 10  ---produces-range-of--->  [0, 9]
- MIN value is 0, MAX value is 4095
- 0 and 4095 are not reserved !
- Allocates previously freed resources
 */

const rangeRegx = /\[([0-9]+)-([0-9]+)\]/

const VLAN_MIN = 0
const VLAN_MAX = 4095

function parse_range(str) {
    let res = rangeRegx.exec(str)
    if (res == null) {
        console.error("VLAN range cannot be parsed from pool name: " + str + ". Not matching pattern: " + rangeRegx)
        return null
    }

    from = parseInt(res[1])
    to = parseInt(res[2])
    if (from < VLAN_MIN || from >= VLAN_MAX) {
        console.error("VLAN range invalid, from end is: " + from)
        return null
    }
    if (to <= VLAN_MIN || to > VLAN_MAX) {
        console.error("VLAN range invalid, to end is: " + from)
        return null
    }

    if (from >= to) {
        console.error("VLAN range invalid, from end: " + from + " and to end: " + to + " do not form a range")
        return null
    }

    return {
        "from": from,
        "to": to
    }
}

function rangeCapacity(vlanRange) {
    return vlanRange.to - vlanRange.from + 1
}

function rangesToStr(currentResources) {
    let subRangesToString = ""
    for (let allocatedRange of currentResources) {
        subRangesToString += rangeToStr(allocatedRange)
    }
    return subRangesToString
}

function rangesCapacity(currentResources) {
    let width = 0
    for (let allocatedRange of currentResources) {
        width += rangeCapacity(allocatedRange)
    }
    return width
}

function freeCapacity(parentRange, utilisedCapacity) {
    return rangeCapacity(parentRange) - utilisedCapacity
}

function utilizedCapacity(allocatedRanges, newlyAllocatedRangeCapacity) {
    return rangesCapacity(allocatedRanges) + newlyAllocatedRangeCapacity
}

function logStats(newlyAllocatedRange, parentRange, allocatedRanges = [], level = "log") {
    let newlyAllocatedRangeCapacity = 0
    if (newlyAllocatedRange) {
        newlyAllocatedRangeCapacity = rangeCapacity(newlyAllocatedRange);
    } else {
        newlyAllocatedRangeCapacity = 0
    }

    let utilisedCapacity = utilizedCapacity(allocatedRanges, newlyAllocatedRangeCapacity)
    let remainingCapacity = freeCapacity(parentRange, utilisedCapacity)
    let utilPercentage
    if (remainingCapacity === 0) {
        utilPercentage = 100.0
    } else {
        utilPercentage = (utilisedCapacity / rangeCapacity(parentRange)) * 100
    }
    console[level]("Remaining capacity: " + remainingCapacity)
    console[level]("Utilised capacity: " + utilisedCapacity)
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
}

function invoke() {
    let parentRangeStr = resourcePool.ResourcePoolName
    let parentRange = parse_range(parentRangeStr)
    if (parentRange == null) {
        console.error("Unable to allocate VLAN range" +
            ". Unable to extract parent vlan range from pool name: " + parentRangeStr)
        return null
    }

    if (!userInput.desiredSize) {
        console.error("Unable to allocate VLAN range from: " + rangeToStr(parentRange) +
            ". Desired size of a new vlan range not provided as userInput.desiredSize")
        return null
    }

    if (userInput.desiredSize < 1) {
        console.error("Unable to allocate VLAN range from: " + rangeToStr(parentRange) +
            ". Desired size is invalid: " + userInput.desiredSize + ". Use values >= 1")
        return null
    }

    // unwrap currentResources
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
    // make sure to sort ranges
    currentResourcesUnwrapped.sort(compareVlanRanges)

    let findingAvailableRange = {
        "from": parentRange.from,
        "to": parentRange.to
    }

    // iterate over allocated ranges and see if a desired new range can be squeezed in
    for (let allocatedRange of currentResourcesUnwrapped) {
        // set to bound to from bound of next range
        findingAvailableRange.to = allocatedRange.from - 1
        // if there is enough space, allocate a chunk of that range
        if (rangeCapacity(findingAvailableRange) >= userInput.desiredSize) {
            findingAvailableRange.to = findingAvailableRange.from + userInput.desiredSize - 1
            // FIXME How to pass these stats ?
            // logStats(findingAvailableRange, parentRange, currentResourcesUnwrapped)
            return findingAvailableRange
        }

        findingAvailableRange.from = allocatedRange.to + 1
        findingAvailableRange.to = allocatedRange.to + 1
    }

    // check if there is some space left at the end of parent range
    findingAvailableRange.to = parentRange.to
    if (rangeCapacity(findingAvailableRange) >= userInput.desiredSize) {
        findingAvailableRange.to = findingAvailableRange.from + userInput.desiredSize - 1
        // FIXME How to pass these stats ?
        // logStats(findingAvailableRange, parentRange, currentResourcesUnwrapped)
        return findingAvailableRange
    }

    // no suitable range found
    console.error("Unable to allocate VLAN range from: " + rangeToStr(parentRange) +
        ". Insufficient capacity to allocate a new range of size: " + userInput.desiredSize)
    console.error("Currently allocated ranges: " + rangesToStr(currentResourcesUnwrapped))
    logStats(null, parentRange, currentResourcesUnwrapped, "error")
    return null
}

function compareVlanRanges(range1, range2) {
    // assuming non overlapping ranges
    return range1.to - range2.to
}

function rangeToStr(range) {
    return ` + "`" + `[${range.from}-${range.to}]` + "`" + `
}


`

const VLAN = `

/*
VLAN allocation strategy

- Expects VLAN resource type to have 1 properties of type int ["vlan"]
- Logs utilisation stats
- MIN value is 0, MAX value is 4095
- 0 and 4095 are not reserved !
- Allocates previously freed resources
 */

const rangeRegx = /\[([0-9]+)-([0-9]+)\]/

const VLAN_MIN = 0
const VLAN_MAX = 4095

function parse_range(str) {
    let res = rangeRegx.exec(str)
    if (res == null) {
        console.error("VLAN range cannot be parsed from pool name: " + str + ". Not matching pattern: " + rangeRegx)
        return null
    }

    from = parseInt(res[1])
    to = parseInt(res[2])
    if (from < VLAN_MIN || from >= VLAN_MAX) {
        console.error("VLAN range invalid, from end is: " + from)
        return null
    }
    if (to <= VLAN_MIN || to > VLAN_MAX) {
        console.error("VLAN range invalid, to end is: " + from)
        return null
    }

    if (from >= to) {
        console.error("VLAN range invalid, from end: " + from + " and to end: " + to + " do not form a range")
        return null
    }

    return {
        "from": from,
        "to": to
    }
}

function rangeCapacity(vlanRange) {
    return vlanRange.to - vlanRange.from + 1
}

function freeCapacity(parentRange, utilisedCapacity) {
    return rangeCapacity(parentRange) - utilisedCapacity
}

function utilizedCapacity(allocatedRanges, newlyAllocatedVlan) {
    return allocatedRanges.length + (newlyAllocatedVlan != null)
}

function logStats(newlyAllocatedVlan, parentRange, allocatedVlans = [], level = "log") {
    let utilisedCapacity = utilizedCapacity(allocatedVlans, newlyAllocatedVlan)
    let remainingCapacity = freeCapacity(parentRange, utilisedCapacity)
    let utilPercentage
    if (remainingCapacity === 0) {
        utilPercentage = 100.0
    } else {
        utilPercentage = (utilisedCapacity / rangeCapacity(parentRange)) * 100
    }
    console[level]("Remaining capacity: " + remainingCapacity)
    console[level]("Utilised capacity: " + utilisedCapacity)
    console[level](` + "`" + `Utilisation: ${utilPercentage.toFixed(1)}%` + "`" + `)
}

function invoke() {
    let parentRangeStr = resourcePool.ResourcePoolName
    let parentRange = parse_range(parentRangeStr)
    if (parentRange == null) {
        console.error("Unable to allocate VLAN" +
            ". Unable to extract parent vlan range from pool name: " + parentRangeStr)
        return null
    }

    // unwrap currentResources
    currentResourcesUnwrapped = currentResources.map(cR => cR.Properties)
    let currentResourcesSet = new Set(currentResourcesUnwrapped.map(vlan => vlan.vlan))

    for (let i = parentRange.from; i <= parentRange.to; i++) {
        if (!currentResourcesSet.has(i)) {
            // FIXME How to pass these stats ?
            // logStats(i, parentRange, currentResourcesUnwrapped)
            return {
                "vlan": i
            }
        }
    }

    // no more vlans
    console.error("Unable to allocate VLAN from: " + rangeToStr(parentRange) +
        ". Insufficient capacity to allocate a new vlan")
    logStats(null, parentRange, currentResourcesUnwrapped, "error")
    return null
}

function rangeToStr(range) {
    return ` + "`" + `[${range.from}-${range.to}]` + "`" + `
}


`
