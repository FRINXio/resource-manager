import tap from 'tap';
import {claimResource} from "../graphql-queries.js";
import {cleanup, createRandomSignedInt32Pool} from "../test-helpers.js";
const test = tap.test;

test('claim resource from only positive random int32 values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: 1, to: 100});

    const resource = await claimResource(poolId, {});

    t.not(resource.Properties.int, null);

    await cleanup();
    t.end();
});

test('claim resource from only negative random int32 values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: -100, to: -1});

    const resource = await claimResource(poolId, {});

    t.not(resource.Properties.int, null);

    await cleanup();
    t.end();
});

test('claim resource from mixed random int32 values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: -100, to: 100});

    const resource = await claimResource(poolId, {});

    t.not(resource.Properties.int, null);

    await cleanup();
    t.end();
});

test('create random signed int32 pool with maximum values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: -2147483648, to: 2147483647});

    t.ok(poolId);

    await cleanup();
    t.end();
});

test('create random signed int32 pool with bad values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: -121897391898798798928374923874928374928374, to: 2147483650});

    t.notOk(poolId);

    await cleanup();
    t.end();
});