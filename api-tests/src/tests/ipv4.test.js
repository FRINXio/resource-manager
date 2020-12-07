import { claimResource, getResourcePool, getPoolHierarchyPath } from '../graphql-queries.js';
import {
    createIpv4PrefixRootPool, createIpv4PrefixNestedPool,
    createSingletonIpv4PrefixNestedPool,
    createIpv4NestedPool, get2ChildrenIds,
    prepareIpv4Pool, allocateFromIPv4PoolSerially, allocateFromIPv4PoolParallelly, queryIPs
} from '../test-helpers.js';

import tap from 'tap';
const test = tap.test;

test('create ipv4 prefix root pool', async (t) => {
    const pool = await createIpv4PrefixRootPool()
    t.ok(pool);
    t.equal(pool.PoolProperties['address'], '10.0.0.0')
    t.equal(pool.PoolProperties['prefix'], 8)
    t.end();
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

test('create ipv4 hierarchy', async (t) => {
    let rootPoolId = (await createIpv4PrefixRootPool()).id;

    let firstParentResourceId = (await claimResource(rootPoolId, { desiredSize: 8388608 })).id;
    let secondParentResourceId = (await claimResource(rootPoolId, { desiredSize: 8388608 })).id;

    let pool21Id = await createIpv4PrefixNestedPool(firstParentResourceId);
    let pool22Id = await createIpv4PrefixNestedPool(secondParentResourceId);

    const poolChildren = await get2ChildrenIds(rootPoolId);

    t.deepEqual(poolChildren, [pool21Id, pool22Id]);

    let resource11Id = (await claimResource(pool21Id, { desiredSize: 4194304 })).id;
    let resource12Id = (await claimResource(pool21Id, { desiredSize: 4194304 })).id;
    let resource21Id = (await claimResource(pool22Id, { desiredSize: 4194304 })).id;
    let resource22Id = (await claimResource(pool22Id, { desiredSize: 4194304 })).id;

    const pool31Id = await createSingletonIpv4PrefixNestedPool(resource11Id);
    const pool32Id = await createIpv4NestedPool(resource12Id);
    const pool33Id = await createIpv4NestedPool(resource21Id);
    const pool34Id = await createIpv4NestedPool(resource22Id);

    const pool21Children = await get2ChildrenIds(pool21Id);
    t.deepEqual(pool21Children, [pool31Id, pool32Id]);

    const pool22Children = await get2ChildrenIds(pool22Id);
    t.deepEqual(pool22Children, [pool33Id, pool34Id]);

    let resource31 = await claimResource(pool31Id, {});
    let resource32 = await claimResource(pool32Id, {});
    let resource33 = await claimResource(pool33Id, {});
    let resource34 = await claimResource(pool34Id, {});

    t.equal(resource31.Properties.address, '10.10.0.0');
    t.equal(resource32.Properties.address, '10.64.0.0');
    t.equal(resource33.Properties.address, '10.128.0.0');
    t.equal(resource34.Properties.address, '10.192.0.0');

    // assert hierarchy queries
    const pool34Queried = await getResourcePool(pool34Id);
    t.equal(pool34Queried.ParentResource.id, resource22Id);
    t.equal(pool34Queried.ParentResource.ParentPool.id, pool22Id);

    const hierarchyPath34 = await getPoolHierarchyPath(pool34Id);
    t.deepEqual(hierarchyPath34.map(it => it.id), [rootPoolId, pool22Id]);

    const hierarchyPathRoot = await getPoolHierarchyPath(rootPoolId);
    t.equal(hierarchyPathRoot.length, 0);

    const hierarchyPath21 = await getPoolHierarchyPath(pool21Id);
    t.deepEqual(hierarchyPath21.map(it => it.id), [rootPoolId]);

    t.end();
});

test('ipv4_pool serially', async (t) => {
    const count = 100;
    const poolId = await prepareIpv4Pool();
    const createdIPs = await allocateFromIPv4PoolSerially(poolId, count);
    const queriedIPs = await queryIPs(poolId, count);
    t.deepEqual(queriedIPs, createdIPs);
    t.equal(queriedIPs.length, count);
    t.end();
});

test('ipv4_pool parallelly', async (t) => {
    const count = 100, retries = 1000;
    const poolId = await prepareIpv4Pool();
    const createdIPs = (await allocateFromIPv4PoolParallelly(poolId, count, retries)).sort();
    const queriedIPs = (await queryIPs(poolId, count)).sort();
    t.deepEqual(queriedIPs, createdIPs);
    t.equal(queriedIPs.length, count);
    t.end();
});
