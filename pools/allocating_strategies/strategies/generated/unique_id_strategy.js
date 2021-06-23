'use strict';

// framework managed constants
//;
//;
// framework managed constants

// STRATEGY_START

/*
Unique id generator
- this strategy accepts text template as "idFormat" and will replace variables in {} for values
- {counter} is mandatory static variable for iterations
- example: 'VPN-{counter}-{network}-{vpn}-local'
 */

String.prototype.format = function(dict) {
    return this.replace(/{(\w+)}/g, function(match, key) {
        return typeof dict[key] !== 'undefined'
            ? dict[key]
            : match
            ;
    });
};

function getNextFreeCounter() {
    let max = 0;
    for (let i = 0; i < currentResources.length; i++) {
        if (currentResources[i].Properties.counter > max) {
            max = currentResources[i].Properties.counter;
        }
    }
    return ++max;
}

// main
function invoke() {
    let nextFreeCounter = getNextFreeCounter();
    if (resourcePoolProperties == null) {
        console.error("Unable to extract resources");
        return null
    }
    if (!("idFormat" in resourcePoolProperties)) {
        console.error("Missing idFormat in resources");
        return null
    }
    if (!resourcePoolProperties["idFormat"].includes("{counter}")) {
        console.error("Missing {counter} in idFormat");
        return null
    }
    const { textFunction, ...poolProperties } = resourcePoolProperties;
    let idFormat = resourcePoolProperties["idFormat"].format(
        {...{counter: nextFreeCounter}, ...poolProperties});
    return { text: idFormat, counter: nextFreeCounter };
}

function capacity() {
    let allocatedCapacity = getNextFreeCounter() - 1;
    let freeCapacity = Number.MAX_SAFE_INTEGER - allocatedCapacity;
    return { freeCapacity: freeCapacity, utilizedCapacity: allocatedCapacity };
}

// STRATEGY_END

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg;
    resourcePoolProperties = resourcePoolArg;
    return invoke()
}

function invokeWithParamsCapacity(currentResourcesArg, resourcePoolArg, userInputArg) {
    currentResources = currentResourcesArg;
    resourcePoolProperties = resourcePoolArg;
    return capacity()
}

exports.invoke = invoke;
exports.capacity = capacity;
exports.invokeWithParams = invokeWithParams;
exports.invokeWithParamsCapacity = invokeWithParamsCapacity;
// For testing purposes
