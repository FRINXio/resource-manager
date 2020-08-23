// framework managed constants
var currentResources = []
var resourcePool = {}
// framework managed constants

function inet_ntoa(addrint) {
    return ((addrint >> 24) & 0xff) + "." +
        ((addrint >> 16) & 0xff) + "." +
        ((addrint >> 8) & 0xff) + "." +
        (addrint & 0xff);
}

function inet_aton(addrstr) {
    var re = /^([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})\.([0-9]{1,3})$/;
    var res = re.exec(addrstr);

    if (res === null) {
        return null;
    }

    for (var i = 1; i <= 4; i++) {
        if (res[i] < 0 || res[i] > 255) {
            return null;
        }
    }

    return (res[1] << 24) | (res[2] << 16) | (res[3] << 8) | res[4];
}

function subnet_addresses(mask)
{
    return 1<<(32-mask);
}

const prefixRegex = /([0-9.]+)\/([0-9]{1,2})/;

function parse_prefix(str) {
    let res = prefixRegex.exec(str);
    if (res == null) {
        return null
    }

    return {
        "address": res[1],
        "prefix": parseInt(res[2], 10)
    }
}

function invoke() {
    let rootPrefix = resourcePool.ResourcePoolName

    let rootPrefixParsed = parse_prefix(rootPrefix)
    if (rootPrefixParsed == null) {
        console.error("Unable to extract root prefix from pool name: " + rootPrefix)
        return null
    }

    let addrStr = rootPrefixParsed.address
    let mask = rootPrefixParsed.prefix

    if (currentResources.length === 0) {
        let addrInt = inet_aton(addrStr);
        return {
            "ip": inet_ntoa(addrInt + 1),
        }
    } else if (currentResources.length + 2 >= subnet_addresses(mask)) {
        console.error("No more free resources")
        return null
    } else {
        let addrInt = inet_aton(addrStr);
        return {
            "ip": inet_ntoa(addrInt + currentResources.length + 1),
        }
    }
}

// For testing purposes
function invokeWithParams(currentResourcesArg, resourcePoolArg) {
    currentResources = currentResourcesArg
    resourcePool = resourcePoolArg
    return invoke()
}
exports.invoke = invoke
exports.invokeWithParams = invokeWithParams
// For testing purposes
