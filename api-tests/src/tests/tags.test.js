import {deleteTag, untagPool, tagPool, createTag, searchPoolsByTags} from '../graphql-queries.js';
import {
    getTag,
    createVlanRangeRootPool,
    getUniqueName,
    createVlanRootPool
} from '../test-helpers.js';
import tap from 'tap';
const test = tap.test;

test('create and delete tag', async (t) => {
    let tagName = getUniqueName('test tag');
    let tagId = await createTag(tagName);
    t.ok(await getTag(tagName));

    await deleteTag(tagId);
    t.notOk(await getTag(tagName));
    t.end();
});

test('tagging and untagging pool', async (t) => {
    let tagName = getUniqueName('test tag');
    let tagId = await createTag(tagName);
    let poolId = await createVlanRangeRootPool();

    await tagPool(tagId, poolId);
    let tag = await getTag(tagName);
    t.equal(tag.Pools[0].id, poolId);

    await untagPool(tagId, poolId);
    tag = await getTag(tagName);
    t.equal(tag.Pools.length, 0);
    t.end();
});

test('searching pools via tags', async (t) => {
    const tagName1 = getUniqueName('test tag');
    const tagName2 = getUniqueName('test tag');
    const tag1Id = await createTag(tagName1);
    const tag2Id = await createTag(tagName2);
    const poolId = await createVlanRangeRootPool();

    await tagPool(tag1Id, poolId);
    await tagPool(tag2Id, poolId);

    let matchedPools = await searchPoolsByTags({matchesAny: [{matchesAll: [tagName1, tagName2]}]});
    t.equal(matchedPools[0].id, poolId);
    t.equal(matchedPools.length, 1);

    matchedPools = await searchPoolsByTags({matchesAny: [{matchesAll: [tagName1]},{matchesAll: [tagName2]}]});
    t.equal(matchedPools[0].id, poolId);
    t.equal(matchedPools.length, 1);

    matchedPools = await searchPoolsByTags({matchesAny: [{matchesAll: [tagName2]}]});
    t.equal(matchedPools[0].id, poolId);
    t.equal(matchedPools.length, 1);

    // test createPool and tag in a single operation
    const tag1 = getUniqueName("tag1")
    const tag2 = getUniqueName("tag2")
    const taggedPoolId = await createVlanRootPool([tag1, tag2]);
    matchedPools = await searchPoolsByTags({matchesAny: [{matchesAll: [tag1, tag2]}]});
    t.equal(matchedPools.length, 1);
    t.equal(matchedPools[0].id, taggedPoolId);
    t.end();
});
