import { claimResource} from '../graphql-queries.js';
import {createVlanRangeRootPool, createVlanNestedPool, get2ChildrenIds} from '../test-helpers.js';
import tap from 'tap';
const test = tap.test;

test('create vlan root pool', async (t) => {
    t.ok(await createVlanRangeRootPool());
    t.end();
});

//             vlan hierarchy
//
//               [0-4095]
//                   |
//          [0-2000]   [2001-4095]

test('create vlan hierarchy', async (t) => {
    let rootPoolId = await createVlanRangeRootPool();
    let firstParentResourceId = (await claimResource(rootPoolId, {desiredSize: 2001})).id;
    let secondParentResourceId = (await claimResource(rootPoolId, {desiredSize: 2095})).id;
    let nestedPool1Id = await createVlanNestedPool(firstParentResourceId);
    let nestedPool2Id = await createVlanNestedPool(secondParentResourceId);

    const children = await get2ChildrenIds(rootPoolId);
    t.deepEqual(children, [nestedPool1Id, nestedPool2Id]);

    let resource1 = await claimResource(nestedPool1Id, {});
    let resource2 = await claimResource(nestedPool2Id, {});

    t.equal(resource1.Properties.vlan, 0);
    t.equal(resource2.Properties.vlan, 2001);
    t.end();
});
