const strat = require('../formater_counter_strategy')

test("formatter output and capacity", () => {

    let printFun = (userInput, counter, resourcePoolProperties) => {
        return `VPN-${counter}-${userInput.someProperty}-${resourcePoolProperties.someProperty}-local`;
    }

    let allocated = [
        {text: 'first', counter: 0},
        {text: 'second', counter: 1},
        {text: 'third', counter: 3}];

    let output = strat.invokeWithParams(
        allocated,
        {someProperty: 'VPN85'},
        {someProperty: 'Network19', textFunction: printFun});

    expect(output)
        .toStrictEqual({'counter': 4, text: 'VPN-4-Network19-VPN85-local'});

    let capacity = strat.invokeWithParamsCapacity( allocated,
        {someProperty: 'VPN85'},
        {someProperty: 'Network19', textFunction: printFun});

    expect(capacity)
        .toStrictEqual({freeCapacity: Number.MAX_SAFE_INTEGER - 3, utilizedCapacity: 3});
})
