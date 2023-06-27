import { claimResource, queryResource } from '../graphql-queries.js';
import {
    cleanup,
    createIpv6NestedPool,
    createIpv6PrefixRootPool,
    get2ChildrenIds
} from '../test-helpers.js';
import tap from 'tap';
const test = tap.test;


test('Claim resource with desired size of 1 and subnet false', async (t) => {
    const poolId = await createIpv6PrefixRootPool("dead::", 64, false);

    const resource = await claimResource(poolId, {
        desiredSize: 1,
    });

    t.equal(resource.Properties.address, "dead::");
    t.equal(resource.Properties.prefix, 128);

    await cleanup();
    t.end();
});

test('Claim resource with desired size of 1 and subnet true', async (t) => {
    const poolId = await createIpv6PrefixRootPool("dead::", 64, true);

    const resource = await claimResource(poolId, {
        desiredSize: 1,
    });

    t.equal(resource.Properties.address, "dead::");
    t.equal(resource.Properties.prefix, 126);

    await cleanup();
    t.end();
});

test('Claim resource with desired size of -1 and subnet false', async (t) => {
    const poolId = await createIpv6PrefixRootPool("dead::", 64, false);

    const resource = await claimResource(poolId, {
        desiredSize: -1,
    });

    t.notOk(resource);

    await cleanup();
    t.end();
});

test('Claim resource with desired size of 0 and subnet true', async (t) => {
    const poolId = await createIpv6PrefixRootPool("dead::", 64, true);

    const resource = await claimResource(poolId, {
        desiredSize: 0,
    });

    t.notOk(resource);

    await cleanup();
    t.end();
});

test('create ipv6 root pool', async (t) => {
    t.ok(await createIpv6PrefixRootPool());

    await cleanup()
    t.end();
});

test('create ipv6 hierarchy', async (t) => {
    let rootPoolId = await createIpv6PrefixRootPool();
    let firstResource = await claimResource(rootPoolId, {desiredSize: 128}, "first");
    let secondResource = await claimResource(rootPoolId, {desiredSize: 128});
    let nestedPool1Id = await createIpv6NestedPool(firstResource.id);
    let nestedPool2Id = await createIpv6NestedPool(secondResource.id);

    const children = await get2ChildrenIds(rootPoolId);
    t.same(children, [nestedPool1Id, nestedPool2Id]);

    let resource1 = await claimResource(nestedPool1Id, {});
    let resource2 = await claimResource(nestedPool2Id, {});

    t.equal(resource1.Properties.address, 'dead::');
    t.equal(resource2.Properties.address, 'dead::80');

    // assert resource query
    t.equal((await queryResource(rootPoolId, firstResource.Properties)).Description, "first");
    t.equal((await queryResource(rootPoolId, secondResource.Properties)).Description, null);

    await cleanup()
    t.end();
});

test('cannot create pool with subnet true and 127 or 128 prefix', async (t) => {
    const pool127 = await createIpv6PrefixRootPool("2001:db8:1::", 127, true);
    const pool128 = await createIpv6PrefixRootPool("2001:db8:1::", 128, true);

    t.notOk(pool127);
    t.notOk(pool128);

    await cleanup();
    t.end();
});
