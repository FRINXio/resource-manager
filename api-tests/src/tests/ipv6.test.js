import {claimResource} from '../graphql-queries';
import {createIpv6PrefixRootPool, createIpv6NestedPool, get2ChildrenIds} from '../test-helpers';

test('create ipv6 root pool', async () => {
    expect(await createIpv6PrefixRootPool()).toBeTruthy();
});

test('create ipv6 hierarchy', async () => {
    let rootPoolId = await createIpv6PrefixRootPool();
    let firstParentResourceId = (await claimResource(rootPoolId, {desiredSize: 128})).id;
    let secondParentResourceId = (await claimResource(rootPoolId, {desiredSize: 128})).id;
    let nestedPool1Id = await createIpv6NestedPool(firstParentResourceId);
    let nestedPool2Id = await createIpv6NestedPool(secondParentResourceId);

    const children = await get2ChildrenIds(rootPoolId);
    expect(children).toHaveLength(2);
    expect(children).toContain(nestedPool1Id);
    expect(children).toContain(nestedPool2Id);

    let resource1 = await claimResource(nestedPool1Id, {});
    let resource2 = await claimResource(nestedPool2Id, {});

    expect(resource1.Properties.address).toBe('dead::');
    expect(resource2.Properties.address).toBe('dead::80');
});
