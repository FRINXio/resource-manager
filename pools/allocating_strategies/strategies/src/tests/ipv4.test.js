const strat = require('../ipv4_strategy')

test("allocate all addresses 24", () => {
    let addresses = []
    for (let i = 1; i < 255; i++) {
        let address = strat.invokeWithParams(addresses,
            { 'prefix': 24, 'address': "192.168.1.0"},
            {"subnet": true})
        addresses.push(addr(address.address))
        expect(address).toStrictEqual(addr("192.168.1." + i).Properties)
    }

    // If treated as subnet, prefix is exhausted
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"subnet": true})
    ).toStrictEqual(null)

    // If treated as a pool, there are still 2 more addresses left
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 24, 'address': "192.168.1.0"},
        {})
    ).toStrictEqual(addr("192.168.1.0").Properties)
    addresses.push(addr("192.168.1.0"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 24, 'address': "192.168.1.0"},
        {})
    ).toStrictEqual(addr("192.168.1.255").Properties)
    addresses.push(addr("192.168.1.255"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 24, 'address': "192.168.1.0"},
        {})
    ).toStrictEqual(null)
})

test("allocate all addresses 19", () => {
    let addresses = []
    for (let i = 0; i < 32; i++) {
        for (let j = 0; j < 256; j++) {
            if (i === 0 && j === 0) {
                // First subnet addr: reserved
                continue
            }
            if (i === 31 && j === 255) {
                // Broadcast: reserved

                continue
            }
            let address = strat.invokeWithParams(addresses,
                { 'prefix': 19, 'address': "192.168.0.0"},
                {"subnet": true})
            addresses.push(addr(address.address))
            expect(address).toStrictEqual(addr("192.168." + i + "." + j).Properties)
        }
    }

    // If treated as subnet, prefix is exhausted
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 19, 'address': "192.168.0.0"},
        {"subnet": true})
    ).toStrictEqual(null)

    // If treated as a pool, there are still 2 more addresses left
    expect(strat.invokeWithParams(addresses,
        { 'prefix': 19, 'address': "192.168.0.0"},
        {})
    ).toStrictEqual(addr("192.168.0.0").Properties)
    addresses.push(addr("192.168.0.0"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 19, 'address': "192.168.0.0"},
        {})
    ).toStrictEqual(addr("192.168.31.255").Properties)
    addresses.push(addr("192.168.31.255"))

    expect(strat.invokeWithParams(addresses,
        { 'prefix': 19, 'address': "192.168.0.0"},
        {})
    ).toStrictEqual(null)
})

test("allocate ipv4 at start with existing resources", () => {
    let subnet = strat.invokeWithParams(
        [addr("192.168.1.2")],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"subnet": true});
    expect(subnet)
        .toStrictEqual(addr("192.168.1.1").Properties)
})

test("ipv4 capacity 24 mask", () => {
    let capacity = strat.invokeWithParamsCapacity(
        [addr("192.168.1.2")],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {subnet: true});
    expect(capacity)
        .toStrictEqual({freeCapacity: "254", utilizedCapacity: "1"})
})

test("ipv4 capacity 16 mask", () => {
    let capacity = strat.invokeWithParamsCapacity(
        [addr("192.168.1.2")],
        { 'prefix': 16, 'address': "192.168.1.0"},
        {});
    expect(capacity)
        .toStrictEqual({freeCapacity: "65533", utilizedCapacity: "1"})
})

function addr(ip) {
    return {"Properties": {"address": ip}}
}

test("ipv4 utilisation", () => {
    expect(strat.utilizedCapacity(
        [addr("192.168.1.128").Properties],
        1))
        .toStrictEqual(2)
})
