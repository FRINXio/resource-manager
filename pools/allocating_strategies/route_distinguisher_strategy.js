// framework managed constants
var currentResources = []
var resourcePool = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

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
    console[level](`Utilisation: ${utilPercentage.toFixed(1)}%`)
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
        let newRd = `${asNumber}:${assignedNumber}`
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
        let newRd = `${ipv4}:${assignedNumber}`
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
        let newRd = `${asNumber}:${assignedNumber}`
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

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg
    resourcePool = resourcePoolArg
    userInput = userInputArg
    return invoke()
}

exports.invoke = invoke
exports.invokeWithParams = invokeWithParams
exports.utilizedCapacity = utilizedCapacity
exports.freeCapacity = freeCapacity
// For testing purposes
