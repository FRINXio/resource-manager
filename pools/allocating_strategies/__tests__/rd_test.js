const strat = require('../route_distinguisher_strategy')

test("allocate rd AS2", () => {
    for (let i = 0; i < 100; i++) {
        let randomAS = getRandomInt(1, 65000)
        for (let i = 0; i < 100; i++) {
            let randomNumber = getRandomInt(1, 4294967296)
            let rd = strat.invokeWithParams([],
                {'ResourcePoolName': "T1"},
                {"asNumber": randomAS, "assignedNumber": randomNumber}).rd
            expect(rd).toStrictEqual(randomAS + ":" + randomNumber)
        }
    }
})

test("allocate rd AS4", () => {
    for (let i = 0; i < 100; i++) {
        let randomAS = getRandomInt(1, 4294967296)
        for (let i = 0; i < 100; i++) {
            let randomNumber = getRandomInt(1, 65000)
            let rd = strat.invokeWithParams([],
                {'ResourcePoolName': "T2"},
                {"asNumber": randomAS, "assignedNumber": randomNumber}).rd
            expect(rd).toStrictEqual(randomAS + ":" + randomNumber)
        }
    }
})

test("allocate rd ipv4", () => {
    for (let i = 0; i < 100; i++) {
        let b1 = getRandomInt(0, 255)
        let b2 = getRandomInt(0, 255)
        let ipv4 = "1.2." + b1 + "." + b2
        for (let i = 0; i < 100; i++) {
            let randomNumber = getRandomInt(1, 65000)
            let rd = strat.invokeWithParams([],
                {'ResourcePoolName': "T2"},
                {"ipv4": ipv4, "assignedNumber": randomNumber}).rd
            expect(rd).toStrictEqual(ipv4 + ":" + randomNumber)
        }
    }
})

test("allocate rd wrong input", () => {
    let resourcePoolArg = {'ResourcePoolName': "T2"}
    expect(strat.invokeWithParams([], resourcePoolArg,
        {"assignedNumber": 1})).
    toStrictEqual(null)
    expect(strat.invokeWithParams([], resourcePoolArg,
        {"ipv4": "1.2.3.4"})).
    toStrictEqual(null)
    expect(strat.invokeWithParams([], resourcePoolArg,
        {"ipv4": "abcd", "assignedNumber": 1})).
    toStrictEqual(null)
    expect(strat.invokeWithParams([], resourcePoolArg,
        {"ipv4": "1.2.3.4", "assignedNumber": "asdasd"})).
    toStrictEqual(null)
    expect(strat.invokeWithParams([], resourcePoolArg,
        {"ipv4": "256.2.2.2", "assignedNumber": 1})).
    toStrictEqual(null)
    expect(strat.invokeWithParams([], resourcePoolArg,
        {"asNumber": 650000, "assignedNumber": 6500000})).
    toStrictEqual(null)
    expect(strat.invokeWithParams([], resourcePoolArg,
        {"ipv4": "22.2.2.2", "assignedNumber": 6500000})).
    toStrictEqual(null)
    expect(strat.invokeWithParams([], resourcePoolArg,
        {})).
    toStrictEqual(null)
})

function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 0.9999) + min);
}