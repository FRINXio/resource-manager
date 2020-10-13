// framework managed constants
var currentResources = []
var resourcePoolProperties = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

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

// number of assignable addresses based on address and mask
function hostsInMask(addressStr, mask) {
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
