import {claimResource, deleteResourcePool, getAllPoolsByTypeOrTag, getCapacityForPool} from '../graphql-queries.js';
import { cleanup, createUniqueIdPool } from '../test-helpers.js';

import tap from 'tap';
const test = tap.test;

test('create unique_id pool', async (t) => {
    const pool = await createUniqueIdPool()
    t.ok(pool);
    t.equal(pool.PoolProperties['from'], 1)
    t.equal(pool.PoolProperties['to'], 15)
    t.equal(pool.PoolProperties['idFormat'], "{counter}")

    await deleteResourcePool(pool.id);
    t.end();
});

test('create unique_id pool and allocate resources', async (t) => {
    const pool = await createUniqueIdPool()
    t.ok(pool);

    let resource1 = await claimResource(pool.id, {desiredValue: 4}, "first");
    let resource2 = await claimResource(pool.id, {}, "second");
    let resource3 = await claimResource(pool.id, {desiredValue: 14}, "third");
    let resource4 = await claimResource(pool.id, {desiredValue: 16}, "value is out of scope");
    let resource5 = await claimResource(pool.id, {desiredValue: 4}, "unique-id was already claimed");

    t.equal(resource1.Properties.counter, 4);
    t.equal(resource1.Properties.text, "4");
    t.equal(resource2.Properties.counter, 1);
    t.equal(resource2.Properties.text, "1");
    t.equal(resource3.Properties.counter, 14);
    t.equal(resource3.Properties.text, "14");

    t.notOk(resource4)
    t.notOk(resource5)

    await cleanup()
    t.end();
});

test('unique_id pool capacity', async (t) => {
    const pool = await createUniqueIdPool()
    t.ok(pool);

    await claimResource(pool.id, {}, "");
    await claimResource(pool.id, {}, "");
    await claimResource(pool.id, {}, "");
    let capacity = await getCapacityForPool(pool.id);
    t.equal(capacity.utilizedCapacity, "3");
    t.equal(capacity.freeCapacity, "12");

    await claimResource(pool.id, {}, "");
    await claimResource(pool.id, {}, "");
    capacity = await getCapacityForPool(pool.id);
    t.equal(capacity.utilizedCapacity, "5");
    t.equal(capacity.freeCapacity, "10");

    for (let i = 1; i <= 10; i++) {
        await claimResource(pool.id, {}, "");
    }
    capacity = await getCapacityForPool(pool.id);
    t.equal(capacity.utilizedCapacity, "15");
    t.equal(capacity.freeCapacity, "0");

    let resource16 = await claimResource(pool.id, {}, "Unique-id pool is full");
    t.notOk(resource16)

    await cleanup()
    t.end();
});

test('test loading multiple unique id pools', async (t) => {
    for (let i = 1; i <= 100; i++) {
        await createUniqueIdPool();
    }

    const pools = await getAllPoolsByTypeOrTag();
    t.equal(pools.edges.length, 100);

    await cleanup();
    t.end();
});
