import {
    claimResource,
    createSingletonPool, createTag, deleteResourcePool,
    findResourceTypeId, createSetPool,
    freeResource,
    getResourcesForPool, searchPoolsByTags, tagPool
} from "../graphql-queries";
import {getUniqueName} from "../test-helpers";

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
