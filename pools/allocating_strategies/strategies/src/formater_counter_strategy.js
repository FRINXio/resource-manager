// framework managed constants
var currentResources = []
var resourcePoolProperties = {}
var userInput = {}
// framework managed constants

// STRATEGY_START

/*

'VPN-{counter}-
{userInput.someProperty}-{resourcePool.someProperty}-local'
 */

function  getNextFreeCounter() {
    let max = 0;
    for (let i = 0; i < currentResources.length; i++) {
        if (currentResources[i].counter > max) {
            max = currentResources[i].counter;
        }
    }

    return ++max;
}

function invoke() {
    let nextFreeCounter = getNextFreeCounter();
    return {text: userInput.textFunction(userInput, nextFreeCounter, resourcePoolProperties), counter: nextFreeCounter};
}

function capacity() {
    let allocatedCapacity = getNextFreeCounter() - 1;
    let freeCapacity = Number.MAX_SAFE_INTEGER - allocatedCapacity;
    return { freeCapacity: freeCapacity, utilizedCapacity: allocatedCapacity };
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
// For testing purposes
