const strat = require('../unique_id_strategy')

test("unique id output and capacity", () => {
    let allocated = [uniqueId(0, 'first'),
        uniqueId(1, 'second'),
        uniqueId(3, 'third')];

    let resource_pool = {
        vpn: "VPN85", network: "Network19",
        idFormat: "VPN-{counter}-{network}-{vpn}-local"}

    let output = strat.invokeWithParams(allocated, resource_pool, {});
    expect(output)
        .toStrictEqual({'counter': 4, text: 'VPN-4-Network19-VPN85-local'});

    let capacity = strat.invokeWithParamsCapacity(allocated, resource_pool, {});
    expect(capacity)
        .toStrictEqual({freeCapacity: Number.MAX_SAFE_INTEGER - 3, utilizedCapacity: 3});
})

test("params without resourcePool", () => {
    let output = strat.invokeWithParams([], null, {});
    expect(output).toStrictEqual(null);
})

test("resourcePool without idFormat", () => {
    let output = strat.invokeWithParams([],
        {someProperty: "SomeUniqueL3VPN"}, {});
    expect(output).toStrictEqual(null)
})

test("idFormat without counter", () => {
    let output = strat.invokeWithParams([],
        {someProperty: "SomeUniqueL3VPN", idFormat: "{someProperty}"}, {});
    expect(output).toStrictEqual(null)
})

test("simple l3vpn counter", () => {
    let output = strat.invokeWithParams([],
        {someProperty: "L3VPN", idFormat: "{someProperty}{counter}"}, {});
    expect(output).toStrictEqual({'counter': 1, text: 'L3VPN1'});

    let capacity = strat.invokeWithParamsCapacity([uniqueId(output.counter, output.text)],
        {someProperty: "L3VPN", idFormat: "{someProperty}{counter}"}, {});
    expect(capacity).toStrictEqual({freeCapacity: Number.MAX_SAFE_INTEGER - 1, utilizedCapacity: 1});

    let next_output = strat.invokeWithParams([uniqueId(output.counter, output.text)],
        {someProperty: "L3VPN", idFormat: "{someProperty}{counter}"}, {});
    expect(next_output).toStrictEqual({'counter': 2, text: 'L3VPN2'});

    let next_capacity = strat.invokeWithParamsCapacity(
        [uniqueId(output.counter, output.text), uniqueId(next_output.counter, next_output.text)],
        {someProperty: "L3VPN", idFormat: "{someProperty}{counter}"}, {});
    expect(next_capacity).toStrictEqual({freeCapacity: Number.MAX_SAFE_INTEGER - 2, utilizedCapacity: 2});
})

test("multiple l3vpn counters", () => {
    let outputs = []
    for (let i = 1; i <= 10; i++) {
        let unique_id = strat.invokeWithParams(outputs,
            {someProperty: "L3VPN", idFormat: "{someProperty}{counter}"},
            {});
        outputs.push(uniqueId(unique_id.counter, unique_id.text))
        expect(unique_id).toStrictEqual({'counter': i, text: 'L3VPN'+ i })
    }

    let capacity = strat.invokeWithParamsCapacity(outputs,
        {someProperty: "L3VPN", idFormat: "{someProperty}{counter}"},
        {});
    expect(capacity).toStrictEqual({freeCapacity: Number.MAX_SAFE_INTEGER - 10, utilizedCapacity: 10});
})

function uniqueId(counter, text) {
    return {"Properties": {"counter": counter, "text": text}}
}