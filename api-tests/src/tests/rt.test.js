import { deleteResourceType, deleteResourcePool, getResourcesForPool, freeResource, findResourceTypeId, createSingletonPool, createResourceType, createSetPool, claimResource } from '../graphql-queries';
import {getUniqueName} from '../test-helpers';

test('create resource type, set pool and claim a resource', async () => {
    let complexId = await createResourceType('complex', {a: 'int', b: 'string'} )
    let poolId = await createSetPool(
        getUniqueName('complex'),
        complexId,
        [{a:11, b:'eleven'}, {a:17, b:'seventeen'}, {a:1, b:'one'}]
        );
    let resource = await claimResource(poolId, {});
    expect(resource.Properties.a).toBe(11)
    expect(resource.Properties.b).toBe('eleven')
});

test('create delete resourceType', async () => {
    let resourceName = getUniqueName('testtype');
    await createResourceType(resourceName, {abc: 'int'} )
    let rtId = await findResourceTypeId(resourceName);
    expect(rtId).toBeTruthy();
    await deleteResourceType(rtId);
    rtId = await findResourceTypeId(resourceName);
    expect(rtId).not.toBeTruthy()
});


//TODO this test fails because deleting singleton pools does not work yet
test.skip('singleton pool API test', async () => {
    let rtId = await findResourceTypeId('ipv4');
    const ipAddress = '192.168.1.1';
    let poolId = await createSingletonPool(
        getUniqueName('singleton'),
        rtId,
        [{address: ipAddress}]
    );
    await claimResource(poolId, {});
    let resource = await claimResource(poolId, {});
    await freeResource(poolId, resource.Properties);
    let rs = await getResourcesForPool(poolId);
    expect(rs).toHaveLength(1)
    expect(rs[0].Properties.address).toBe(ipAddress)
    await freeResource(poolId, resource.Properties);
    await deleteResourcePool(poolId);
    rs = await getResourcesForPool(poolId);
    expect(rs).toHaveLength(0)
});

