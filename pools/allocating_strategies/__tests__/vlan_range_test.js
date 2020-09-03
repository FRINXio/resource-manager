const strat = require('../vlan_range');
test("missing parent range", () => {
    expect(strat.invokeWithParams([], {}, {}))
        .toStrictEqual(null);
})

test("missing desiredSize", () => {
    expect(strat.invokeWithParams([], {'ResourcePoolName': "[0-4095]"}, {}))
        .toStrictEqual(null);
})

test("allocate range 0", () => {
    expect(strat.invokeWithParams([], {'ResourcePoolName': "[0-4095]"}, {}))
        .toStrictEqual(null);
})


test("allocate range 4096", () => {
    expect(strat.invokeWithParams([], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 4096}))
        .toStrictEqual(range(0, 4095));
})

test("allocate range 4097", () => {
    expect(strat.invokeWithParams([], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 4097}))
        .toStrictEqual(null);
})

test("allocate range no capacity", () => {
    expect(strat.invokeWithParams([range(0, 2000), range(2001, 4090)], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 100}))
        .toStrictEqual(null);
})

test("allocate range 1", () => {
    expect(strat.invokeWithParams([], {'ResourcePoolName': "[0-33]"}, {"desiredSize": 1}))
        .toStrictEqual(range(0, 0));
})

test("allocate range 784", () => {
    expect(strat.invokeWithParams([], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 784}))
        .toStrictEqual(range(0, 783));
})

function range(from, to) {
    return {"Properties": {"from": from, "to": to}}
}

test("allocate released range", () => {
    expect(strat.invokeWithParams([range(0, 100), range(200, 300)], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 10}))
        .toStrictEqual(range(101, 110));

    expect(strat.invokeWithParams([range(0, 100), range(200, 300)], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 1000}))
        .toStrictEqual(range(301, 1300));

    expect(strat.invokeWithParams([range(100, 200)], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 10}))
        .toStrictEqual(range(0, 9));
})

test("allocate range at the end", () => {
    // allocate range of 1 at the end of parent range
    expect(strat.invokeWithParams([range(0, 1000), range(1001, 3000), range(3001, 4090)], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 1}))
        .toStrictEqual(range(4091, 4091));

    // allocate range of 4 at the end of parent range, totally exhausting the range
    expect(strat.invokeWithParams([range(0, 1000), range(1001, 3000), range(3001, 4090), range(4091, 4091)], {'ResourcePoolName': "[0-4095]"}, {"desiredSize": 4}))
        .toStrictEqual(range(4092, 4095));
})

test("compare vlans", () => {
    expect(strat.compareVlanRanges(range(1, 10), range(11, 19))).toBeLessThanOrEqual(-1)
    expect(strat.compareVlanRanges(range(11, 19), range(1, 10))).toBeGreaterThanOrEqual(1)
    expect(strat.compareVlanRanges(range(1, 10), range(1, 10))).toStrictEqual(0)
})

test("free capacity", () => {
    expect(strat.freeCapacity(range(100, 4055), 100)).toStrictEqual(3856)
})

test("utilisation", () => {
    expect(strat.utilizedCapacity(
        [range(0, 1000), range(1001, 3000), range(3001, 4090), range(4091, 4091)],
        1))
        .toStrictEqual(4093)
})
