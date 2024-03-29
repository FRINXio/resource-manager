import { record, bench } from './util/bench.js';
import {
    createIpv4PrefixRootPool, prepareIpv4Pool,
    allocateFromIPv4PoolSerially, allocateFromIPv4PoolParallelly,
    getUniqueName,
} from '../test-helpers.js';
import {
    findResourceTypeId, findAllocationStrategyId, createAllocationPool,
} from "../graphql-queries.js";

bench('Create ipv4 prefix pool 100x serially',
    async (histograms) => {
        const count = 100;
        let resourceTypeId, strategyId;
        await record(histograms, 'setup',
            async () => {
                resourceTypeId = await findResourceTypeId('ipv4_prefix');
                strategyId = await findAllocationStrategyId('ipv4_prefix');
            });
        await record(histograms, 'createAllocationPool', async () => {
            for (let i = 0; i < count; i++) {
                await createAllocationPool(
                    getUniqueName('ipv4-root'),
                    resourceTypeId,
                    strategyId,
                    { prefix: "int", address: "string", subnet: "bool"},
                    { prefix: 8, address: "10.0.0.0", subnet: false},
                    null,
                    false);
            }
        });
    }
);

bench('Create ipv4 prefix pool 100x parallelly',
    async (histograms) => {
        let resourceTypeId, strategyId;
        await record(histograms, 'setup',
            async () => {
                resourceTypeId = await findResourceTypeId('ipv4_prefix');
                strategyId = await findAllocationStrategyId('ipv4_prefix');
            });
        const count = 100;
        const getPromises = (count) => {
            const ppResult = [];
            for (let i = 0; i < count; i++) {
                ppResult.push(
                    createAllocationPool(
                        getUniqueName('ipv4-root'),
                        resourceTypeId,
                        strategyId,
                        { prefix: "int", address: "string", subnet: "bool"},
                        { prefix: 8, address: "10.0.0.0", subnet: false},
                        null,
                        true));
            }
            return ppResult;
        };
        let result = [];
        while (result.length < count) {
            // Some promises might be rejected, loop until we have `count` items created.
            result.push(...getPromises(count - result.length));
            result = (await Promise.all(result)).filter(it => it);
        }
    }
);

bench('allocate 100 ipv4_prefix_pool resources serially',
    async (histograms) => {
        const iterations = 100;
        let poolId = await record(histograms, 'setup', async () => (await createIpv4PrefixRootPool()).id);
        await record(histograms, 'allocate',
            async () => await allocateFromIPv4PoolSerially(poolId, iterations, { desiredSize: 2 }));
    }
);

bench('allocate 100 ipv4_prefix_pool resources parallelly',
    async (histograms) => {
        const iterations = 100;
        let poolId = await record(histograms, 'setup', async () => (await createIpv4PrefixRootPool()).id);
        await record(histograms, 'awaitPromisses',
            async () => allocateFromIPv4PoolParallelly(poolId, iterations, iterations * 10, { desiredSize: 2 }));
    }
);

bench('allocate 100 ipv4_pool resources serially',
    async (histograms) => {
        const iterations = 100;
        const poolId = await record(histograms, 'setup', prepareIpv4Pool);
        await record(histograms, 'allocate',
            async () => await allocateFromIPv4PoolSerially(poolId, iterations, {}));
    }
);

bench('allocate 100 ipv4_pool resources parallelly',
    async (histograms) => {
        const iterations = 100;
        const poolId = await record(histograms, 'setup', prepareIpv4Pool);
        await record(histograms, 'awaitPromisses',
            async () => allocateFromIPv4PoolParallelly(poolId, iterations, iterations * 10, {}));
    }
);
