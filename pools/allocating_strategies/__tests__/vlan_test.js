const strat = require('../vlan')

test("missing parent range", () => {
    expect(strat.invokeWithParams([], {}))
        .toStrictEqual(null)
})

test("allocate vlan", () => {
    expect(strat.invokeWithParams([], {'ResourcePoolName': "[0-4095]"}))
        .toStrictEqual(vlan(0).Properties)
    expect(strat.invokeWithParams([vlan(1)], {'ResourcePoolName': "[0-4095]"}))
        .toStrictEqual(vlan(0).Properties)
    expect(strat.invokeWithParams([vlan(278)], {'ResourcePoolName': "[278-333]"}))
        .toStrictEqual(vlan(279).Properties)
    expect(strat.invokeWithParams(vlans(0, 4094), {'ResourcePoolName': "[0-4095]"}))
        .toStrictEqual(vlan(4095).Properties)
})

test("allocate vlan full", () => {
    expect(strat.invokeWithParams(vlans(0, 4095), {'ResourcePoolName': "[0-4095]"}))
        .toStrictEqual(null)
})

test("free capacity", () => {
    expect(strat.freeCapacity(range(100, 900).Properties, 100)).toStrictEqual(701)
})

test("utilisation", () => {
    expect(strat.utilizedCapacity(
        [vlan(0), vlan(1), vlan(1000)],
        1))
        .toStrictEqual(4)
})

function vlan(vlan) {
    return {"Properties": {"vlan": vlan}}
}

function vlans(from, to) {
    let vlans = []
    for (let i = from; i <= to; i++) {
        vlans.push(vlan(i))
    }

    return vlans
}

function range(from, to) {
    return {"Properties": {"from": from, "to": to}}
}