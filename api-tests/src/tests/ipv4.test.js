import { claimResource, getResourcePool, getPoolHierarchyPath } from '../graphql-queries';
import {createIpv4PrefixRootPool, createIpv4PrefixNestedPool, createSingletonIpv4PrefixNestedPool,
createIpv4NestedPool, get2ChildrenIds} from '../test-helpers';

test('create ipv4 prefix root pool', async () => {
    expect(await createIpv4PrefixRootPool()).toBeTruthy();
});

//# Test pool hierarchies
//# Creates following hierarchy of ipv4 prefixes with ips at the end
//
//#                                           +---------------------------+
//#                                           |                           |
//#                                           |  Name: 10.0.0.0/8         |
//#                                           |  Type: Ipv4Prefix         |
//#                                           |  Allocation strategy:     |
//#                                           |   split_ipv4_prefix_in_2  |
//#                                           |                           |
//#                                           +--------------+------------+
//#                                                          |
//#                                                          |
//#                                +-------------------------+------------------------------+
//#                                |                                                        |
//#                                |                                                        |
//#                 +--------------+------------+                        +------------------+----------------+
//#                 |                           |                        |                                   |
//#                 |  Name: 10.0.0.0/9         |                        |  Name: 10.128.0.0/9               |
//#                 |  Type: Ipv4Prefix         |                        |  Type: Ipv4Prefix                 |
//#                 |  Allocation strategy:     |                        |  Set:                             |
//#                 |   split_ipv4_prefix_in_2  |                        |   [10.128.0.0/10, 10.192.0.0/10]  |
//#                 |                           |                        |                                   |
//#                 +-------------+-------------+                        +------------------+----------------+
//#                               |                                                         |
//#                               |                                                         |
//#              +----------------+----------------+                                    +---+----------------------------------+
//#              |                                 |                                    |                                      |
//#              |                                 |                                    |                                      |
//#+-------------+-------------+     +-------------+-------------+       +--------------+------------+           +-------------+-------------+
//#|                           |     |                           |       |                           |           |                           |
//#|  Name: 10.0.0.0/10        |     |  Name: 10.64.0.0/10       |       |  Name: 10.128.0.0/10      |           |  Name: 10.192.0.0/10      |
//#|  Type: Ipv4Prefix         |     |  Type: Ipv4               |       |  Type: Ipv4               |           |  Type: Ipv4               |
//#|  Singleton: 10.0.0.0/11   |     |  Allocation strategy:     |       |  Allocation strategy:     |           |  Allocation strategy:     |
//#|                           |     |   ipv4_from_subnet        |       |   ipv4_from_subnet        |           |   ipv4_from_subnet        |
//#|         # Unused          |     |                           |       |                           |           |                           |
//#+---------------------------+     +---------------------------+       +---------------------------+           +---------------------------+

test('create ipv4 hierarchy', async () => {
    let rootPoolId = await createIpv4PrefixRootPool();

    let firstParentResourceId = (await claimResource(rootPoolId, {desiredSize: 8388608})).id;
    let secondParentResourceId = (await claimResource(rootPoolId, {desiredSize: 8388608})).id;

    let pool21Id = await createIpv4PrefixNestedPool(firstParentResourceId);
    let pool22Id = await createIpv4PrefixNestedPool(secondParentResourceId);

    const poolChildren = await get2ChildrenIds(rootPoolId);

    expect(poolChildren).toHaveLength(2);
    expect(poolChildren).toContain(pool21Id);
    expect(poolChildren).toContain(pool22Id);

    let resource11Id = (await claimResource(pool21Id, {desiredSize: 4194304})).id;
    let resource12Id = (await claimResource(pool21Id, {desiredSize: 4194304})).id;
    let resource21Id = (await claimResource(pool22Id, {desiredSize: 4194304})).id;
    let resource22Id = (await claimResource(pool22Id, {desiredSize: 4194304})).id;

    const pool31Id = await createSingletonIpv4PrefixNestedPool(resource11Id);
    const pool32Id = await createIpv4NestedPool(resource12Id);
    const pool33Id = await createIpv4NestedPool(resource21Id);
    const pool34Id = await createIpv4NestedPool(resource22Id);

    const pool21Children = await get2ChildrenIds(pool21Id);
    expect(pool21Children).toHaveLength(2);
    expect(pool21Children).toContain(pool31Id);
    expect(pool21Children).toContain(pool32Id);

    const pool22Children = await get2ChildrenIds(pool22Id);
    expect(pool22Children).toHaveLength(2);
    expect(pool22Children).toContain(pool33Id);
    expect(pool22Children).toContain(pool34Id);

    let resource31 = await claimResource(pool31Id, {});
    let resource32 = await claimResource(pool32Id, {});
    let resource33 = await claimResource(pool33Id, {});
    let resource34 = await claimResource(pool34Id, {});

    expect(resource31.Properties.address).toBe('10.10.0.0');
    expect(resource32.Properties.address).toBe('10.64.0.0');
    expect(resource33.Properties.address).toBe('10.128.0.0');
    expect(resource34.Properties.address).toBe('10.192.0.0');

    // assert hierarchy queries
    const pool34Queried = await getResourcePool(pool34Id)
    expect(pool34Queried.ParentResource.id).toBe(resource22Id)
    expect(pool34Queried.ParentResource.ParentPool.id).toBe(pool22Id)

    const hierarchyPath34 = await getPoolHierarchyPath(pool34Id)
    expect(hierarchyPath34.length).toBe(2)
    expect(hierarchyPath34[0].id).toBe(rootPoolId)
    expect(hierarchyPath34[1].id).toBe(pool22Id)

    const hierarchyPathRoot = await getPoolHierarchyPath(rootPoolId)
    expect(hierarchyPathRoot.length).toBe(0)

    const hierarchyPath21 = await getPoolHierarchyPath(pool21Id)
    expect(hierarchyPath21.length).toBe(1)
    expect(hierarchyPath21[0].id).toBe(rootPoolId)
});
