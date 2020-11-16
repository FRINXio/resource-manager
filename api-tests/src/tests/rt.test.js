import { deleteResourceType, findResourceTypeId, createResourceType, createSetPool, claimResource } from '../graphql-queries';
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



