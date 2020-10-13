import {
    claimResource,
    createSingletonPool, createTag, deleteResourcePool,
    findResourceTypeId, createSetPool,
    freeResource,
    getResourcesForPool, searchPoolsByTags, tagPool, getCapacityForPool
} from "../graphql-queries";
import {
    createIpv4PrefixRootPool,
    createIpv4RootPool,
    createIpv6PrefixRootPool,
    createIpv6RootPool, createRandomIntRootPool, createRdRootPool,
    createVlanRangeRootPool,
    createVlanRootPool,
    getUniqueName
} from "../test-helpers";

test('singleton claim and free resource', async () => {
    let rtId = await findResourceTypeId('ipv4');
    const ipAddress = '192.168.1.1';
    let poolId = await createSingletonPool(
        getUniqueName('singleton'),
        rtId,
        [{address: ipAddress}]
    );
    let resource = await claimResource(poolId, {});
    let rs = await getResourcesForPool(poolId);
    expect(rs).toHaveLength(1)
    expect(rs[0].Properties.address).toBe(ipAddress)
    await freeResource(poolId, resource.Properties);

    rs = await getResourcesForPool(poolId);
    expect(rs).toHaveLength(0)
});

test('create and delete singleton pool', async () => {
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
    expect(foundPool).toHaveLength(1);
    expect(foundPool[0].id).toBe(poolId);

    let resource1 = await claimResource(poolId, {});
    let resource2 = await claimResource(poolId, {});

    expect(resource1).toStrictEqual(resource2); //the same resource

    await freeResource(poolId, resource2.Properties);

    await deleteResourcePool(poolId);
    foundPool = await searchPoolsByTags({matchesAny: [{matchesAll: [tagText]}]});
    expect(foundPool).toHaveLength(0);
});

test('create and delete resources in set pool', async () => {
    let rtId = await findResourceTypeId('ipv4');
    let poolId = await createSetPool(
        getUniqueName('singleton'),
        rtId,
        [{address: '192.168.1.1'}, {address: '192.168.1.2'}]
    );
    let resource = await claimResource(poolId, {});
    let rs = await getResourcesForPool(poolId);
    expect(rs).toHaveLength(1);
    await freeResource(poolId, resource.Properties)
    rs = await getResourcesForPool(poolId);
    expect(rs).toHaveLength(0);
});

test('create and delete set pool', async () => {
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
    expect(foundPool).toHaveLength(1);
    expect(foundPool[0].id).toBe(poolId);
    let resource = await claimResource(poolId, {});
    await freeResource(poolId, resource.Properties)
    await deleteResourcePool(poolId);
    foundPool = await searchPoolsByTags({matchesAny: [{matchesAll: [tagText]}]});
    expect(foundPool).toHaveLength(0);
});

test('capacity for allocating vlan-range pool', async () => {
    const poolId = await createVlanRangeRootPool();

    await claimResource(poolId, {desiredSize: 1});
    await claimResource(poolId, {desiredSize: 1});
    await claimResource(poolId, {desiredSize: 3});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(5);
    expect(capacity.freeCapacity).toBe(4091);
});

test('capacity for allocating vlan pool', async () => {
    const poolId = await createVlanRootPool();

    await claimResource(poolId, {});
    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(3);
    expect(capacity.freeCapacity).toBe(4093);
});

test('capacity for allocating ipv6-prefix pool', async () => {
    const poolId = await createIpv6PrefixRootPool();

    await claimResource(poolId, {desiredSize: 4})
    await claimResource(poolId, {desiredSize: 4});
    await claimResource(poolId, {desiredSize: 4});
    await claimResource(poolId, {desiredSize: 4});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(16);
    expect(capacity.freeCapacity).toBe(240);
});

test('capacity for allocating ipv6 pool', async () => {
    const poolId = await createIpv6RootPool();

    await claimResource(poolId, {})
    await claimResource(poolId, {});
    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(4);
    expect(capacity.freeCapacity).toBe(5.192296858534828e+33);
});

test('capacity for allocating ipv4 pool', async () => {
    const poolId = await createIpv4RootPool('192.168.3.0', 16);

    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(2);
    expect(capacity.freeCapacity).toBe(65532);
});

test('capacity for random pool', async () => {
    const poolId = await createRandomIntRootPool();

    await claimResource(poolId, {});
    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(2);
    expect(capacity.freeCapacity).toBe(997);
});

test('capacity for allocating ipv4-prefix pool', async () => {
    const poolId = await createIpv4PrefixRootPool();

    await claimResource(poolId, {desiredSize: 2});
    await claimResource(poolId, {desiredSize: 2});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(4);
    expect(capacity.freeCapacity).toBe(16777210);
});

test('capacity for allocating RD pool', async () => {
    const rdPoolId = await createRdRootPool();

    await claimResource(rdPoolId, {asNumber: 1985, assignedNumber: 5891});
    await claimResource(rdPoolId, {asNumber: 2020, assignedNumber: 2020});
    await claimResource(rdPoolId, {asNumber: 47, assignedNumber: 47});

    const capacity = await getCapacityForPool(rdPoolId);
    expect(capacity.utilizedCapacity).toBe(3);
    expect(capacity.freeCapacity).toBe(281474976710656);
});

test('capacity for set pool', async () => {
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
    expect(capacity.utilizedCapacity).toBe(3);
    expect(capacity.freeCapacity).toBe(1);
});

test('capacity for singleton pool', async () => {
    let rtId = await findResourceTypeId('ipv4');
    let poolId = await createSingletonPool(
        getUniqueName('singleton'),
        rtId,
        [{address: '192.168.1.1'}]
    );

    await claimResource(poolId, {});

    const capacity = await getCapacityForPool(poolId);
    expect(capacity.utilizedCapacity).toBe(1);
    expect(capacity.freeCapacity).toBe(0);
});