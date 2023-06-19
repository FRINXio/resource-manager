const strat = require('../ipv4_prefix_strategy')
import {hostsInMask, parsePrefix, subnetAddresses} from "../ipv4-utils";

test("single allocation pool", () => {
    for (let i = 1; i <= 8; i++) {
        let desiredSize = Math.pow(2, i)
        let subnet = strat.invokeWithParams([],
            { 'prefix': 24, 'address': "192.168.1.0", "subnet": false},
            {"desiredSize": desiredSize});
        expect(subnet)
            .toStrictEqual(prefixWithSubnet("192.168.1.0", 32 - i, false).Properties)
    }
})

test("single allocation subnet", () => {
    for (let i = 1; i <= 7; i++) {
        let desiredSize = Math.pow(2, i)
        let subnet = strat.invokeWithParams([],
            { 'prefix': 24, 'address': "192.168.1.0", "subnet": true},
            {"desiredSize": desiredSize});
        expect(subnet)
            .toStrictEqual(prefixWithSubnet("192.168.1.0", 32 - i - 1, true).Properties)
    }
})

test("allocate range at start with existing resources", () => {
    let subnet = strat.invokeWithParams(
        [prefixWithSubnet("192.168.1.16", 28, false)],
        { 'prefix': 24, 'address': "192.168.1.0", "subnet": false},
        {"desiredSize": 10});
    expect(subnet)
        .toStrictEqual(prefixWithSubnet("192.168.1.0", 28, false).Properties)
})

test("ipv4-prefix capacity 24 mask", () => {
    let capacity = strat.invokeWithParamsCapacity(
        [prefixWithSubnet("192.168.1.16", 28, true)],
        { 'prefix': 24, 'address': "192.168.1.0", subnet: true},
        {});

    expect(capacity)
        .toStrictEqual({freeCapacity: "240", utilizedCapacity: "16"})
})

test("ipv4 prefix allocation subnet vs. pool", () => {
    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0", "subnet": true},
        {"desiredSize": 2}))
        .toStrictEqual(prefixWithSubnet("192.168.1.0", 30, true).Properties)

    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0", "subnet": false},
        {"desiredSize": 2}))
        .toStrictEqual(prefixWithSubnet("192.168.1.0", 31, false).Properties)

    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0", "subnet": false},
        {"desiredSize": 256}))
        .toStrictEqual(prefixWithSubnet("192.168.1.0", 24, false).Properties)

    // 256 desired size for a subnet does not fit into 192.168.1.0/24
    expect(strat.invokeWithParams([],
        { 'prefix': 24, 'address': "192.168.1.0", "subnet": true},
        {"desiredSize": 256}))
        .toStrictEqual(null)
})

test("ipv4 prefix allocation 24", () => {
    let resourcePoolArg = { 'prefix': 24, 'address': "192.168.1.0", "subnet": false};
    let subnets = []
    let expectedSubnets = [
        prefixWithSubnet("192.168.1.0", 28, false),      // 10 ->   0 -  15
        prefixWithSubnet("192.168.1.32", 27, false),     // 19 ->  32 -  63
        prefixWithSubnet("192.168.1.64", 26, false),     // 39 ->  64 - 127
        prefixWithSubnet("192.168.1.16", 31, false),     //  2 ->  16 -  17
        prefixWithSubnet("192.168.1.128", 28, false),    // 14 -> 128 - 143
        prefixWithSubnet("192.168.1.24", 29, false),     //  8 ->  24 -  31
        prefixWithSubnet("192.168.1.192", 26, false),    // 64 -> 196 - 255
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
        prefixWithSubnet("192.168.1.20", 30, false),
        prefixWithSubnet("192.168.1.160", 27, false),
        prefixWithSubnet("192.168.1.144", 29, false),
        prefixWithSubnet("192.168.1.152", 30, false),
        prefixWithSubnet("192.168.1.18", 31, false),
        prefixWithSubnet("192.168.1.156", 30, false)
    ]
    counter = 0
    for (const i of [4, 32, 8, 4, 2, 4]) {
        let subnet = strat.invokeWithParams(
            subnets,
            resourcePoolArg,
            {"desiredSize": i})
        subnets.push(prefixWithSubnet(subnet.address, subnet.prefix, false))
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
    let resourcePoolArg = { 'prefix': 8, 'address': "10.0.0.0", "subnet": false};
    let subnets = []
    let expectedSubnets = [
        prefixWithSubnet("10.0.0.0", 12, false),
        prefixWithSubnet("10.32.0.0", 11, false),
        prefixWithSubnet("10.64.0.0", 10, false),
        prefixWithSubnet("10.16.0.0", 15, false),
        prefixWithSubnet("10.128.0.0", 12, false),
        prefixWithSubnet("10.24.0.0", 13, false),
        prefixWithSubnet("10.192.0.0", 10, false),
    ]
    let counter = 0
    for (const i of [655360, 1245184, 2555904, 131072, 917504, 524288, 4194304]) {
        let subnet = strat.invokeWithParams(
            subnets,
            resourcePoolArg,
            {"desiredSize": i})
        subnets.push(prefixWithSubnet(subnet.address, subnet.prefix, false))
        expect(subnet).toStrictEqual(expectedSubnets[counter].Properties)
        counter++
    }

    // utilisation should be 202/256=78.9%

    // Round 2, try to squeeze in additional subnets
    let expectedSubnets2 = [
        prefixWithSubnet("10.20.0.0", 14, false),
        prefixWithSubnet("10.160.0.0", 11, false),
        prefixWithSubnet("10.144.0.0", 13, false),
        prefixWithSubnet("10.152.0.0", 14, false),
        prefixWithSubnet("10.18.0.0", 15, false),
        prefixWithSubnet("10.156.0.0", 14, false)
    ]
    counter = 0
    for (const i of [262144, 2097152, 524288, 262144, 131072, 262144]) {
        let subnet = strat.invokeWithParams(
            subnets,
            resourcePoolArg,
            {"desiredSize": i})
        subnets.push(prefixWithSubnet(subnet.address, subnet.prefix, false))
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
        { 'prefix': 24, 'address': "192.168.1.0", "subnet": false},
        {"desiredSize": 256}))
        .toStrictEqual(prefixWithSubnet("192.168.1.0", 24, false).Properties)
})

test("parse prefix invalid", () => {
    expect(parsePrefix("255.255.255.256/32")).toStrictEqual(null)
    expect(parsePrefix("1.2.3.4/36")).toStrictEqual(null)
    expect(parsePrefix("abcd")).toStrictEqual(null)
})

test("parse prefix", () => {
    expect(parsePrefix("192.168.1.0/16"))
        .toStrictEqual(prefix("192.168.0.0", 16).Properties)

    expect(parsePrefix("192.168.1.0/8"))
        .toStrictEqual(prefix("192.0.0.0", 8).Properties)

    expect(parsePrefix("255.168.1.0/1"))
        .toStrictEqual(prefix("128.0.0.0", 1).Properties)

    expect(parsePrefix("255.168.1.0/2"))
        .toStrictEqual(prefix("192.0.0.0", 2).Properties)

    expect(parsePrefix("192.168.1.8/32"))
        .toStrictEqual(prefix("192.168.1.8", 32).Properties)

    expect(parsePrefix("192.168.1.0/0"))
        .toStrictEqual(prefix("0.0.0.0", 0).Properties)

    expect(parsePrefix("0.0.0.0/0"))
        .toStrictEqual(prefix("0.0.0.0", 0).Properties)

    expect(parsePrefix("255.255.255.255/32"))
        .toStrictEqual(prefix("255.255.255.255", 32).Properties)
})

test("parse capacity without subnet", () => {
    expect(strat.invokeWithParamsCapacity([prefixWithSubnet("10.0.0.0", 31, false), prefixWithSubnet("10.0.0.2", 31, false)],
        { 'prefix': 8, 'address': "10.0.0.0", subnet: false},
        {"desiredSize": 2}))
        .toStrictEqual({"freeCapacity": "16777212", "utilizedCapacity": "4"})
})

test("parse capacity with subnet", () => {
    expect(strat.invokeWithParamsCapacity([prefixWithSubnet("10.0.0.0", 31, true), prefixWithSubnet("10.0.0.4", 31, true)],
        { 'prefix': 8, 'address': "10.0.0.0", subnet: true},
        {"desiredSize": 2}))
        .toStrictEqual({"freeCapacity": "16777208", "utilizedCapacity": "8"})
})

function prefixWithSubnet(ip, prefix, isSubnet) {
    return {"Properties": {"address": ip, "prefix": prefix, "subnet": isSubnet}}
}

function prefix(ip, prefix) {
    return {"Properties": {"address": ip, "prefix": prefix}}
}

test("free ipv4 prefix capacity", () => {
    expect(strat.freeCapacity(prefixWithSubnet("192.168.1.0", 24, false).Properties, 100)).toStrictEqual(156)
})

test("ipv4 prefix utilisation", () => {
    expect(strat.utilizedCapacity(
        [prefixWithSubnet("192.168.1.0", 28, false).Properties, prefixWithSubnet("192.168.1.128", 27, false).Properties],
        32))
        .toStrictEqual(16+32+32)
})

test('Claim resources from the same pool in parallel way', async (t) => {
    const poolId = await createIpv4PrefixRootPool();

    const promises = [];

    for (let i = 0; i < 100; i++) {
        promises.push(claimResource(poolId, {desiredSize: 2}));
    }

    await Promise.all(promises);

    const pool = await getResourcePool(poolId, undefined, undefined, 120);
    const allocatedResourceProperties = pool.allocatedResources.edges.map(({node}) => node.Properties);

    t.equal(allocatedResourceProperties[0].from, 0);
    t.equal(allocatedResourceProperties[99].from, 99);

    await cleanup();
    t.end();
});
