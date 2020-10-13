const strat = require('../ipv6_strategy')

test("allocate all addresses /120 ipv6", () => {
    addresses = []
    for (let i = 1; i < 255; i++) {
        let address = strat.invokeWithParams(addresses,
            { 'prefix': 120, 'address': "dddd::"},
            {"subnet": true})
        addresses.push(addr(address.address))
        expect(address).toStrictEqual(addr("dddd::" + i.toString(16)).Properties)
    }

    // If treated as subnet, prefix is exhausted
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 120, 'address': "dddd::"},
        {"subnet": true})
    ).toStrictEqual(null)

    // If treated as a pool, there are still 2 more addresses left
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 120, 'address': "dddd::"},
        {})
    ).toStrictEqual(addr("dddd::").Properties)
    addresses.push(addr("dddd::"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 120, 'address': "dddd::"},
        {})
    ).toStrictEqual(addr("dddd::ff").Properties)
    addresses.push(addr("dddd::ff"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 120, 'address': "dddd::"},
        {})
    ).toStrictEqual(null)
})

test("allocate all addresses /117 ipv6", () => {
    addresses = []
    for (let i = 0; i < 8; i++) {
        for (let j = 0; j < 256; j++) {
            if (i === 0 && j === 0) {
                // First subnet addr: reserved
                continue
            }
            if (i === 7 && j === 255) {
                // Broadcast: reserved
                continue
            }
            let address = strat.invokeWithParams(addresses,
                { 'prefix': 115, 'address': "dddd::"},
                {"subnet": true})
            addresses.push(addr(address.address))

            // some formatting adjustments
            let byteI = i.toString(16)
            let byteJ = j.toString(16)
            if (byteI === "0") {
                byteI = ""
                if (byteJ === "0") {
                    byteJ = ""
                }
            } else if (byteJ.length === 1) {
                byteJ = "0" + byteJ
            }

            expect(address).toStrictEqual(addr("dddd::" + byteI + "" + byteJ).Properties)
        }
    }

    // If treated as subnet, prefix is exhausted
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 117, 'address': "dddd::"},
        {"subnet": true})
    ).toStrictEqual(null)

    // If treated as a pool, there are still 2 more addresses left
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 117, 'address': "dddd::"},
        {})
    ).toStrictEqual(addr("dddd::").Properties)
    addresses.push(addr("dddd::"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 117, 'address': "dddd::"},
        {})
    ).toStrictEqual(addr("dddd::7ff").Properties)
    addresses.push(addr("dddd::7ff"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 117, 'address': "dddd::"},
        {})
    ).toStrictEqual(null)
})

test("allocate ipv6 at start with existing resources", () => {
    let subnet = strat.invokeWithParams(
        [addr("dead::2")],
        { 'prefix': 24, 'address': "dead::"},
        {"subnet": true});
    expect(subnet)
        .toStrictEqual(addr("dead::1").Properties)
})

test("ipv6 capacity 24 mask", () => {
    let capacity = strat.invokeWithParamsCapacity(
        [addr("dead::2")],
        { 'prefix': 110, 'address': "dead::"},
        {"subnet": true});
    expect(capacity)
        .toStrictEqual({freeCapacity: 262144, utilizedCapacity: 1})
})

function addr(ip) {
    return {"Properties": {"address": ip}}
}

test("free ipv6 capacity", () => {
    expect(strat.freeCapacity({"prefix": 120}, 100)).toStrictEqual(BigInt(156))
})

test("ipv6 utilisation", () => {
    expect(strat.utilizedCapacity(
        [addr("dead::beed").Properties],
        1))
        .toStrictEqual(BigInt(2))
})
