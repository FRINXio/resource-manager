test("prefix splitting", () => {
    const strat = require('../ipv4_prefix_splin_in_2_strategy');
    console.log(strat.invokeWithParams([], {'ResourcePoolName': "10.0.0.0/8"}));
})