import { claimResource} from '../graphql-queries';
import {createVlanRangeRootPool, createVlanNestedPool, get2ChildrenIds} from '../test-helpers';

test('create vlan root pool', async () => {
    expect(await createVlanRangeRootPool()).toBeTruthy();
});

//             vlan hierarchy
//
//               [0-4095]
//                   |
//          [0-2000]   [2001-4095]

test('create vlan hierarchy', async () => {
    let rootPoolId = await createVlanRangeRootPool();
    let firstParentResourceId = (await claimResource(rootPoolId, {desiredSize: 2001})).id;
    let secondParentResourceId = (await claimResource(rootPoolId, {desiredSize: 2095})).id;
    let nestedPool1Id = await createVlanNestedPool(firstParentResourceId);
    let nestedPool2Id = await createVlanNestedPool(secondParentResourceId);

    const children = await get2ChildrenIds(rootPoolId);
    expect(children).toHaveLength(2);
    expect(children).toContain(nestedPool1Id);
    expect(children).toContain(nestedPool2Id);

    let resource1 = await claimResource(nestedPool1Id, {});
    let resource2 = await claimResource(nestedPool2Id, {});

    expect(resource1.Properties.vlan).toBe(0);
    expect(resource2.Properties.vlan).toBe(2001);
});
