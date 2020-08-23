test("ip allocation from prefix", () => {
    const strat = require('../ipv4_from_subnet');
    expect(strat.invokeWithParams([], {'ResourcePoolName': "10.0.0.0/8"}))
        .toStrictEqual({
            "ip": "10.0.0.1"
        });

    expect(strat.invokeWithParams(["", "", ""], {'ResourcePoolName': "10.0.0.0/8"}))
        .toStrictEqual({
            "ip": "10.0.0.4"
        });
})

test("ip allocation from small prefix", () => {
    const strat = require('../ipv4_from_subnet');
    expect(strat.invokeWithParams([], {'ResourcePoolName': "10.0.0.0/30"}))
        .toStrictEqual({
            "ip": "10.0.0.1"
        });

    expect(strat.invokeWithParams([""], {'ResourcePoolName': "10.0.0.0/30"}))
        .toStrictEqual({
            "ip": "10.0.0.2"
        });

    expect(strat.invokeWithParams(["", ""], {'ResourcePoolName': "10.0.0.0/30"}))
        .toBe(null)
})