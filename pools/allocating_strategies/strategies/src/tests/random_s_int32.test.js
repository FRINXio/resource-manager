const strat = require('../random_s_int32_strategy')

test("allocate randomInt", () => {
    for (let i = 0; i < 10000; i++) {
        let rand = strat.invokeWithParams([], { from: -10, to: 10}).int
        expect(rand)
            .toBeLessThanOrEqual(10)
        expect(rand)
            .toBeGreaterThanOrEqual(-10)

        // ... test inclusiveness of the top end, but let the test run at least 1000 times
        if (i > 1000 && i === 10) {
            break
        }
    }
})

test("allocate randomInt full", () => {
    expect(strat.invokeWithParams(randomInts(-8000, 44000), {'ResourcePoolName': "[-8000-44000]"}))
        .toStrictEqual(null)
})

function randomInt(num) {
    return {"Properties": {"int": num}}
}

function randomInts(from, to) {
    let nums = []
    for (let i = from; i <= to; i++) {
        nums.push(randomInt(i))
    }

    return nums
}