import {
    addressesToStr,
    hostsInMask,
    inet_aton,
    inet_ntoa,
    parsePrefix,
    prefixToStr,
    subnetAddresses
} from "./ipv4utils";

// framework managed constants
var currentResources = []
var resourcePoolProperties = {}
var userInput = {}
// framework managed constants

/*
IPv4 address allocation strategy

- Expects IPv4 prefix resource type to have 2 properties of type int ["address:string", "mask:int"]
- userInput.subnet is an optional parameter specifying whether root prefix will be used as a real subnet or just
  as an IP pool. Essentially whether to consider subnet address and broadcast when allocating addresses.
- Logs utilisation stats
- Allocates previously freed prefixes
- All addresses from parent prefix are used, including the first and last one
 */



// calculate utilized capacity based on previously allocated prefixes + a newly allocated prefix
function utilizedCapacity(allocatedAddresses, newlyAllocatedRangeCapacity) {
    return allocatedAddresses.length + newlyAllocatedRangeCapacity
}

// calculate free capacity based on previously allocated prefixes
function freeCapacity(address, mask, utilisedCapacity) {
    let subnetItself = userInput.subnet ? 1 : 0;
    return hostsInMask(address, mask) - utilisedCapacity + subnetItself;
}

function capacity() {
    return { freeCapacity: freeCapacity(resourcePoolProperties.address, resourcePoolProperties.prefix, currentResources.length), utilizedCapacity: currentResources.length };
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
    console[level](`Utilisation: ${utilPercentage.toFixed(1)}%`)
}

// main
function invoke() {
    let rootPrefixParsed = resourcePoolProperties
    if (rootPrefixParsed == null) {
        console.error("Unable to extract root prefix from pool name: " + resourcePoolProperties)
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
    const resourceCount = userInput.resourceCount?userInput.resourceCount:1
    const result = []
    for (let resourceIdx = 0; resourceIdx < resourceCount; resourceIdx++) {
        let firstPossibleAddr = 0
        let lastPossibleAddr = 0
        if (userInput.subnet === true) {
            firstPossibleAddr = rootAddressNum + 1
            lastPossibleAddr = rootAddressNum + rootCapacity - 1
        } else {
            firstPossibleAddr = rootAddressNum
            lastPossibleAddr = rootAddressNum + rootCapacity
        }
        let found = false
        for (let i = firstPossibleAddr; i < lastPossibleAddr && !found; i++) {
            const ipInAscii = inet_ntoa(i)
            if (!currentResourcesSet.has(ipInAscii)) {
                // FIXME How to pass these stats ?
                // logStats(inet_ntoa(i), rootPrefixParsed, userInput.subnet === true, currentResourcesUnwrapped)
                result.push({"address": ipInAscii})
                found = true
                // add it to existing resources
                currentResourcesSet.add(ipInAscii)
            }
        }
        if (!found) {
            // no suitable range found
            console.error(`Unable to allocate Ipv4 address from: ${rootPrefixStr}. ` +
                `Insufficient capacity to allocate ${resourceCount} new address(es)`)
            console.error("Currently allocated addresses: " + addressesToStr(currentResourcesUnwrapped))
            logStats(null, rootPrefixParsed, userInput.subnet === true, currentResourcesUnwrapped, "error")
            return null
        }
    }
    return result
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
exports.utilizedCapacity = utilizedCapacity
exports.freeCapacity = freeCapacity
// For testing purposes
