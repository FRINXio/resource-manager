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
//;
// framework managed constants

// STRATEGY_START

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

function rangesToStr(currentResources) {
    let subRangesToString = "";
    for (let allocatedRange of currentResources) {
        subRangesToString += rangeToStr(allocatedRange);
    }
    return subRangesToString
}

function rangesCapacity(currentResources) {
    let width = 0;
    for (let allocatedRange of currentResources) {
        width += rangeCapacity(allocatedRange);
    }
    return width
}

function utilizedCapacity(allocatedRanges, newlyAllocatedRangeCapacity) {
    return rangesCapacity(allocatedRanges) + newlyAllocatedRangeCapacity
}

function logStats(newlyAllocatedRange, parentRange, allocatedRanges = [], level = "log") {
    let newlyAllocatedRangeCapacity = 0;
    if (newlyAllocatedRange) {
        newlyAllocatedRangeCapacity = rangeCapacity(newlyAllocatedRange);
    } else {
        newlyAllocatedRangeCapacity = 0;
    }

    let utilisedCapacity = utilizedCapacity(allocatedRanges, newlyAllocatedRangeCapacity);
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
        console.error("Unable to allocate VLAN range" +
            ". Unable to extract parent vlan range from pool name: " + resourcePoolProperties);
        return null
    }

    if (!userInput.desiredSize) {
        console.error("Unable to allocate VLAN range from: " + rangeToStr(parentRange) +
            ". Desired size of a new vlan range not provided as userInput.desiredSize");
        return null
    }

    if (userInput.desiredSize < 1) {
        console.error("Unable to allocate VLAN range from: " + rangeToStr(parentRange) +
            ". Desired size is invalid: " + userInput.desiredSize + ". Use values >= 1");
        return null
    }

    // unwrap currentResources
    let currentResourcesUnwrapped = currentResources.map(cR => cR.Properties);
    // make sure to sort ranges
    currentResourcesUnwrapped.sort(compareVlanRanges);

    let findingAvailableRange = {
        "from": parentRange.from,
        "to": parentRange.to
    };

    // iterate over allocated ranges and see if a desired new range can be squeezed in
    for (let allocatedRange of currentResourcesUnwrapped) {
        // set to bound to from bound of next range
        findingAvailableRange.to = allocatedRange.from - 1;
        // if there is enough space, allocate a chunk of that range
        if (rangeCapacity(findingAvailableRange) >= userInput.desiredSize) {
            findingAvailableRange.to = findingAvailableRange.from + userInput.desiredSize - 1;
            // FIXME How to pass these stats ?
            // logStats(findingAvailableRange, parentRange, currentResourcesUnwrapped)
            return findingAvailableRange
        }

        findingAvailableRange.from = allocatedRange.to + 1;
        findingAvailableRange.to = allocatedRange.to + 1;
    }

    // check if there is some space left at the end of parent range
    findingAvailableRange.to = parentRange.to;
    if (rangeCapacity(findingAvailableRange) >= userInput.desiredSize) {
        findingAvailableRange.to = findingAvailableRange.from + userInput.desiredSize - 1;
        // FIXME How to pass these stats ?
        // logStats(findingAvailableRange, parentRange, currentResourcesUnwrapped)
        return findingAvailableRange
    }

    // no suitable range found
    console.error("Unable to allocate VLAN range from: " + rangeToStr(parentRange) +
        ". Insufficient capacity to allocate a new range of size: " + userInput.desiredSize);
    console.error("Currently allocated ranges: " + rangesToStr(currentResourcesUnwrapped));
    logStats(null, parentRange, currentResourcesUnwrapped, "error");
    return null
}

function compareVlanRanges(range1, range2) {
    // assuming non overlapping ranges
    return range1.to - range2.to
}

function capacity() {
    let allocatedCapacity = 0;
    let resource;
    for (resource of currentResources) {
        allocatedCapacity += rangeCapacity(resource.Properties);
    }
    return { freeCapacity: freeCapacity(resourcePoolProperties, allocatedCapacity), utilizedCapacity: allocatedCapacity };
}

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg;
    resourcePoolProperties = resourcePoolArg;
    userInput = userInputArg;
    return invoke()
}

function invokeWithParamsCapacity(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg;
    resourcePoolProperties = resourcePoolArg;
    userInput = userInputArg;
    return capacity()
}

exports.invoke = invoke;
exports.capacity = capacity;
exports.invokeWithParams = invokeWithParams;
exports.invokeWithParamsCapacity = invokeWithParamsCapacity;
exports.compareVlanRanges = compareVlanRanges;
exports.utilizedCapacity = utilizedCapacity;
// For testing purposes
