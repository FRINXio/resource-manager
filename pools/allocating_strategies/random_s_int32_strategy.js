// framework managed constants
var currentResources = []
var resourcePool = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

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
    console[level](`Utilisation: ${utilPercentage.toFixed(1)}%`)
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
