import { deleteResourceType, findResourceTypeId, createResourceType, createSetPool, claimResource } from '../graphql-queries.js';
import {getUniqueName} from '../test-helpers.js';
import tap from 'tap';
const test = tap.test;

test('create resource type, set pool and claim a resource', async (t) => {
    let complexId = await createResourceType('complex', {a: 'int', b: 'string'} )
    let poolId = await createSetPool(
        getUniqueName('complex'),
        complexId,
        [{a:11, b:'eleven'}, {a:17, b:'seventeen'}, {a:1, b:'one'}]
        );
    let resource = await claimResource(poolId, {});
    t.equal(resource.Properties.a, 11)
    t.equal(resource.Properties.b, 'eleven');
    t.end();
});

test('create delete resourceType', async (t) => {
    let resourceName = getUniqueName('testtype');
    await createResourceType(resourceName, {abc: 'int'} )
    let rtId = await findResourceTypeId(resourceName);
    t.ok(rtId);
    await deleteResourceType(rtId);
    rtId = await findResourceTypeId(resourceName);
    t.notOk(rtId);
    t.end();
});



