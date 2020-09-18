const strat = require('../ipv6_prefix_strategy')

function prefix(ip, prefix) {
    return {"Properties": {"address": ip, "prefix": prefix}}
}

test("ipv6 parse", () => {
    for (let ipv6 of [
        "dead::beef",
        "::1",
        "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
        "::",
        "a897:fedc:1111:9999:f999::abcd"
    ]) {
        console.log(ipv6)
        expect(strat.inet_ntoa(strat.inet_aton(ipv6))).toStrictEqual(ipv6)
    }
})

test("ipv6 parse prefix", () => {
    let expected = [
        prefix("dead::", 64),
        prefix("::", 19),
        prefix("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", 128),
        prefix("ffff:ffff:ffff:ffff:ffff:ffff:ffff:0", 112),
        prefix("ff00::", 8),
        prefix("ffff::", 16),
        prefix("a897:fedc:1111:9999:f999:abcc::", 95),
    ]
    let counter = 0

    for (let ipv6 of [
        "dead::beef/64",
        "::1/19",
        "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128",
        "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/112",
        "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/8",
        "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/16",
        "a897:fedc:1111:9999:f999:abcd::/95",
    ]) {
        expect(strat.parsePrefix(ipv6)).toStrictEqual(expected[counter].Properties)
        counter++
    }
})

test("ipv6 parse invalid", () => {
    for (let ipv6 of [
        "xxxx::yyyy",
        "z",
        "888878468945"
    ]) {
        console.log(ipv6)
        expect(strat.inet_aton(ipv6)).toStrictEqual(null)
    }
})

test("single allocation pool ipv6", () => {
    for (let i = 1; i <= 128-8; i++) {
        let desiredSize = BigInt(2) ** BigInt(i)
        let subnet = strat.invokeWithParams([],
            { 'prefix': 8, 'address': "bb00::"},
            {"desiredSize": desiredSize});
        expect(subnet)
            .toStrictEqual(prefix("bb00::", 128 - i).Properties)
    }
})

test("single allocation pool ipv6 subnet", () => {
    for (let i = 1; i <= 128-8-1; i++) {
        let desiredSize = BigInt(2) ** BigInt(i)
        let subnet = strat.invokeWithParams([],
            { 'prefix': 8, 'address': "bb00::"},
            {"desiredSize": desiredSize, "subnet": true});
        expect(subnet)
            .toStrictEqual(prefix("bb00::", 128 - i - 1).Properties)
    }
})

test("allocate range at start with existing resources ipv6", () => {
    let subnet = strat.invokeWithParams(
        [prefix("dead::be02", 127)],
        { 'prefix': 120, 'address': "dead::be00"},
        {"desiredSize": 2});
    expect(subnet)
        .toStrictEqual(prefix("dead::be00", 127).Properties)
})

test("desired size > than root ipv6", () => {
    expect(strat.invokeWithParams([],
        { 'prefix': 120, 'address': "dead::be00"},
        {"desiredSize": 300}))
        .toStrictEqual(null)
})

test("desired size === root ipv6", () => {
    expect(strat.invokeWithParams([],
        { 'prefix': 104, 'address': "dead::"},
        {"desiredSize": 16777216}))
        .toStrictEqual(prefix("dead::", 104).Properties)
})

test("ipv6 prefix allocation /104", () => {
    let resourcePoolArg = { 'prefix': 104, 'address': "abcd:ef01:2345:6789::"};
    let subnets = []
    let expectedSubnets = [
        prefix("abcd:ef01:2345:6789::", 108),
        prefix("abcd:ef01:2345:6789::20:0", 107),
        prefix("abcd:ef01:2345:6789::40:0", 106),
        prefix("abcd:ef01:2345:6789::10:0", 111),
        prefix("abcd:ef01:2345:6789::80:0", 108),
        prefix("abcd:ef01:2345:6789::18:0", 109),
        prefix("abcd:ef01:2345:6789::c0:0", 106),
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

    // utilisation should be 78.9%

    // Round 2, try to squeeze in additional subnets
    let expectedSubnets2 = [
        prefix("abcd:ef01:2345:6789::14:0", 110),
        prefix("abcd:ef01:2345:6789::a0:0", 107),
        prefix("abcd:ef01:2345:6789::90:0", 109),
        prefix("abcd:ef01:2345:6789::98:0", 110),
        prefix("abcd:ef01:2345:6789::12:0", 111),
        prefix("abcd:ef01:2345:6789::9c:0", 110),
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
