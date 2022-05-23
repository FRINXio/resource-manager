'use strict';

// framework managed constants
//;
//;
// framework managed constants

// STRATEGY_START

/*
signed int32 random allocation strategy

- Expects random signed int32 resource type to have 1 properties of type int ["int"]
- Logs utilisation stats
- MIN value is  -2147483648, MAX value is 2147483647
- Allocates previously freed resources
 */

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
    let utilisedCapacity = utilizedCapacity(allocatedS_int32s, newlyAllocatedS_int32);
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
        console.error("Unable to allocate random s_int32" +
            ". Unable to extract parent int range from pool name: " + resourcePoolProperties);
        return null
    }

    // unwrap currentResources
    let currentResourcesUnwrapped = currentResources.map(cR => cR.Properties);
    let currentResourcesSet = new Set(currentResourcesUnwrapped.map(s_int32 => s_int32.int));

    for (let i = parentRange.from; i <= parentRange.to; i++) {
        let newInt = getRandomInt(parentRange.from, parentRange.to);
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
        ". Insufficient capacity to allocate a new s_int32");
    logStats(null, parentRange, currentResourcesUnwrapped, "error");
    return null
}

// returns an int between min and max (inclusive)
function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 0.9999) + min);
}

function rangeToStr(range) {
    return `[${range.from}-${range.to}]`
}

function capacity() {
    return {
        freeCapacity: String(freeCapacity(resourcePoolProperties, currentResources.length)),
        utilizedCapacity: String(currentResources.length)
    };
}

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg) {
    currentResources = currentResourcesArg;
    resourcePoolProperties = resourcePoolArg;
    return invoke()
}

exports.invoke = invoke;
exports.capacity = capacity;
exports.invokeWithParams = invokeWithParams;
exports.utilizedCapacity = utilizedCapacity;
exports.freeCapacity = freeCapacity;
// For testing purposes
