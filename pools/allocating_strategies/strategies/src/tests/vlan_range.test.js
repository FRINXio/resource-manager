const strat = require('../vlan_range_strategy');
import {freeCapacity} from "../vlanutils";

test("missing parent range", () => {
    expect(strat.invokeWithParams([], {}, {}))
        .toStrictEqual(null);
})

test("missing desiredSize", () => {
    expect(strat.invokeWithParams([], {'from': 0, 'to':4095}, {}))
        .toStrictEqual(null);
})

test("allocate range 0", () => {
    expect(strat.invokeWithParams([], {'from': 0, 'to':4095}, {}))
        .toStrictEqual(null);
})


test("allocate range 4096", () => {
    expect(strat.invokeWithParams([], {'from': 0, 'to':4095}, {"desiredSize": 4096}))
        .toStrictEqual(range(0, 4095).Properties);
})

test("allocate range 4097", () => {
    expect(strat.invokeWithParams([], {'from': 0, 'to':4095}, {"desiredSize": 4097}))
        .toStrictEqual(null);
})

test("allocate range no capacity", () => {
    expect(strat.invokeWithParams([range(0, 2000), range(2001, 4090)], {'from': 0, 'to':4095}, {"desiredSize": 100}))
        .toStrictEqual(null);
})

test("allocate range 1", () => {
    expect(strat.invokeWithParams([], {'from': 0, 'to':33}, {"desiredSize": 1}))
        .toStrictEqual(range(0, 0).Properties);
})

test("allocate range 784", () => {
    expect(strat.invokeWithParams([], {'from': 0, 'to':4095}, {"desiredSize": 784}))
        .toStrictEqual(range(0, 783).Properties);
})

test("vlan range capacity measure", () => {
    const capacity = strat.invokeWithParamsCapacity([range(0, 2000), range(2001, 4090)], {'from': 0, 'to':4095}, {"desiredSize": 100});
    expect(capacity)
        .toStrictEqual({freeCapacity: "5", utilizedCapacity: "4091"})
})

function range(from, to) {
    return {"Properties": {"from": from, "to": to}}
}

test("allocate released range", () => {
    expect(strat.invokeWithParams([range(0, 100), range(200, 300)], {'from': 0, 'to':4095}, {"desiredSize": 10}))
        .toStrictEqual(range(101, 110).Properties);

    expect(strat.invokeWithParams([range(0, 100), range(200, 300)], {'from': 0, 'to':4095}, {"desiredSize": 1000}))
        .toStrictEqual(range(301, 1300).Properties);

    expect(strat.invokeWithParams([range(100, 200)], {'from': 0, 'to':4095}, {"desiredSize": 10}))
        .toStrictEqual(range(0, 9).Properties);
})

test("allocate range at the end", () => {
    // allocate range of 1 at the end of parent range
    expect(strat.invokeWithParams([range(0, 1000), range(1001, 3000), range(3001, 4090)], {'from': 0, 'to':4095}, {"desiredSize": 1}))
        .toStrictEqual(range(4091, 4091).Properties);

    // allocate range of 4 at the end of parent range, totally exhausting the range
    expect(strat.invokeWithParams([range(0, 1000), range(1001, 3000), range(3001, 4090), range(4091, 4091)], {'from': 0, 'to':4095}, {"desiredSize": 4}))
        .toStrictEqual(range(4092, 4095).Properties);
})

test("compare vlans", () => {
    expect(strat.compareVlanRanges(range(1, 10).Properties, range(11, 19).Properties)).toBeLessThanOrEqual(-1)
    expect(strat.compareVlanRanges(range(11, 19).Properties, range(1, 10).Properties)).toBeGreaterThanOrEqual(1)
    expect(strat.compareVlanRanges(range(1, 10).Properties, range(1, 10).Properties)).toStrictEqual(0)
})

test("free capacity", () => {
    expect(freeCapacity(range(100, 4055).Properties, 100)).toStrictEqual(3856)
})

test("utilisation", () => {
    expect(strat.utilizedCapacity(
        [range(0, 1000).Properties, range(1001, 3000).Properties, range(3001, 4090).Properties, range(4091, 4091).Properties],
        1))
        .toStrictEqual(4093)
})
