import { claimResource, queryResource } from '../graphql-queries';
import { createIpv6NestedPool, createIpv6PrefixRootPool, get2ChildrenIds } from '../test-helpers';

test('create ipv6 root pool', async () => {
    expect(await createIpv6PrefixRootPool()).toBeTruthy();
});

test('create ipv6 hierarchy', async () => {
    let rootPoolId = await createIpv6PrefixRootPool();
    let firstResource = await claimResource(rootPoolId, {desiredSize: 128}, "first");
    let secondResource = await claimResource(rootPoolId, {desiredSize: 128});
    let nestedPool1Id = await createIpv6NestedPool(firstResource.id);
    let nestedPool2Id = await createIpv6NestedPool(secondResource.id);

    const children = await get2ChildrenIds(rootPoolId);
    expect(children).toHaveLength(2);
    expect(children).toContain(nestedPool1Id);
    expect(children).toContain(nestedPool2Id);

    let resource1 = await claimResource(nestedPool1Id, {});
    let resource2 = await claimResource(nestedPool2Id, {});

    expect(resource1.Properties.address).toBe('dead::');
    expect(resource2.Properties.address).toBe('dead::80');

    // assert resource query
    expect((await queryResource(rootPoolId, firstResource.Properties)).Description).toBe("first")
    expect((await queryResource(rootPoolId, secondResource.Properties)).Description).toBeNull()
});
