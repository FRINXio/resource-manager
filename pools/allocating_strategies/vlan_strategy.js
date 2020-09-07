// framework managed constants
var currentResources = []
var resourcePool = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

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
    console[level](`Utilisation: ${utilPercentage.toFixed(1)}%`)
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
    return `[${range.from}-${range.to}]`
}

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg) {
    currentResources = currentResourcesArg
    resourcePool = resourcePoolArg
    return invoke()
}

exports.invoke = invoke
exports.invokeWithParams = invokeWithParams
exports.utilizedCapacity = utilizedCapacity
exports.freeCapacity = freeCapacity
// For testing purposes
