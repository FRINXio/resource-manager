const strat = require('../ipv4_prefix_strategy')

test("single allocation pool", () => {
    for (let i = 1; i <= 8; i++) {
        let desiredSize = Math.pow(2, i)
        let subnet = strat.invokeWithParams([],
            { 'prefix': 24, 'address': "192.168.1.0"},
            {"desiredSize": desiredSize});
        expect(subnet)
            .toStrictEqual(prefix("192.168.1.0", 32 - i).Properties)
    }
})

test("single allocation subnet", () => {
    for (let i = 1; i <= 7; i++) {
        let desiredSize = Math.pow(2, i)
        let subnet = strat.invokeWithParams([],
            { 'prefix': 24, 'address': "192.168.1.0"},
            {"desiredSize": desiredSize, "subnet": true});
        expect(subnet)
            .toStrictEqual(prefix("192.168.1.0", 32 - i - 1).Properties)
    }
})

test("allocate range at start with existing resources", () => {
    let subnet = strat.invokeWithParams(
        [prefix("192.168.1.16", 28)],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"desiredSize": 10});
    expect(subnet)
        .toStrictEqual(prefix("192.168.1.0", 28).Properties)
})

test("ipv6-prefix capacity 24 mask", () => {
    let capacity = strat.invokeWithParamsCapacity(
        [prefix("192.168.1.16", 28)],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {});

    expect(capacity)
        .toStrictEqual({freeCapacity: 240, utilizedCapacity: 14})
})

test("ipv4 prefix allocation subnet vs. pool", () => {
    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"desiredSize": 2, "subnet": true}))
        .toStrictEqual(prefix("192.168.1.0", 30).Properties)

    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"desiredSize": 2, "subnet": false}))
        .toStrictEqual(prefix("192.168.1.0", 31).Properties)

    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"desiredSize": 256, "subnet": false}))
        .toStrictEqual(prefix("192.168.1.0", 24).Properties)

    // 256 desired size for a subnet does not fit into 192.168.1.0/24
    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"desiredSize": 256, "subnet": true}))
        .toStrictEqual(null)
})

test("ipv4 prefix allocation 24", () => {
    let resourcePoolArg = { 'prefix': 24, 'address': "192.168.1.0"};
    let subnets = []
    let expectedSubnets = [
        prefix("192.168.1.0", 28),      // 10 ->   0 -  15
        prefix("192.168.1.32", 27),     // 19 ->  32 -  63
        prefix("192.168.1.64", 26),     // 39 ->  64 - 127
        prefix("192.168.1.16", 31),     //  2 ->  16 -  17
        prefix("192.168.1.128", 28),    // 14 -> 128 - 143
        prefix("192.168.1.24", 29),     //  8 ->  24 -  31
        prefix("192.168.1.192", 26),    // 64 -> 196 - 255
    ]
    let counter = 0
    for (const i of [10, 19, 39, 2, 14, 8, 64]) {
        let subnet = strat.invokeWithParams(
            subnets,
            resourcePoolArg,
            {"desiredSize": i})
        subnets.push(prefix(subnet.address, subnet.prefix))
        expect(subnet).toStrictEqual(expectedSubnets[counter].Properties)
        counter++
    }

    // utilised capacity should be 16+32+64+2+16+8+64=202
    // free capacity should be 256-202=54
    // utilisation should be 202/256=78.9%
    // free blocks should be: 19-23, 144-195

    // Round 2, try to squeeze in additional subnets
    let expectedSubnets2 = [
        prefix("192.168.1.20", 30),
        prefix("192.168.1.160", 27),
        prefix("192.168.1.144", 29),
        prefix("192.168.1.152", 30),
        prefix("192.168.1.18", 31),
        prefix("192.168.1.156", 30)
    ]
    counter = 0
    for (const i of [4, 32, 8, 4, 2, 4]) {
        let subnet = strat.invokeWithParams(
            subnets,
            resourcePoolArg,
            {"desiredSize": i})
        subnets.push(prefix(subnet.address, subnet.prefix))
        expect(subnet).toStrictEqual(expectedSubnets2[counter].Properties)
        counter++
    }

    // utilised capacity should be 202+4+32+8+4+2+4=256
    // free capacity should be 0
    // utilisation should be 100%

    // Round 3, no more capacity at utilisation 100%
    expect(strat.invokeWithParams(
        subnets,
        resourcePoolArg,
        {"desiredSize": 2})).toStrictEqual(null)
})

// This test is the same as "ipv4 prefix allocation 24" everything is just multiplied by 256*256 to simplify the assertions
test("ipv4 prefix allocation 8", () => {
    let resourcePoolArg = { 'prefix': 8, 'address': "10.0.0.0"};
    let subnets = []
    let expectedSubnets = [
        prefix("10.0.0.0", 12),
        prefix("10.32.0.0", 11),
        prefix("10.64.0.0", 10),
        prefix("10.16.0.0", 15),
        prefix("10.128.0.0", 12),
        prefix("10.24.0.0", 13),
        prefix("10.192.0.0", 10),
    ]
    let counter = 0
    for (const i of [655360, 1245184, 2555904, 131072, 917504, 524288, 4194304]) {
        let subnet = strat.invokeWithParams(
            subnets,
            resourcePoolArg,
            {"desiredSize": i})
        subnets.push(prefix(subnet.address, subnet.prefix))
        expect(subnet).toStrictEqual(expectedSubnets[counter].Properties)
        counter++
    }

    // utilisation should be 202/256=78.9%

    // Round 2, try to squeeze in additional subnets
    let expectedSubnets2 = [
        prefix("10.20.0.0", 14),
        prefix("10.160.0.0", 11),
        prefix("10.144.0.0", 13),
        prefix("10.152.0.0", 14),
        prefix("10.18.0.0", 15),
        prefix("10.156.0.0", 14)
    ]
    counter = 0
    for (const i of [262144, 2097152, 524288, 262144, 131072, 262144]) {
        let subnet = strat.invokeWithParams(
            subnets,
            resourcePoolArg,
            {"desiredSize": i})
        subnets.push(prefix(subnet.address, subnet.prefix))
        expect(subnet).toStrictEqual(expectedSubnets2[counter].Properties)
        counter++
    }

    // utilisation should be 100%

    // Round 3, no more capacity at utilisation 100%
    expect(strat.invokeWithParams(
        subnets,
        resourcePoolArg,
        {"desiredSize": 2})).toStrictEqual(null)
})


test("desired size > than root", () => {
    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"desiredSize": 300}))
        .toStrictEqual(null)
})

test("desired size === than root", () => {
    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0"},
        {"desiredSize": 256}))
        .toStrictEqual(prefix("192.168.1.0", 24).Properties)
})

test("parse prefix invalid", () => {
    expect(strat.parsePrefix("255.255.255.256/32")).toStrictEqual(null)
    expect(strat.parsePrefix("1.2.3.4/36")).toStrictEqual(null)
    expect(strat.parsePrefix("abcd")).toStrictEqual(null)
})

test("parse prefix", () => {
    expect(strat.parsePrefix("192.168.1.0/16"))
        .toStrictEqual(prefix("192.168.0.0", 16).Properties)

    expect(strat.parsePrefix("192.168.1.0/8"))
        .toStrictEqual(prefix("192.0.0.0", 8).Properties)

    expect(strat.parsePrefix("255.168.1.0/1"))
        .toStrictEqual(prefix("128.0.0.0", 1).Properties)

    expect(strat.parsePrefix("255.168.1.0/2"))
        .toStrictEqual(prefix("192.0.0.0", 2).Properties)

    expect(strat.parsePrefix("192.168.1.8/32"))
        .toStrictEqual(prefix("192.168.1.8", 32).Properties)

    expect(strat.parsePrefix("192.168.1.0/0"))
        .toStrictEqual(prefix("0.0.0.0", 0).Properties)

    expect(strat.parsePrefix("0.0.0.0/0"))
        .toStrictEqual(prefix("0.0.0.0", 0).Properties)

    expect(strat.parsePrefix("255.255.255.255/32"))
        .toStrictEqual(prefix("255.255.255.255", 32).Properties)
})

function prefix(ip, prefix) {
    return {"Properties": {"address": ip, "prefix": prefix}}
}

test("free ipv4 prefix capacity", () => {
    expect(strat.freeCapacity(prefix("192.168.1.0", 24).Properties, 100)).toStrictEqual(156)
})

test("ipv4 prefix utilisation", () => {
    expect(strat.utilizedCapacity(
        [prefix("192.168.1.0", 28).Properties, prefix("192.168.1.128", 27).Properties],
        32))
        .toStrictEqual(16+32+32)
})
