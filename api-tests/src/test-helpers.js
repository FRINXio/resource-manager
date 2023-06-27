// Contains helper functions used in tests (so that test functions are more readable)

import {
    claimResource,
    createAllocationPool,
    createNestedAllocationPool,
    createNestedSingletonPool,
    deleteResourcePool,
    findAllocationStrategyId,
    findResourceTypeId,
    freeResource, getAllPoolsByTypeOrTag,
    getAllTags,
    getLeafPools,
    getResourcePool,
    getResourcesForPool
} from "./graphql-queries.js";
import _ from "underscore";

export async function getTag(tagName) {
    let allTags = await getAllTags();
    return _.find(allTags, d => d.Tag === tagName);
}

export async function createRandomIntRootPool() {
    let resourceTypeId = await findResourceTypeId('random_signed_int32');
    let strategyId = await findAllocationStrategyId('random_signed_int32');
    let poolName = getUniqueName('root-random-int');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {from: "int", to: "int"},
        {from: '1', to: 999},)

    return pool.id;
}

export async function createIpv4RootPool(address, prefix, subnet = false) {
    let resourceTypeId = await findResourceTypeId('ipv4');
    let strategyId = await findAllocationStrategyId('ipv4');
    let poolName = getUniqueName('root-ipv4');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {address: "string", prefix: "int", subnet: "bool"},
        {address: address, prefix: prefix, subnet: subnet})

    if (pool == null) {
        return null;
    }

    return pool.id
}

export async function createRdRootPool() {
    let resourceTypeId = await findResourceTypeId('route_distinguisher');
    let strategyId = await findAllocationStrategyId('route_distinguisher');
    let poolName = getUniqueName('root-rd');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {rd: "int"},
        {rd: 0},)

    return pool.id;
}


export async function createIpv6RootPool() {
    let resourceTypeId = await findResourceTypeId('ipv6');
    let strategyId = await findAllocationStrategyId('ipv6');
    let poolName = getUniqueName('root-ipv6-range');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {address: "string", prefix: "int", subnet: "bool"},
        {address: "dead::", prefix: 16, subnet: false})

    if (pool == null) {
        return null;
    }

    return pool.id;
}

export async function createIpv6PrefixRootPool(address = "dead::", prefix = 120, subnet = false) {
    let resourceTypeId = await findResourceTypeId('ipv6_prefix');
    let strategyId = await findAllocationStrategyId('ipv6_prefix');
    let poolName = getUniqueName('root-ipv6-range');
    const pool = await createAllocationPool(
        poolName,
        resourceTypeId,
        strategyId,
        {address: "string", prefix: "int", subnet: "bool"},
        {address, prefix, subnet})

    if (pool == null) {
        return null;
    }

    return pool.id;
}


export async function createIpv6NestedPool(parentResourceId) {
    let resourceTypeId = await findResourceTypeId('ipv6');
    let strategyId = await findAllocationStrategyId('ipv6');
    return await createNestedAllocationPool(
        getUniqueName('nested-ipv6'),
        resourceTypeId,
        strategyId,
        parentResourceId);
}

export async function createIpv4PrefixRootPool(address = "10.0.0.0", prefix = 8, subnet = false) {
    let resourceTypeId = await findResourceTypeId('ipv4_prefix');
    let strategyId = await findAllocationStrategyId('ipv4_prefix');
    return await createAllocationPool(
        getUniqueName('ipv4-root'),
        resourceTypeId,
        strategyId,
        {prefix: "int", address: "string", subnet: "bool"},
        {prefix: prefix, address: address, subnet: subnet});
}

export async function createUniqueIdPool() {
    let resourceTypeId = await findResourceTypeId('unique_id');
    let strategyId = await findAllocationStrategyId('unique_id');
    return await createAllocationPool(
        getUniqueName('unique_id'),
        resourceTypeId,
        strategyId,
        {from: "int", to: "int", idFormat: "string"},
        {from: 1, to: 15, idFormat: "{counter}"});
}

export async function createIpv4NestedPool(parentResourceId) {
    let resourceTypeId = await findResourceTypeId('ipv4');
    let strategyId = await findAllocationStrategyId('ipv4');
    return await createNestedAllocationPool(
        getUniqueName('ipv4-nested'),
        resourceTypeId,
        strategyId,
        parentResourceId);
}

export async function createIpv4PrefixNestedPool(parentResourceId) {
    let resourceTypeId = await findResourceTypeId('ipv4_prefix');
    let strategyId = await findAllocationStrategyId('ipv4_prefix');
    return await createNestedAllocationPool(
        getUniqueName('ipv4prefix-nested'),
        resourceTypeId,
        strategyId,
        parentResourceId);
}

export async function get2ChildrenIds(poolId) {
    let resourceForPool = await getResourcesForPool(poolId);
    return [resourceForPool.edges[0].node.NestedPool.id, resourceForPool.edges[1].node.NestedPool.id]
}

export async function createSingletonIpv4PrefixNestedPool(parentResourceId) {
    let resourceTypeId = await findResourceTypeId('ipv4_prefix');
    return await createNestedSingletonPool(
        getUniqueName('singleton-ipv4prefix-nested'),
        resourceTypeId,
        [{address: "10.10.0.0", prefix: 11, subnet:false}],
        parentResourceId);
}

export async function createVlanRangeRootPool() {
    let resourceTypeId = await findResourceTypeId('vlan_range');
    let strategyId = await findAllocationStrategyId('vlan_range');
    const pool = await createAllocationPool(
        getUniqueName('root-vlan-range'),
        resourceTypeId,
        strategyId,
        {from: "int", to: "int"},
        {from: 0, to: 4095},)

    return pool.id;
}

export async function createVlanRootPool(tags = null) {
    let resourceTypeId = await findResourceTypeId('vlan');
    let strategyId = await findAllocationStrategyId('vlan');
    const pool = await createAllocationPool(
        getUniqueName('root-vlan'),
        resourceTypeId,
        strategyId,
        {from: "int", to: "int"},
        {from: 0, to: 4095},
        tags)

    return pool.id;
}

export async function createVlanNestedPool(parentResourceId, tags = null) {
    let resourceTypeId = await findResourceTypeId('vlan');
    let strategyId = await findAllocationStrategyId('vlan');
    return await createNestedAllocationPool(
        getUniqueName('nested'),
        resourceTypeId,
        strategyId,
        parentResourceId,
        tags)
}

export async function createRandomSignedInt32Pool({from, to}) {
    let resourceTypeId = await findResourceTypeId('random_signed_int32');
    let strategyId = await findAllocationStrategyId('random_signed_int32');
    return await createAllocationPool(
        getUniqueName('random-signed-int32'),
        resourceTypeId,
        strategyId,
        {from: "int", to: "int"},
        {from, to});
}

export function getUniqueName(prefix) {
    return prefix + Math.random().toString(36).substring(5);
}

export async function prepareIpv4Pool() {
    let rootPoolId = (await createIpv4PrefixRootPool()).id;
    let resource12Id = (await claimResource(rootPoolId, {desiredSize: 4194304})).id;
    return await createIpv4NestedPool(resource12Id);
}

export async function allocateFromIPv4PoolSerially(poolId, count, claimResourceParams) {
    const result = [];
    for (let i = 0; i < count; i++) {
        result.push(await claimResource(poolId, claimResourceParams));
    }
    const ips = result.filter(it => it != null).map(it => it.Properties.address);
    const uniqIPs = [...new Set(ips)];
    if (uniqIPs.length != count) {
        console.error({result, ips, uniqIPs});
        throw new Error("Unexpected number of IPs:" + uniqIPs.length);
    }
    return uniqIPs;
}

export async function allocateFromIPv4PoolParallelly(poolId, count, retries, claimResourceParams) {
    const promises = [];
    for (let i = 0; i < count; i++) {
        promises.push(claimResource(poolId, claimResourceParams, null, true));
    }
    const result = await Promise.all(promises);
    const ips = result.filter(it => it != null).map(it => it.Properties.address);
    if (ips.length != count) {
        if (ips.length == 0) {
            throw new Error("Nothing alocated at " + JSON.stringify({count, retries}))
        } else if (retries-- > 0) {
            // call recursively
            const moreIPs = await allocateFromIPv4PoolParallelly(poolId, count - ips.length, retries, claimResourceParams);
            ips.push(...moreIPs);
        } else {
            console.error({result, ips});
            throw new Error("Unexpected number of IPs:" + ips.length);
        }
    }
    return ips;
}

export async function queryIPs(poolId, count) {
    const queried = await getResourcePool(poolId, null, null, count);
    const allocated = queried.allocatedResources;
    if (!allocated || !allocated.edges) {
        console.error("Unexpected query result", {queried});
        throw new Error("Unexpected query result");
    }
    const foundIPs = allocated.edges.map(it => it.node.Properties.address);
    const uniqIPs = [...new Set(foundIPs)];
    if (foundIPs.length != uniqIPs.length) {
        console.error("Query returned duplicates", {foundIPs, uniqIPs});
        throw new Error("Query returned duplicates");
    }
    return uniqIPs;
}

export async function cleanup() {
    const resourceIds = new Set()
    resourceIds.add(await findResourceTypeId('ipv4'));
    resourceIds.add(await findResourceTypeId('ipv4_prefix'));
    resourceIds.add(await findResourceTypeId('ipv6'));
    resourceIds.add(await findResourceTypeId('ipv6_prefix'));
    resourceIds.add(await findResourceTypeId('vlan'));
    resourceIds.add(await findResourceTypeId('vlan_range'));
    resourceIds.add(await findResourceTypeId('route_distinguisher'));
    resourceIds.add(await findResourceTypeId('random_signed_int32'));
    resourceIds.add(await findResourceTypeId('unique_id'));

    for (let i = 0; i < resourceIds.size; i++) {
        let pools = await getAllPoolsByTypeOrTag(resourceIds[i]);
        let leafPools = await getLeafPools(resourceIds[i]);

        while (leafPools?.edges?.length > 0) {
            for (let j = 0; j < leafPools?.edges?.length; j++) {
                pools = pools?.edges?.filter(function (e) {return e.node.id !== leafPools?.edges?.[j].node.id});
                await cleanPool(leafPools?.edges?.[j])
            }
            leafPools = await getLeafPools(resourceIds[i]);
        }
        for (let j = 0; j < pools?.edges?.length; j++) {
            await cleanPool(pools?.edges?.[j]);
        }
    }
}

export async function cleanPool(pool) {
    for (let j = 0; j < pool.node.Resources.length; j++) {
        await freeResource(pool.node.id, pool.node.Resources[j].Properties);
    }
    await deleteResourcePool(pool.node.id);
}

export async function prepareDataForFiltering() {
    const ipv6 = await createIpv6PrefixRootPool("dead::", 64, true);
    const ipv4 = await createIpv4RootPool("10.0.0.0", 24, false);

    await claimResource(ipv6, {
        desiredSize: 2
    });
    await claimResource(ipv6, {
        desiredSize: 2
    });

    await claimResource(ipv4, {
        desiredSize: 2
    });
    await claimResource(ipv4, {
        desiredSize: 2
    });

    return {
        ipv4,
        ipv6
    }
}