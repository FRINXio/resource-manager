import {
    addressesToStr,
    hostsInMask,
    inet_aton,
    inet_ntoa,
    parsePrefix,
    prefixToStr,
    subnetAddresses
} from "./ipv6-utils";

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


// calculate utilized capacity based on previously allocated prefixes + a newly allocated prefix
function utilizedCapacity(allocatedAddresses, newlyAllocatedRangeCapacity) {
    return BigInt(allocatedAddresses.length) + BigInt(newlyAllocatedRangeCapacity)
}

// calculate free capacity based on previously allocated prefixes
function freeCapacity(parentPrefix, utilisedCapacity) {
    return subnetAddresses(parentPrefix.prefix) - BigInt(utilisedCapacity)
}

function capacity() {
    let subnetItself = userInput.subnet ? BigInt(1) : BigInt(0);
    let freeInTotal = hostsInMask(resourcePoolProperties.address, resourcePoolProperties.prefix) + subnetItself;
    return {
        freeCapacity: String(freeInTotal - BigInt(currentResources.length)),
        utilizedCapacity: String(currentResources.length)
    };
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
exports.capacity = capacity
exports.invokeWithParams = invokeWithParams
exports.invokeWithParamsCapacity = invokeWithParamsCapacity
exports.parsePrefix = parsePrefix
exports.utilizedCapacity = utilizedCapacity
exports.freeCapacity = freeCapacity
// For testing purposes
