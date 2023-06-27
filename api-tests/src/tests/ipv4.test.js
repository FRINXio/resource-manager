import {
    claimResource,
    getResourcePool,
    getPoolHierarchyPath,
    deleteResourcePool,
    claimResourceWithAltId, queryResourcesByAltId, queryResourcesByAltIdAndPoolId
} from '../graphql-queries.js';
import {
    createIpv4PrefixRootPool,
    createIpv4PrefixNestedPool,
    createSingletonIpv4PrefixNestedPool,
    createIpv4NestedPool,
    get2ChildrenIds,
    prepareIpv4Pool,
    allocateFromIPv4PoolSerially,
    allocateFromIPv4PoolParallelly,
    queryIPs,
    cleanup,
    createIpv4RootPool,
    createVlanRangeRootPool
} from '../test-helpers.js';

import tap from 'tap';
const test = tap.test;

test('create ipv4 prefix root pool', async (t) => {
    const pool = await createIpv4PrefixRootPool()
    t.ok(pool);
    t.equal(pool.PoolProperties['address'], '10.0.0.0')
    t.equal(pool.PoolProperties['prefix'], 8)
    await deleteResourcePool(pool.id);
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

    t.same(poolChildren, [pool21Id, pool22Id]);

    let resource11Id = (await claimResource(pool21Id, { desiredSize: 4194304 })).id;
    let resource12Id = (await claimResource(pool21Id, { desiredSize: 4194304 })).id;
    let resource21Id = (await claimResource(pool22Id, { desiredSize: 4194304 })).id;
    let resource22Id = (await claimResource(pool22Id, { desiredSize: 4194304 })).id;

    const pool31Id = await createSingletonIpv4PrefixNestedPool(resource11Id);
    const pool32Id = await createIpv4NestedPool(resource12Id);
    const pool33Id = await createIpv4NestedPool(resource21Id);
    const pool34Id = await createIpv4NestedPool(resource22Id);

    const pool21Children = await get2ChildrenIds(pool21Id);
    t.same(pool21Children, [pool31Id, pool32Id]);

    const pool22Children = await get2ChildrenIds(pool22Id);
    t.same(pool22Children, [pool33Id, pool34Id]);

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
    t.same(hierarchyPath34.map(it => it.id), [rootPoolId, pool22Id]);

    const hierarchyPathRoot = await getPoolHierarchyPath(rootPoolId);
    t.equal(hierarchyPathRoot.length, 0);

    const hierarchyPath21 = await getPoolHierarchyPath(pool21Id);
    t.same(hierarchyPath21.map(it => it.id), [rootPoolId]);

    await cleanup()
    t.end();
});

test('ipv4_prefix_pool serially', async (t) => {
    const count = 100;
    const poolId = (await createIpv4PrefixRootPool()).id;
    const createdIPs = await allocateFromIPv4PoolSerially(poolId, count, { desiredSize: 2 });
    const queriedIPs = await queryIPs(poolId, count);
    t.same(queriedIPs, createdIPs);
    t.equal(queriedIPs.length, count);

    await cleanup()
    t.end();
});

test('ipv4_pool serially', async (t) => {
    const count = 100;
    const poolId = await prepareIpv4Pool();
    const createdIPs = await allocateFromIPv4PoolSerially(poolId, count, {});
    const queriedIPs = await queryIPs(poolId, count);
    t.same(queriedIPs, createdIPs);
    t.equal(queriedIPs.length, count);

    await cleanup()
    t.end();
});

test('ipv4_prefix_pool parallelly', async (t) => {
    const count = 100, retries = 1000;
    const poolId = (await createIpv4PrefixRootPool()).id;
    const createdIPs = (await allocateFromIPv4PoolParallelly(
        poolId, count, retries, { desiredSize: 2 })).sort();
    const queriedIPs = (await queryIPs(poolId, count)).sort();
    t.same(queriedIPs, createdIPs);
    t.equal(queriedIPs.length, count);

    await cleanup()
    t.end();
});

test('ipv4_pool parallelly', async (t) => {
    const count = 100, retries = 1000;
    const poolId = await prepareIpv4Pool();
    const createdIPs = (await allocateFromIPv4PoolParallelly(
        poolId, count, retries, {})).sort();
    const queriedIPs = (await queryIPs(poolId, count)).sort();
    t.same(queriedIPs, createdIPs);
    t.equal(queriedIPs.length, count);

    await cleanup()
    t.end();
});

test('allocate resource from ipv4 prefix pool with desired value', async (t) => {
    const pool = await createIpv4PrefixRootPool();
    const allocatedResource = await claimResource(pool.id, {desiredSize: 6, desiredValue: "10.0.0.16"});

    t.same(allocatedResource.Properties, {address: "10.0.0.16", prefix: 29, subnet: false});

    console.log(allocatedResource.Properties);

    await cleanup();
    t.end();
});

test('insufficient capacity to claim resource in ipv4 prefix pool', async (t) => {
    const pool = await createIpv4PrefixRootPool();
    await claimResource(pool.id, {desiredSize: 64, desiredValue: "10.0.0.128"});
    const allocatedResource = await claimResource(pool.id, {desiredSize: 128, desiredValue: "10.0.0.128"});

    t.notOk(allocatedResource);

    await cleanup();
    t.end();
});

test('invalid desired value when claiming resource from ipv4 prefix pool', async (t) => {
    const pool = await createIpv4PrefixRootPool();
    const allocatedResource = await claimResource(pool.id, {desiredSize: 128, desiredValue: "10.0.0.53"});

    t.notOk(allocatedResource);

    await cleanup();
    t.end();
});

test('overlapping subnet with provided desired value', async (t) => {
    const pool = await createIpv4PrefixRootPool();
    await claimResource(pool.id, {desiredSize: 30, desiredValue: "10.0.0.0"});
    await claimResource(pool.id, {desiredSize: 6, desiredValue: "10.0.0.24"});
    const allocatedResource = await claimResource(pool.id, {desiredSize: 14, desiredValue: "10.0.0.16"});

    t.notOk(allocatedResource);

    await cleanup();
    t.end();
});

test('cannot create pool with subnet true and 31 or 32 prefix', async (t) => {
    const pool31 = await createIpv4RootPool("10.0.0.0", 31, true);
    const pool32 = await createIpv4RootPool("10.0.0.0", 32, true);

    t.equal(pool31, null);
    t.equal(pool32, null);

    await cleanup();
    t.end();
});

test('search resources by Optional<poolId> and altId', async (t) => {
    const pool1 = await createIpv4PrefixRootPool("10.0.0.0", 24, true);
    const pool2 = await createIpv4PrefixRootPool("10.0.0.0", 24, true);

    await claimResourceWithAltId(pool1.id, {
        desiredSize: 2,
    }, {
        altId: "altId1",
    })

    await claimResourceWithAltId(pool1.id, {
        desiredSize: 2,
    }, {
        altId: "altId1",
    })

    await claimResourceWithAltId(pool1.id, {
        desiredSize: 2,
    }, {
        altId: "altId1",
    })

    await claimResourceWithAltId(pool2.id, {
        desiredSize: 2,
    }, {
        altId: "altId1",
    })

    await claimResourceWithAltId(pool2.id, {
        desiredSize: 2,
    }, {
        altId: "altId1",
    })

    await claimResourceWithAltId(pool2.id, {
        desiredSize: 2,
    }, {
        altId: "altId1",
    })

    const resources = await queryResourcesByAltId({
        altId: "altId1",
    });

    const resources2 = await queryResourcesByAltIdAndPoolId(pool2.id, {
        altId: "altId1",
    });

    t.equal(resources.edges.length, 6);
    t.equal(resources2.edges.length, 3);

    await cleanup();
    t.end();
});

test('Claim resource with desired size of 1 and subnet false', async (t) => {
    const pool = await createIpv4PrefixRootPool("10.0.0.0", 24, false);

    const resource = await claimResource(pool.id, {
        desiredSize: 1,
    });

    t.equal(resource.Properties.address, "10.0.0.0");
    t.equal(resource.Properties.prefix, 32);

    await cleanup();
    t.end();
});

test('Claim resource with desired size of 1 and subnet true', async (t) => {
    const pool = await createIpv4PrefixRootPool("10.0.0.0", 24, true);

    const resource = await claimResource(pool.id, {
        desiredSize: 1,
    });

    t.equal(resource.Properties.address, "10.0.0.0");
    t.equal(resource.Properties.prefix, 30);

    await cleanup();
    t.end();
});

test('Claim resource with desired size of -1 and subnet false', async (t) => {
    const pool = await createIpv4PrefixRootPool("10.0.0.0", 24, false);

    const resource = await claimResource(pool.id, {
        desiredSize: -1,
    });

    t.notOk(resource);

    await cleanup();
    t.end();
});

test('Claim resource with desired size of 0 and subnet true', async (t) => {
    const pool = await createIpv4PrefixRootPool("10.0.0.0", 24, true);

    const resource = await claimResource(pool.id, {
        desiredSize: 0,
    });

    t.notOk(resource);

    await cleanup();
    t.end();
});