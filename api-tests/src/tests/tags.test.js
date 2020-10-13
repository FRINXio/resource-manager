import {deleteTag, untagPool, tagPool, createTag, searchPoolsByTags} from '../graphql-queries';
import {getTag, createVlanRangeRootPool, getUniqueName} from '../test-helpers';


test('create and delete tag', async () => {
    let tagName = getUniqueName('test tag');
    let tagId = await createTag(tagName);
    expect(await getTag(tagName)).toBeTruthy();

    await deleteTag(tagId);
    expect(await getTag(tagName)).not.toBeTruthy();
});

test('tagging and untagging pool', async () => {
    let tagName = getUniqueName('test tag');
    let tagId = await createTag(tagName);
    let poolId = await createVlanRangeRootPool();

    await tagPool(tagId, poolId);
    let tag = await getTag(tagName);
    expect(tag.Pools[0].id).toBe(poolId);

    await untagPool(tagId, poolId);
    tag = await getTag(tagName);
    expect(tag.Pools).toHaveLength(0);
});

test('searching pools via tags', async () => {
    const tagName1 = getUniqueName('test tag');
    const tagName2 = getUniqueName('test tag');
    const tag1Id = await createTag(tagName1);
    const tag2Id = await createTag(tagName2);
    const poolId = await createVlanRangeRootPool();

    await tagPool(tag1Id, poolId);
    await tagPool(tag2Id, poolId);

    let matchedPools = await searchPoolsByTags({matchesAny: [{matchesAll: [tagName1, tagName2]}]});
    expect(matchedPools[0].id).toBe(poolId);
    expect(matchedPools).toHaveLength(1);

    matchedPools = await searchPoolsByTags({matchesAny: [{matchesAll: [tagName1]},{matchesAll: [tagName2]}]});
    expect(matchedPools[0].id).toBe(poolId);
    expect(matchedPools).toHaveLength(1);

    matchedPools = await searchPoolsByTags({matchesAny: [{matchesAll: [tagName2]}]});
    expect(matchedPools[0].id).toBe(poolId);
    expect(matchedPools).toHaveLength(1);
});
