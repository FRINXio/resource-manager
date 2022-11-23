import tap from 'tap';
import {claimResource} from "../graphql-queries.js";
import {cleanup, createRandomSignedInt32Pool} from "../test-helpers.js";
const test = tap.test;

test('claim resource from only positive random int32 values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: 1, to: 100});

    const resources = await claimResource(poolId, {});

    t.equal(resources.length, 1);

    await cleanup();
    t.end();
});

test('claim resource from only negative random int32 values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: -100, to: -1});

    const resources = await claimResource(poolId, {});

    t.equal(resources.length, 1);

    await cleanup();
    t.end();
});

test('claim resource from mixed random int32 values', async (t) => {
    const poolId = await createRandomSignedInt32Pool({from: -100, to: 100});

    const resources = await claimResource(poolId, {});

    t.equal(resources.length, 1);

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
    const poolId = await createRandomSignedInt32Pool({from: -2147483650, to: 2147483650});

    t.notOk(poolId);

    await cleanup();
    t.end();
});