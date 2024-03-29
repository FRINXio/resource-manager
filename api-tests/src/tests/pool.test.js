import {
    claimResource,
    createSingletonPool,
    createTag,
    deleteResourcePool,
    findResourceTypeId,
    createSetPool,
    createResourceType,
    freeResource,
    claimResourceWithAltId,
    getResourcesForPool,
    searchPoolsByTags,
    tagPool,
    getCapacityForPool,
    getResourcePool,
    getResourcesByDatetimeRange,
    updateResourceAltId,
    getEmptyPools,
    getPaginatedResourcesForPool,
    getRequiredPoolProperties,
    findAllocationStrategyId,
    createAllocationPool,
    queryResourcesByAltIdAndPoolId, queryResourcePools
} from "../graphql-queries.js";
import {
    cleanup,
    createIpv4PrefixRootPool, createIpv4PrefixRootPoolWithName,
    createIpv4RootPool,
    createIpv6PrefixRootPool,
    createIpv6RootPool, createRandomIntRootPool, createRdRootPool,
    createVlanRangeRootPool,
    createVlanRootPool,
    getUniqueName, prepareDataForFiltering
} from "../test-helpers.js";
import tap from 'tap';
const test = tap.test;

test('singleton claim and free resource', async (t) => {
    let rtId = await findResourceTypeId('ipv4');
    const ipAddress = '192.168.1.1';
    const poolName = getUniqueName('singleton');
    let poolId = await createSingletonPool(poolName, rtId, [{address: ipAddress}]);
    let resourceTypeOfPoolName = await findResourceTypeId(poolName + "-ResourceType");
    t.equal(resourceTypeOfPoolName, undefined);

    let resource = await claimResource(poolId, {});
    let rs = await getResourcesForPool(poolId);
    t.equal(rs.edges.length, 1);
    t.equal(rs.edges[0].node.Properties.address, ipAddress)
    await freeResource(poolId, resource.Properties);

    rs = await getResourcesForPool(poolId);
    t.equal(rs.edges.length, 0);

    await cleanup()
    t.end();
});

test('create and delete singleton pool', async (t) => {
    let rtId = await findResourceTypeId('ipv4');
    let poolId = await createSingletonPool(
        getUniqueName('singleton'),
        rtId,
        [{address: '192.168.1.1'}]
    );

    const tagText = getUniqueName("singletonTag");
    const tagId = await createTag(tagText);
    await tagPool(tagId, poolId);
    let foundPool = await searchPoolsByTags({matchesAny: [{matchesAll: [tagText]}]});
    t.equal(foundPool?.edges?.length, 1);
    t.equal(foundPool?.edges?.[0].node.id, poolId);

    let resource1 = await claimResource(poolId, {});
    let resource2 = await claimResource(poolId, {});

    t.same(resource1, resource2); //the same resource

    await freeResource(poolId, resource2.Properties);

    await deleteResourcePool(poolId);
    foundPool = await searchPoolsByTags({matchesAny: [{matchesAll: [tagText]}]});
    t.equal(foundPool?.edges?.length, 0);

    await cleanup()
    t.end();
});

test('get pools by range of datetime', async (t) => {
    const poolId = await createIpv4RootPool('192.168.3.0', 16);

    await claimResource(poolId, {});
    await claimResource(poolId, {});

    let today = new Date().toISOString().slice(0, 13);
    today = today.replace("T", "-");
    let pools = await getResourcesByDatetimeRange(today, "");
    //from today 00:00 to current moment expected 2 resources
    t.equal(pools.edges.length, 2);

    let tomorrow = new Date()
    tomorrow.setDate(new Date().getDate() + 1)
    tomorrow = tomorrow.toISOString().slice(0,13)
    tomorrow = tomorrow.replace("T", "-");
    pools = await getResourcesByDatetimeRange(today, tomorrow);
    //from today 00:00 to tomorrow expected 2 resources
    t.equal(pools.edges.length, 2);

    let yesterday = new Date()
    yesterday.setDate(new Date().getDate() - 1)
    yesterday = yesterday.toISOString().slice(0,13)
    yesterday = yesterday.replace("T", "-");
    pools = await getResourcesByDatetimeRange(yesterday, today);
    t.equal(pools.edges.length, 0);

    await cleanup()
    t.end();
});

test('create and delete resources in set pool', async (t) => {
    let rtId = await findResourceTypeId('ipv4');
    let poolId = await createSetPool(
        getUniqueName('singleton'),
        rtId,
        [{address: '192.168.1.1'}, {address: '192.168.1.2'}]
    );
    let resource = await claimResource(poolId, {});
    let rs = await getResourcesForPool(poolId);
    t.equal(rs.edges.length, 1);
    await freeResource(poolId, resource.Properties)
    rs = await getResourcesForPool(poolId);
    t.equal(rs.edges.length, 0);

    await cleanup()
    t.end();
});

test('create and delete set pool', async (t) => {
    let rtId = await findResourceTypeId('ipv4');
    let poolId = await createSetPool(
        getUniqueName('singleton'),
        rtId,
        [{address: '192.168.1.1'}, {address: '192.168.1.2'}]
    );
    let tagText = getUniqueName("setTag");
    const tagId = await createTag(tagText);
    await tagPool(tagId, poolId);
    let foundPool = await searchPoolsByTags({matchesAny: [{matchesAll: [tagText]}]});
    t.equal(foundPool?.edges?.length, 1);
    t.equal(foundPool?.edges?.[0].node.id, poolId);
    let resource = await claimResource(poolId, {});
    await freeResource(poolId, resource.Properties)
    await deleteResourcePool(poolId);
    foundPool = await searchPoolsByTags({matchesAny: [{matchesAll: [tagText]}]});
    t.equal(foundPool?.edges?.length, 0);
    t.end();
});

test('allocating pool with incorrect properties', async (t) => {
    let resourceTypeId = await findResourceTypeId('ipv4');
    let strategyId = await findAllocationStrategyId('ipv4');
    let poolName = getUniqueName('root-ipv4');
    let pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {prefix: "int"},
        {prefix: 24},)
    t.notOk(pool)

    pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {address: "string"},
        {address: "192.168.3.0"},)
    t.notOk(pool)

    pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {address: "string", prefix: "int", subnet: "bool"},
        {address: "192.168.3.0", prefix: 24, subnet: true})
    t.ok(pool)

    await cleanup()
    t.end();
});

test('capacity for allocating vlan-range pool', async (t) => {
    const poolId = await createVlanRangeRootPool();

    await claimResource(poolId, {desiredSize: 1});
    await claimResource(poolId, {desiredSize: 1});
    await claimResource(poolId, {desiredSize: 3});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "5");
    t.equal(capacity.freeCapacity, "4091");

    await cleanup()
    t.end();
});

test('capacity for allocating vlan pool', async (t) => {
    const poolId = await createVlanRootPool();

    await claimResource(poolId, {});
    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "3");
    t.equal(capacity.freeCapacity, "4093");

    await cleanup()
    t.end();
});

test('capacity for allocating ipv6-prefix pool', async (t) => {
    const poolId = await createIpv6PrefixRootPool();

    await claimResource(poolId, {desiredSize: 4})
    await claimResource(poolId, {desiredSize: 4});
    await claimResource(poolId, {desiredSize: 4});
    await claimResource(poolId, {desiredSize: 4});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "16");
    t.equal(capacity.freeCapacity, "240");

    await cleanup()
    t.end();
});

test('capacity for allocating ipv6 pool', async (t) => {
    const poolId = await createIpv6RootPool();

    await claimResource(poolId, {})
    await claimResource(poolId, {});
    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "4");
    t.equal(capacity.freeCapacity, "5192296858534827628530496329220092");

    await cleanup()
    t.end();
});

test('capacity for allocating ipv4 pool', async (t) => {
    const poolId = await createIpv4RootPool('192.168.3.0', 16);

    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "2");
    t.equal(capacity.freeCapacity, "65534");

    await cleanup()
    t.end();
});

test('capacity for random pool', async (t) => {
    const poolId = await createRandomIntRootPool();

    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "2");
    t.equal(capacity.freeCapacity, "997");

    await cleanup()
    t.end();
});

test('capacity for allocating ipv4-prefix pool', async (t) => {
    const poolId = (await createIpv4PrefixRootPool()).id;

    await claimResource(poolId, {desiredSize: 2});
    await claimResource(poolId, {desiredSize: 2});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "4");
    t.equal(capacity.freeCapacity, "16777212");

    await cleanup()
    t.end();
});

test('capacity for allocating RD pool', async (t) => {
    const rdPoolId = await createRdRootPool();

    await claimResource(rdPoolId, {asNumber: 1985, assignedNumber: 5891});
    await claimResource(rdPoolId, {asNumber: 2020, assignedNumber: 2020});
    await claimResource(rdPoolId, {asNumber: 47, assignedNumber: 47});

    const capacity = await getCapacityForPool(rdPoolId);
    t.equal(capacity.utilizedCapacity, "3");
    t.equal(capacity.freeCapacity, "281474976710656");

    await cleanup()
    t.end();
});

test('capacity for set pool', async (t) => {
    let rtId = await findResourceTypeId('ipv4');
    let poolId = await createSetPool(
        getUniqueName('singleton-ipv4'),
        rtId,
        [{address: '192.168.1.1'}, {address: '192.168.1.2'}, {address: '192.168.1.3'}, {address: '192.168.1.4'}]
    );

    await claimResource(poolId, {});
    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "3");
    t.equal(capacity.freeCapacity, "1");

    await cleanup()
    t.end();
});

test('capacity for singleton pool', async (t) => {
    let rtId = await findResourceTypeId('ipv4');
    let poolId = await createSingletonPool(
        getUniqueName('singleton'),
        rtId,
        [{address: '192.168.1.1'}]
    );

    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    t.equal(capacity.utilizedCapacity, "1");
    t.equal(capacity.freeCapacity, "0");

    await cleanup()
    t.end();
});

test('pagination of allocated resources in vlan-pool', async (t) => {
    const poolId = await createVlanRootPool();

    let resourceIds = [];

    for (let i = 0; i < 20; i++) {
        resourceIds.push((await claimResource(poolId, {})).id);
    }

    //get 3 first resources
    let pool = await getResourcePool(poolId, null, null, 3, null);
    t.equal(pool.allocatedResources.edges.length, 3);
    let thirdResource = pool.allocatedResources.pageInfo.endCursor;

    //get 3 resources after the 3rd resource
    pool = await getResourcePool(poolId, null,  thirdResource, 3, null);
    t.equal(resourceIds[3], pool.allocatedResources.edges[0].node.id);
    t.equal(resourceIds[4], pool.allocatedResources.edges[1].node.id);
    t.equal(resourceIds[5], pool.allocatedResources.edges[2].node.id);
    t.ok(pool.allocatedResources.pageInfo.hasNextPage);

    //get all resources after the 3rd resource
    pool = await getResourcePool(poolId, null, thirdResource, 1000, null);
    t.equal(pool.allocatedResources.edges.length, 17);
    t.notOk(pool.allocatedResources.pageInfo.hasNextPage);

    //get 1 resource before the 3rd resource
    pool = await getResourcePool(poolId, thirdResource,  null, null, 1);
    t.equal(pool.allocatedResources.edges.length, 1);
    t.equal(resourceIds[1], pool.allocatedResources.edges[0].node.id);
    t.ok(pool.allocatedResources.pageInfo.hasPreviousPage);
    let secondResource = pool.allocatedResources.pageInfo.startCursor;

    //get all resources before the 2nd resource
    pool = await getResourcePool(poolId, secondResource,  null, null, 1000);
    t.equal(pool.allocatedResources.edges.length, 1);
    t.equal(resourceIds[0], pool.allocatedResources.edges[0].node.id);
    t.notOk(pool.allocatedResources.pageInfo.hasPreviousPage);

    await cleanup()
    t.end();
});

test('allocation pool test alternative ID', async (t) => {
    const poolId = await createVlanRootPool();
    await claimResourceWithAltId(poolId, {}, {vlanAltId: getUniqueName('first allocation vlan'), order: getUniqueName('first')});

    let altId = {vlanAltId: getUniqueName('second allocation vlan'), order: getUniqueName('second')};
    await claimResourceWithAltId(poolId, {}, altId);

    let altId2 = {vlanAltId: Math.floor(Math.random() * 100000), order: getUniqueName('third')};
    await claimResourceWithAltId(poolId, {}, altId2);

    let altId3 = {vlanAltId: Math.floor(Math.random() * 100000), order: getUniqueName('third')};
    let resource = await claimResourceWithAltId(poolId, {}, altId2);

    await updateResourceAltId(poolId, resource.Properties, altId3)

    let res3 = await queryResourcesByAltIdAndPoolId(poolId, altId3);

    t.equal(res3.edges[0].node.Properties.vlan, 3);
    t.equal(res3.edges[0].node.AlternativeId.vlanAltId, altId3.vlanAltId)
    t.equal(res3.edges[0].node.AlternativeId.order, altId3.order)

    //test nothing found
    let res = await queryResourcesByAltIdAndPoolId(poolId, {someKey: 'this does not exist :('});
    t.notOk(res[0]);

    //test string-only comparison
    res = await queryResourcesByAltIdAndPoolId(poolId, altId);
    t.equal(res.edges[0].node.Properties.vlan, 1);

    //test string-and-number comparison
    res = await queryResourcesByAltIdAndPoolId(poolId, altId2);
    t.equal(res.edges[0].node.Properties.vlan, 2);

    await cleanup()
    t.end();
});

test('test pagination resources', async (t) => {
    const poolId = await createVlanRootPool();
    await claimResourceWithAltId(poolId, {}, {vlanAltId: getUniqueName('first allocation vlan'), order: getUniqueName('first')});

    let altId = {vlanAltId: getUniqueName('second allocation vlan'), order: getUniqueName('second')};
    await claimResourceWithAltId(poolId, {}, altId);

    let altId2 = {vlanAltId: Math.floor(Math.random() * 100000), order: getUniqueName('third')};
    await claimResourceWithAltId(poolId, {}, altId2);

    let altId3 = {vlanAltId: Math.floor(Math.random() * 100000), order: getUniqueName('third')};
    await claimResourceWithAltId(poolId, {}, altId3);

    let resources = await getPaginatedResourcesForPool(poolId, 1, null, null, null);
    t.equal(resources.edges[0].node.Properties.vlan, 0);

    resources = await getPaginatedResourcesForPool(poolId, null, 1, null, null);
    t.equal(resources.edges[0].node.Properties.vlan, 3);

    await cleanup()
    t.end();
});

test('set pool test alternative ID', async (t) => {
    let rtId = await createResourceType(getUniqueName('simplevalue'), {avalue: 'int'} )
    let poolId = await createSetPool(
        getUniqueName('simple_value_set_pool'),
        rtId,
        [{avalue:1}, {avalue:2}, {avalue:3}]
    );
    let altId = {id: getUniqueName('avalue1')};

    await claimResourceWithAltId(poolId, {}, altId);
    let res = await queryResourcesByAltIdAndPoolId(poolId, altId);
    t.equal(res.edges[0].node.Properties.avalue, 1)

    let duplicate = await claimResourceWithAltId(poolId, {}, altId, null, true);
    t.ok(duplicate);

    await cleanup()
    t.end();
});

test('empty pools test', async (t) => {
    const poolId = await createIpv4RootPool('192.168.3.0', 16);
    await createIpv4RootPool('192.168.24.0', 16);

    let emptyPools = await getEmptyPools(null);
    t.equal(emptyPools?.edges?.length, 2)

    await claimResource(poolId, {});
    await claimResource(poolId, {});

    emptyPools = await getEmptyPools(null);
    t.equal(emptyPools?.edges?.length, 1)

    await cleanup()
    t.end();
});

test('test filtering of allocated ipv6 resources by pool id', async (t) => {
    const poolId = await createIpv6PrefixRootPool();
    const poolId2 = await createIpv6PrefixRootPool();

    await claimResourceWithAltId(poolId, { desiredSize: 4 }, {});
    await claimResourceWithAltId(poolId2, { desiredSize: 4 }, {});

    const allocResourcesForPoolId = await queryResourcesByAltIdAndPoolId(poolId, {});
    const allocResourcesForPoolId2 = await queryResourcesByAltIdAndPoolId(poolId2, {});

    t.equal(allocResourcesForPoolId.edges.length, 1);
    t.equal(allocResourcesForPoolId2.edges.length, 1);

    await cleanup();
    t.end();
});

test('test filtering of allocated ipv4 resources by pool id', async (t) => {
    const poolId = await createIpv4PrefixRootPool();
    const poolId2 = await createIpv4PrefixRootPool();

    await claimResourceWithAltId(poolId.id, { desiredSize: 4 }, {});
    await claimResourceWithAltId(poolId2.id, { desiredSize: 4 }, {});

    const allocResourcesForPoolId = await queryResourcesByAltIdAndPoolId(poolId.id, {});
    const allocResourcesForPoolId2 = await queryResourcesByAltIdAndPoolId(poolId2.id, {});

    t.equal(allocResourcesForPoolId.edges.length, 1);
    t.equal(allocResourcesForPoolId2.edges.length, 1);

    await cleanup();
    t.end();
});

test('test filtering of allocated resources when pool id is null', async (t) => {
    const poolId = await createIpv6PrefixRootPool()

    await claimResourceWithAltId(poolId, {
        desiredSize: 4
    }, {});

    const allocResourcesForPoolId = await queryResourcesByAltIdAndPoolId(null, null);

    t.equal(allocResourcesForPoolId, null);

    await cleanup();
    t.end();
});

test('query resources that do not have an alternative id', async (t) => {
    const poolId = await createIpv4RootPool('10.0.0.0', 16);
    await claimResource(poolId, {});
    await claimResourceWithAltId(poolId, {}, {vlanAltId: getUniqueName('vlan')});

    const res = await queryResourcesByAltIdAndPoolId(poolId, {});
    const res2 = await getResourcesForPool(poolId);

    t.equal(res.edges.every((node) => node.AlternativeId != null || node.AlternativeId == null), true);
    t.equal(res2.edges.every((node) => node.AlternativeId != null || node.AlternativeId == null), true);
});

test('claim resource with integer value sent as string', async (t) => {
    const pool = await createIpv4PrefixRootPool();
    const claimedResource = await claimResource(pool.id, { desiredSize: '4' });
    const poolWithResources = await getResourcesForPool(pool.id);

    t.ok(claimedResource);
    t.equal(poolWithResources.edges.length, 1);
});

test('claim resource with altId and integer value sent as string', async (t) => {
    const pool = await createIpv4PrefixRootPool();
    const claimedResource = await claimResourceWithAltId(pool.id, { desiredSize: '4' }, { vlanAltId: '123' });
    const poolWithResources = await getResourcesForPool(pool.id);

    t.ok(claimedResource);
    t.equal(poolWithResources.edges.length, 1);
});

test('validate format of provided ipv4 address when creating pool', async (t) => {
    const pool1 = await createIpv4RootPool("10.0.0.0", 24, false);
    const pool2 = await createIpv4RootPool("10.0.0.1", 24, false);
    const pool3 = await createIpv4RootPool("oidjqwodij", 24, false);

    t.ok(pool1);
    t.notOk(pool2);
    t.notOk(pool3);

    await cleanup();
    t.end();
});

test('validate format of provided ipv6 address when creating pool', async (t) => {
    const pool1 = await createIpv6PrefixRootPool();
    const pool3 = await createIpv6PrefixRootPool("something::1", 64, false);

    t.ok(pool1);
    t.equal(pool3, null);

    await cleanup();
    t.end();
});

test('test connection structure of resource pool query', async (t) => {
    await prepareDataForFiltering();
    const pools = await queryResourcePools();

    t.equal(pools.edges.length, 2);

    await cleanup();
    t.end();
});

test('test pagination of resource pools', async (t) => {
    await prepareDataForFiltering();
    const pools = await queryResourcePools(1);

    t.equal(pools.edges.length, 1);

    const pools1 = await queryResourcePools();
    const startCursor = pools1.pageInfo.startCursor;
    const endCursor = pools1.pageInfo.endCursor;

    const pools2 = await queryResourcePools(1, undefined, undefined, startCursor);
    t.equal(pools2.edges[0].cursor, endCursor);
    t.equal(pools2.pageInfo.hasPreviousPage, true);

    await cleanup();
    t.end();
});

test('test filtering by allocated resource properties', async (t) => {
    await prepareDataForFiltering();
    const byIpv6Addr = await queryResourcePools(undefined, undefined, undefined, undefined, {
        address: "dead::"
    });
    const byIpv6Addr2 = await queryResourcePools(undefined, undefined, undefined, undefined, {
        address: "dead::4"
    });
    const byIpv4Addr = await queryResourcePools(undefined, undefined, undefined, undefined, {
        address: "10.0.0.0"
    });


    t.equal(byIpv6Addr.edges[0].node.allocatedResources.edges[0].node.Properties["address"], "dead::");
    t.equal(byIpv6Addr2.edges[0].node.allocatedResources.edges[1].node.Properties["address"], "dead::4");
    t.equal(byIpv4Addr.edges[0].node.allocatedResources.edges[0].node.Properties["address"], "10.0.0.0");

    t.equal(byIpv6Addr.edges.length, 1);
    t.equal(byIpv6Addr2.edges.length, 1);
    t.equal(byIpv4Addr.edges.length, 1);

    await cleanup();
    t.end();
});

test('test ordering of resource pools by name in descending direction', async (t) => {
    for (let i = 0; i < 3; i++) {
        await createIpv4PrefixRootPoolWithName(`pool-${i}`);
    }

    const pools = await queryResourcePools(undefined, undefined, undefined, undefined, undefined, {
        field: "name",
        direction: "DESC"
    });

    t.equal(pools.edges[0].node.Name.substring(0, 6), "pool-2");

    await cleanup();
    t.end();
});

test('test ordering of resource pools by name in ascending direction', async (t) => {
    for (let i = 0; i < 3; i++) {
        await createIpv4PrefixRootPoolWithName(`pool-${i}`);
    }

    const pools = await queryResourcePools(undefined, undefined, undefined, undefined, undefined, {
        field: "name",
        direction: "ASC"
    });

    t.equal(pools.edges[0].node.Name.substring(0, 6), "pool-0");

    await cleanup();
    t.end();
});

test('test ordering of resource pools by dealocation safety period in descending direction', async (t) => {
    for (let i = 0; i < 3; i++) {
        await createIpv4PrefixRootPool("10.0.0.0", 24, false, i);
    }

    const pools = await queryResourcePools(undefined, undefined, undefined, undefined, undefined, {
        field: "dealocationSafetyPeriod",
        direction: "DESC"
    });

    t.equal(pools.edges[0].node.DealocationSafetyPeriod, 2);

    await cleanup();
    t.end();
});

test('test ordering of resource pools by dealocation safety period in ascending direction', async (t) => {
    for (let i = 0; i < 3; i++) {
        await createIpv4PrefixRootPool("10.0.0.0", 24, false, i);
    }

    const pools = await queryResourcePools(undefined, undefined, undefined, undefined, undefined, {
        field: "dealocationSafetyPeriod",
        direction: "ASC"
    });

    t.equal(pools.edges[0].node.DealocationSafetyPeriod, 0);

    await cleanup();
    t.end();
});

test('test ordering of resource pools by non existing field in ascending direction', async (t) => {
    for (let i = 0; i < 3; i++) {
        await createIpv4PrefixRootPool("10.0.0.0", 24, false, i);
    }

    const pools = await queryResourcePools(undefined, undefined, undefined, undefined, undefined, {
        field: "nonExistingSortKey",
        direction: "ASC"
    });

    t.notOk(pools);

    await cleanup();
    t.end();
});

test('test ordering of resource pools by dealocation safety period in ascending direction with pagination', async (t) => {
    for (let i = 0; i < 3; i++) {
        await createIpv4PrefixRootPool("10.0.0.0", 24, false, i);
    }

    const pools = await queryResourcePools(20, undefined, undefined, undefined, undefined, {
        field: "dealocationSafetyPeriod",
        direction: "ASC"
    });

    t.ok(pools);

    await cleanup();
    t.end();
});
