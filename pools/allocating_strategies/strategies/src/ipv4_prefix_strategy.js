import {
    hostsInMask,
    inet_aton,
    inet_ntoa,
    prefixToStr,
    subnetAddresses
} from "./ipv4utils";

// framework managed constants
var currentResources = []
var resourcePoolProperties = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

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

function capacity() {
    let totalCapacity = hostsInMask(resourcePoolProperties.address, resourcePoolProperties.prefix);
    let allocatedCapacity = 0;
    let resource;
    let subnetItself = userInput.subnet ? 1 : 0;

    for (resource of currentResources) {
        allocatedCapacity += hostsInMask(resource.Properties.address, resource.Properties.prefix);
    }

    return { freeCapacity: totalCapacity - allocatedCapacity + subnetItself, utilizedCapacity: allocatedCapacity };
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

function prefixToRangeStr(prefix) {
    return `[${prefix.address}-${inet_ntoa(inet_aton(prefix.address) + subnetAddresses(prefix.prefix) - 1)}]`
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

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg
    resourcePoolProperties = resourcePoolArg
    userInput = userInputArg
    return invoke()
}

// For testing purposes
function invokeWithParamsCapacity(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg
    resourcePoolProperties = resourcePoolArg
    userInput = userInputArg
    return capacity()
}
exports.invoke = invoke
exports.capacity = capacity
exports.invokeWithParams = invokeWithParams
exports.invokeWithParamsCapacity = invokeWithParamsCapacity
exports.utilizedCapacity = utilizedCapacity
exports.freeCapacity = freeCapacity
// For testing purposes
