import { executeBenchmark, record } from './util/bench.js';
import { testStrategy, createAllocationStrategy } from '../graphql-queries.js';
import { getUniqueName } from '../test-helpers.js';

async function initAndGetJsTestFunction(histograms) {
    return record(histograms, 'setup', async () => {
        let poolName = getUniqueName('testJSstrategy');
        let strategyId = await createAllocationStrategy(
            poolName,
            'function invoke() {return {vlan: userInput.desiredVlan};}',
            'js');
        return async () => record(histograms, 'testStrategy', async () =>
            await testStrategy(strategyId, { ResourcePoolName: 'testpool' }, poolName, [], { desiredVlan: 85 }));
    });
}

async function initAndGetPyTestFunction(histograms) {
    return record(histograms, 'setup', async () => {
        let poolName = getUniqueName('testPYstrategy');
        let strategyId = await createAllocationStrategy(
            poolName,
            'return {\'vlan\': userInput[\'desiredVlan\']}',
            'py');
        return async () => record(histograms, 'testStrategy', async () =>
            await testStrategy(strategyId, { ResourcePoolName: poolName }, poolName, [], { desiredVlan: 85 }));
    });
}

async function executeInParallel(histograms, fn, iterations) {
    const promises = [];
    for (let i = 0; i < iterations; i++) {
        promises.push(fn());
    }
    const results = await record(histograms, 'awaitPromisses', () => Promise.all(promises));
    for (const strategyOutput of results) {
        if (strategyOutput.stdout.vlan != 85) {
            throw new Error('Unexpected vlan' + strategyOutput.stdout.vlan);
        }
    }
}

async function runAll() {

    await executeBenchmark('create and test simple JS strategy once', {},
        async (histograms) => {
            const fn = await initAndGetJsTestFunction(histograms);
            let strategyOutput = await fn();
            if (strategyOutput.stdout.vlan != 85) {
                throw new Error('Unexpected vlan' + strategyOutput.stdout.vlan);
            }
        });

    await executeBenchmark('create and test simple PY strategy once', {},
        async (histograms) => {
            const fn = await initAndGetPyTestFunction(histograms);
            let strategyOutput = await fn();
            if (strategyOutput.stdout.vlan != 85) {
                throw new Error('Unexpected vlan' + strategyOutput.stdout.vlan);
            }
        });

    await executeBenchmark('create and test simple JS strategy 100x', {},
        async (histograms) => {
            const fn = await initAndGetJsTestFunction(histograms);
            await executeInParallel(histograms, fn, 100);
        });

    await executeBenchmark('create and test simple PY strategy 100x', {},
        async (histograms) => {
            const fn = await initAndGetPyTestFunction(histograms);
            await executeInParallel(histograms, fn, 100);
        });
}

runAll();
