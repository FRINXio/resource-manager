'use strict';

function rangeCapacity(vlanRange) {
    return vlanRange.to - vlanRange.from + 1
}

function rangeToStr(range) {
    return `[${range.from}-${range.to}]`
}

function freeCapacity(parentRange, utilisedCapacity) {
    return rangeCapacity(parentRange) - utilisedCapacity
}

// framework managed constants
//;
//;

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

function utilizedCapacity(allocatedRanges, newlyAllocatedVlan) {
    return allocatedRanges.length + (newlyAllocatedVlan != null)
}

function logStats(newlyAllocatedVlan, parentRange, allocatedVlans = [], level = "log") {
    let utilisedCapacity = utilizedCapacity(allocatedVlans, newlyAllocatedVlan);
    let remainingCapacity = freeCapacity(parentRange, utilisedCapacity);
    let utilPercentage;
    if (remainingCapacity === 0) {
        utilPercentage = 100.0;
    } else {
        utilPercentage = (utilisedCapacity / rangeCapacity(parentRange)) * 100;
    }
    console[level]("Remaining capacity: " + remainingCapacity);
    console[level]("Utilised capacity: " + utilisedCapacity);
    console[level](`Utilisation: ${utilPercentage.toFixed(1)}%`);
}

function invoke() {
    let parentRange = resourcePoolProperties;
    if (parentRange == null) {
        console.error("Unable to allocate VLAN" +
            ". Unable to extract parent vlan range from pool name: " + resourcePoolProperties);
        return null
    }

    // unwrap currentResources
    let currentResourcesUnwrapped = currentResources.map(cR => cR.Properties);
    let currentResourcesSet = new Set(currentResourcesUnwrapped.map(vlan => vlan.vlan));

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
        ". Insufficient capacity to allocate a new vlan");
    logStats(null, parentRange, currentResourcesUnwrapped, "error");
    return null
}

function capacity() {
    return { freeCapacity: freeCapacity(resourcePoolProperties, currentResources.length), utilizedCapacity: currentResources.length };
}

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg) {
    currentResources = currentResourcesArg;
    resourcePoolProperties = resourcePoolArg;
    return invoke()
}

function invokeWithParamsCapacity(currentResourcesArg, resourcePoolArg) {
    currentResources = currentResourcesArg;
    resourcePoolProperties = resourcePoolArg;
    return capacity()
}

exports.invoke = invoke;
exports.capacity = capacity;
exports.invokeWithParams = invokeWithParams;
exports.invokeWithParamsCapacity = invokeWithParamsCapacity;
exports.utilizedCapacity = utilizedCapacity;
// For testing purposes
