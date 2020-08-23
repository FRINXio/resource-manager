test("ip allocation from small prefix", () => {
    const strat = require('../ipv4_from_pool');
    const resourcePoolArg = {'ResourcePoolName': "10.0.0.0/30"};

    expect(strat.invokeWithParams([], resourcePoolArg))
        .toStrictEqual({
            "ip": "10.0.0.0"
        });

    expect(strat.invokeWithParams([""], resourcePoolArg))
        .toStrictEqual({
            "ip": "10.0.0.1"
        });

    expect(strat.invokeWithParams(["", ""], resourcePoolArg))
        .toStrictEqual({
            "ip": "10.0.0.2"
        });

    expect(strat.invokeWithParams(["", "", ""], resourcePoolArg))
        .toStrictEqual({
            "ip": "10.0.0.3"
        });

    expect(strat.invokeWithParams(["", "", "", ""], resourcePoolArg))
        .toBe(null)
})