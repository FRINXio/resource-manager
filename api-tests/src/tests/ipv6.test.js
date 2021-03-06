import { claimResource, queryResource } from '../graphql-queries.js';
import { createIpv6NestedPool, createIpv6PrefixRootPool, get2ChildrenIds } from '../test-helpers.js';
import tap from 'tap';
const test = tap.test;

test('create ipv6 root pool', async (t) => {
    t.ok(await createIpv6PrefixRootPool());
    t.end();
});

test('create ipv6 hierarchy', async (t) => {
    let rootPoolId = await createIpv6PrefixRootPool();
    let firstResource = await claimResource(rootPoolId, {desiredSize: 128}, "first");
    let secondResource = await claimResource(rootPoolId, {desiredSize: 128});
    let nestedPool1Id = await createIpv6NestedPool(firstResource.id);
    let nestedPool2Id = await createIpv6NestedPool(secondResource.id);

    const children = await get2ChildrenIds(rootPoolId);
    t.deepEqual(children, [nestedPool1Id, nestedPool2Id]);

    let resource1 = await claimResource(nestedPool1Id, {});
    let resource2 = await claimResource(nestedPool2Id, {});

    t.equal(resource1.Properties.address, 'dead::');
    t.equal(resource2.Properties.address, 'dead::80');

    // assert resource query
    t.equal((await queryResource(rootPoolId, firstResource.Properties)).Description, "first");
    t.equal((await queryResource(rootPoolId, secondResource.Properties)).Description, null);
    t.end();
});
